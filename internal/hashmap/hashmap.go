package hashmap

import "github.com/cespare/xxhash/v2"

type entry struct {
	key   string
	value any
}

type bucket struct {
	entries [8]entry
	tophash [8]uint8
	next    *bucket
}

type Entry struct {
	Key   string
	Value any
}

type Map struct {
	buckets []*bucket
	size    int
}

func New() *Map {
	buckets := make([]*bucket, 8)
	for i := range buckets {
		buckets[i] = &bucket{}
	}
	return &Map{buckets: buckets}
}

func hashKey(key string) uint64 {
	return xxhash.Sum64String(key)
}

func (hm *Map) Set(key string, value any) {
	if float32(hm.size)/float32(len(hm.buckets)) > 0.75 {
		hm.resize()
	}
	myHash := hashKey(key)
	topHash := uint8((myHash >> 56) + 2)
	idx := int(myHash & uint64(len(hm.buckets)-1))

	b := hm.buckets[idx]
	for {
		for i := range b.tophash {
			if b.tophash[i] < 2 {
				b.entries[i] = entry{key: key, value: value}
				b.tophash[i] = topHash
				hm.size++
				return
			}
			if b.tophash[i] == topHash && b.entries[i].key == key {
				b.entries[i].value = value
				return
			}
		}
		if b.next == nil {
			b.next = &bucket{}
		}
		b = b.next
	}
}

func (hm *Map) Get(key string) (any, bool) {
	myHash := hashKey(key)
	topHash := uint8((myHash >> 56) + 2)
	idx := int(myHash & uint64(len(hm.buckets)-1))

	b := hm.buckets[idx]
	for b != nil {
		for i := range b.tophash {
			if b.tophash[i] >= 2 && b.tophash[i] == topHash && b.entries[i].key == key {
				return b.entries[i].value, true
			}
		}
		b = b.next
	}
	return nil, false
}

func (hm *Map) Delete(key string) bool {
	myHash := hashKey(key)
	topHash := uint8((myHash >> 56) + 2)
	idx := int(myHash & uint64(len(hm.buckets)-1))

	b := hm.buckets[idx]
	for b != nil {
		for i := range b.tophash {
			if b.tophash[i] >= 2 && b.tophash[i] == topHash && b.entries[i].key == key {
				b.tophash[i] = 1
				b.entries[i] = entry{}
				hm.size--
				return true
			}
		}
		b = b.next
	}
	return false
}

func (hm *Map) Len() int {
	return hm.size
}

func (hm *Map) Entries() []Entry {
	entries := make([]Entry, 0, hm.size)
	for _, bucket := range hm.buckets {
		for b := bucket; b != nil; b = b.next {
			for i := range b.tophash {
				if b.tophash[i] >= 2 {
					entries = append(entries, Entry{
						Key:   b.entries[i].key,
						Value: b.entries[i].value,
					})
				}
			}
		}
	}
	return entries
}

func (hm *Map) resize() {
	newBuckets := make([]*bucket, len(hm.buckets)*2)
	for i := range newBuckets {
		newBuckets[i] = &bucket{}
	}
	for _, b := range hm.buckets {
		for ; b != nil; b = b.next {
			for i := range b.tophash {
				if b.tophash[i] >= 2 {
					e := b.entries[i]
					myHash := hashKey(e.key)
					topHash := uint8((myHash >> 56) + 2)
					idx := int(myHash & uint64(len(newBuckets)-1))

					nb := newBuckets[idx]
					for {
						placed := false
						for j := range nb.tophash {
							if nb.tophash[j] < 2 {
								nb.tophash[j] = topHash
								nb.entries[j] = e
								placed = true
								break
							}
						}
						if placed {
							break
						}
						if nb.next == nil {
							nb.next = &bucket{}
						}
						nb = nb.next
					}
				}
			}
		}
	}
	hm.buckets = newBuckets
}
