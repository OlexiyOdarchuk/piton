# Vizualizator

Piton уміє не тільки виконувати код, а й будувати SVG-блок-схеми за **ДСТУ**.

Рендером займається бібліотека [rombik](https://rombik.ishawyha.dev): Piton зводить
AST програми до її проміжного представлення, а rombik сам робить розкладку,
маршрутизацію стрілок і малює чистий SVG (без зовнішніх залежностей).

## Komandy

```bash
./piton -draw ./main.piton
./piton -all ./main.piton
./piton -draw -split ./main.piton
./piton -draw -target=myFunc ./main.piton
```

## Shcho robyt

- збирає AST програми
- може підтягувати імпортовані модулі (`vykorystaty`)
- будує **окрему схему на кожну функцію**; у загальному режимі вони зводяться в
  один файл одна під одною, з підписом над кожною

Це корисно для навчання, рев’ю логіки та демонстрації алгоритмів.

## Fihury (ДСТУ)

Кожна конструкція Piton має свою фігуру:

| Фігура | Що означає | Звідки |
|---|---|---|
| овал «Початок/Кінець» | термінатор | межі функції |
| прямокутник | дія | присвоєння, вираз |
| ромб | розв’язок «Так/Ні» | `yaksho` / `inackshe` |
| паралелограм | ввід / вивід | `vvid`, `drukuvaty` |
| прямокутник із боковими рисками | підпрограма | виклик визначеної у файлі функції |

Цикл `poky` малюється ромбом-умовою з дугою повернення.

## Pryklad: povna skhema prohramy

Нижче - згенерована схема для `examples/session-tracker.piton`. У загальному
режимі головна програма і всі функції йдуть одним файлом, кожна під своїм підписом.

![Session tracker flowchart](./assets/session-tracker.svg)

## Pryklad: skhema proyektu z importamy

Нижче - згенерована схема для `examples/vykorystaty-demo.piton`, побудована через `-all`,
тобто разом з імпортованим модулем. Функції з імпортів підписані іменем модуля
(напр. «Модуль hello — функція Hello»).

![Imports project flowchart](./assets/imports-project.svg)

## Pryklad: okrema funktsiya u split-rezhymi

`-split` корисний, коли повна схема завелика: кожна функція потрапляє в окремий
файл `{модуль}_{функція}.svg`. Ось функція `average` із `session-tracker`:

![Average function flowchart](./assets/visualizer/main_average.svg)

І ще одна окрема функція `deviation`:

![Deviation function flowchart](./assets/visualizer/main_deviation.svg)

## Koly shcho vykorystovuvaty

- `-draw` - коли працюєш з одним файлом
- `-all` - коли хочеш бачити програму разом з імпортами
- `-split` - коли функцій багато і потрібні окремі діаграми
- `-target=name` - коли треба проаналізувати тільки одну функцію

SVG, які бачиш у цій книжці, згенеровані реальним `piton` з прикладів репозиторію, а не намальовані вручну.
