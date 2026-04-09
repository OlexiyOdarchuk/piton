# Operatory ta Porivnyannya

## Aryfmetychni operatory

- `+`
- `-`
- `*`
- `/`
- `stupin`

```piton
drukuvaty 1 + 2
drukuvaty 10 - 3
drukuvaty 4 * 5
drukuvaty 8 / 2
drukuvaty 2 stupin 3
```

Для рядків:

- `+` виконує конкатенацію
- `*` у парі `string * number` або `number * string` повторює рядок

```piton
drukuvaty "hello" + "!"
drukuvaty "ha" * 3
drukuvaty 2 * "go"
```

## Porivnyannya

- `>`
- `<`
- `>=`
- `<=`
- `==`
- `!=`

```piton
drukuvaty 5 > 3
drukuvaty 5 == 5
drukuvaty 3 != 4
```

Поточна реалізація операторів порівняння орієнтована на числові значення.

## Logika

- `ta`
- `abo`
- `ne`

```piton
drukuvaty true ta false
drukuvaty true abo false
drukuvaty ne true
```

`ta` і `abo` очікують булеві значення з обох боків.

## Prefixni matematychni operatory

Деякі математичні дії виглядають як префіксні оператори:

- `korin`
- `loh10`
- `abs`
- `arksyn`
- `kosynus`

```piton
drukuvaty korin 9
drukuvaty abs -5
drukuvaty loh10 100
```
