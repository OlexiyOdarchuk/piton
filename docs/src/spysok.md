# Spysok

Список задається через `[]` і може містити значення різних типів.

```piton
s = [1, 2, "tri", true]
drukuvaty s
```

## Dostup za indeksom

```piton
s = [10, 20, 30]
drukuvaty s[0]
```

Індекс має бути цілим числом у межах списку.

## Zmina elementa

```piton
s = [10, 20, 30]
s[1] = 99
drukuvaty s
```

## Zrizy

Підтримуються:

- `s[start:end]`
- `s[:end]`
- `s[start:]`
- `s[:]`

```piton
s = [1, 2, 3, 4]
drukuvaty s[1:3]
drukuvaty s[:2]
drukuvaty s[2:]
```

## Dodavannya

`dodaty(list, x)` додає один елемент або зливає два списки.

```piton
a = [1, 2]
b = [3, 4]
drukuvaty dodaty(a, b)
drukuvaty dodaty(a, 5)
```

## Vydalennya elementa

Окремого `delete` для списків зараз немає. Типовий спосіб - склеїти два зрізи:

```piton
s = [1, 2, 3, 4]
trimmed = dodaty(s[:2], s[3:])
drukuvaty trimmed
```
