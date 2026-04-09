# Typy Danykh

Piton - динамічно типізована мова. Основні типи значень:

- `float64` для чисел
- `string`
- `bool`
- `[]` для списків
- `{}` для словників

## Chysla

```piton
a = 10
b = 2.5
drukuvaty a / b
```

## Boolean

```piton
enabled = true
drukuvaty ne enabled
```

Доступні обидва значення:

```piton
drukuvaty true
drukuvaty false
```

## Ryadky

```piton
name = "Piton"
drukuvaty "Hello " + name
```

## Kolektsii

```piton
s = [1, "two", true]
m = {"name": "Piton", "year": 2026}
```

Списки та словники можна вкладати один в одного.

## Shcho potribno pamyataty

- окремих декларацій типів немає
- одна й та сама змінна може містити різні типи в різний момент часу
- словники очікують **рядкові ключі**
