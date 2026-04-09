# Tur po Movi

Цей розділ проходить по головних можливостях Piton на одному компактному прикладі.

```piton
sessions = [35, 40, 28]
target = 40

functia classify(value):
    yaksho value > target:
        vernuty "bilshe za plan"
    inackshe:
        vernuty "v mezhakh planu"

i = 0
poky i < dovzhyna(sessions):
    current = sessions[i]
    drukuvaty current
    drukuvaty classify(current)
    i = i + 1
```

## Shcho tut vidbuvayetsya

1. `sessions` - це список чисел.
2. `target` - звичайна змінна.
3. `functia classify(value):` оголошує функцію з одним параметром.
4. `yaksho ... inackshe` працює як `if/else`.
5. `poky` виконує тіло циклу, поки умова істинна.
6. `dovzhyna(sessions)` повертає кількість елементів у списку.
7. `sessions[i]` читає елемент за індексом.

## Mentalna model movy

Piton варто сприймати як невелику динамічну мову з такими правилами:

- усе крутиться навколо виразів, присвоєнь і викликів функцій
- колекції мутабельні: елементи списків і словників можна змінювати
- типи перевіряються під час виконання, а не на етапі парсингу
- помилки не кидають винятки у стилі Python, а друкуються повідомленням інтерпретатора

Після цього розділу переходь до [Bazovyi Syntaksys](./syntax.md) або одразу до [Operatory ta Porivnyannya](./operators.md).
