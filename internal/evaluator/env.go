package evaluator

type Environment struct {
	outer *Environment
	store map[string]any
}

func NewEnv(outer *Environment) *Environment {
	return &Environment{store: make(map[string]interface{}), outer: outer}
}

func (e *Environment) Get(name string) (interface{}, bool) {
	v, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return v, ok
}

func (e *Environment) Set(name string, val interface{}) {
	e.store[name] = val
}

func (e *Environment) ForEach(fn func(name string, val interface{})) {
	for name, val := range e.store {
		fn(name, val)
	}
}
