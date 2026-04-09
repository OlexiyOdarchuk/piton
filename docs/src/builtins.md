# Vbudovani Funktsii

Цей розділ описує все, що вже реалізовано в інтерпретаторі як builtin або спеціальний оператор.

## Matematyka

- `korin x`
- `loh10 x`
- `abs x`
- `arksyn x`
- `kosynus x`
- `zaokruhlennya(x, precision)`
- `stupin` як інфіксний оператор

```piton
drukuvaty korin 9
drukuvaty loh10 100
drukuvaty 2 stupin 3
drukuvaty zaokruhlennya(123.4567, 2)
```

`korin`, `loh10`, `abs`, `arksyn`, `kosynus` поводяться як префіксні оператори. `zaokruhlennya()` викликається як звичайна функція.

## Kolektsii

- `dovzhyna(x)` для списків і словників
- `dodaty(list, itemOrList)` для списків
- `delete(slovnyk, key)` для словників

```piton
s = [1, 2]
s = dodaty(s, 3)
drukuvaty dovzhyna(s)

m = {"a": 1, "b": 2}
delete(m, "a")
drukuvaty m
```

## Chas i vypadkovist

- `vypadkovo(max)` повертає випадкове число з проміжку `[0, max)`
- `vypadkovo(min, max)` повертає число з проміжку `[min, max)`
- `vypadkovo(list)` повертає випадковий елемент списку
- `chas()`
- `zatrymka(seconds)`

```piton
drukuvaty chas()
drukuvaty vypadkovo(10)
drukuvaty vypadkovo(5, 10)
drukuvaty vypadkovo([10, 20, 30])
```

`chas()` повертає час у секундах як число. `zatrymka(x)` призупиняє виконання.

## Kolory

`kolor(name, value)` повертає ANSI-рядок для кольорового друку.

```piton
drukuvaty kolor("red", "uvaha")
drukuvaty kolor("bright_cyan", 123)
```

Назва кольору має бути рядком. Якщо термінал не підтримує ANSI або `NO_COLOR` увімкнений, колір може не відобразитися.

Доступні назви кольорів:

- `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`
- `bright_black`, `bright_red`, `bright_green`, `bright_yellow`
- `bright_blue`, `bright_magenta`, `bright_cyan`, `bright_white`
