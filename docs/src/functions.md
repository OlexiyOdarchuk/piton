# Funktsii

Функція оголошується через `functia`, повернення значення робиться через `vernuty`.

```piton
functia add(a, b):
    vernuty a + b

drukuvaty add(2, 3)
```

## Parametry i vyklyk

```piton
functia repeat(label, count):
    i = 0
    poky i < count:
        drukuvaty label
        i = i + 1

repeat("go", 3)
```

## Povernennya znachennya

```piton
functia square(x):
    vernuty x * x
```

## Obsyah vydymosti

Параметри функції обчислюються в місці виклику і передаються в нове локальне оточення. Функція має власний scope, але бачить глобальні значення через зовнішнє оточення.

```piton
outer = 10

functia show():
    drukuvaty outer

show()
```

## Funktsii z moduliv

Після `vykorystaty "hello"` можна викликати функцію через селектор:

```piton
hello.Hello("svit")
```

## Rekursiya

Функції можуть викликати самі себе, тобто рекурсія підтримується.

```piton
functia fact(n):
    yaksho n <= 1:
        vernuty 1
    inackshe:
        vernuty n * fact(n - 1)

drukuvaty fact(5)
```

Це працює для задач на кшталт обходу, факторіала або рекурсивної обробки вкладених структур.
