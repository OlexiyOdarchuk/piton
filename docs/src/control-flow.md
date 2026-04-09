# Keruvannya Potokom

## Yaksho / Inackshe

```piton
score = 75

yaksho score >= 90:
    drukuvaty "vidminno"
inackshe yaksho score >= 60:
    drukuvaty "zarakhovano"
inackshe:
    drukuvaty "nezarakhovano"
```

Для `else if` використовується форма `inackshe yaksho`:

```piton
yaksho x > 10:
    drukuvaty "big"
inackshe yaksho x > 5:
    drukuvaty "medium"
inackshe:
    drukuvaty "small"
```

## Poky

```piton
i = 0

poky i < 3:
    drukuvaty i
    i = i + 1
```

Умови очікують булеві вирази. Для складніших умов комбінуй `ta`, `abo` і `ne`.

## Typovyi pattern dlya spyskiv

У мові зараз немає `for`, тому обхід списку робиться через `poky`:

```piton
i = 0
poky i < dovzhyna(items):
    drukuvaty items[i]
    i = i + 1
```
