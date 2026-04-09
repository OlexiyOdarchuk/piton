# Bazovyi Syntaksys

## Zminni

Окремого ключового слова для оголошення немає. Перше присвоєння одразу створює змінну:

```piton
x = 10
name = "Piton"
ready = true
```

Булеві значення в мові - це `true` і `false`.

## Komentari

Однорядковий коментар починається з `#`:

```piton
# tse komentar
drukuvaty "ok"
```

## Vidstupy i bloky

Блоки в Piton задаються відступами, як у Python:

```piton
yaksho true:
    drukuvaty "inside"

drukuvaty "outside"
```

Без правильних відступів програма не розбереться коректно.

## Vyrazy

Підтримуються:

- `+`, `-`, `*`, `/`
- `stupin`
- `>`, `<`, `>=`, `<=`, `==`, `!=`
- `ta`, `abo`, `ne`

```piton
drukuvaty 2 + 3
drukuvaty 2 stupin 3
drukuvaty ne false
drukuvaty (2 < 3) ta true
```

## Ryadky

Рядки задаються в подвійних лапках:

```piton
msg = "Pryvit"
drukuvaty msg + "!"
```

`+` також уміє конкатенувати рядки з числами та булевими значеннями.

Оператор `*` теж має спеціальну поведінку: рядок можна повторити числом.

```piton
drukuvaty "ha" * 3
```

## Vstup i vyvid

```piton
drukuvaty "Vvedy chislo"
vvid x
drukuvaty x
```

Поточна реалізація `vvid` читає число. Для рядкового введення окремого режиму поки немає.
