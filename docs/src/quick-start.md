# Shvydkyi Start

Найкоротший шлях до першого запуску:

1. Скачати бінарник з [Releases](https://github.com/OlexiyOdarchuk/piton/releases)
2. Запустити REPL:

```bash
./piton
```

Або запустити відразу готовий приклад:

```bash
./piton examples/slovnyk-demo.piton
```

Якщо хочеш саме зібрати інтерпретатор самостійно з вихідного коду, дивись [Vstanovlennya ta CLI](./install-and-cli.md).

## Naiprostisha programa

```piton
drukuvaty "Pryvit, svite!"
```

Запуск:

```bash
./piton hello.piton
```

## Shcho vazhlyvo z pershykh khvylyn

- відступи мають значення
- змінні створюються при першому присвоєнні
- числа всередині мови працюють як числа з плаваючою крапкою
- для списків і словників є індексація через `[]`
- якщо запустити `./piton` без файла, відкриється REPL

## Pryklad zi zminnymy, umovoyu i tsyklom

```piton
values = [1, 2, 3]
i = 0

poky i < dovzhyna(values):
    yaksho values[i] > 1:
        drukuvaty "bilshe za odyn"
    inackshe:
        drukuvaty "odyn abo menshe"
    i = i + 1
```

## Hotovi skrypty v repozytorii

- `examples/session-tracker.piton`
- `examples/spysok-demo.piton`
- `examples/slovnyk-demo.piton`
- `examples/vykorystaty-demo.piton`
- `examples/chas-demo.piton`

Найкращий наступний крок - [Tur po Movi](./tour.md).
