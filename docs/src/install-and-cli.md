# Vstanovlennya ta CLI

## Hotovi binarnyky z Releases

Найпростіший варіант - скачати готовий бінарник з [Releases](https://github.com/OlexiyOdarchuk/piton/releases) репозиторію.

У релізах публікуються збірки для:

- `Linux`
- `macOS`
- `Windows`
- архітектур `amd64` та `arm64`

Після завантаження достатньо розпакувати файл і запускати `piton` напряму.

## Zbirka z vykhidnoho kodu

```bash
git clone https://github.com/OlexiyOdarchuk/piton.git
cd piton
go build -o piton ./cmd/piton
```

## Zapusk skrypta

```bash
./piton script.piton
```

## Osnovni rezhymy CLI

- `./piton file.piton` - виконання програми
- `./piton -draw file.piton` - блок-схема для одного файла
- `./piton -all file.piton` - блок-схема для проєкту з імпортами
- `./piton -draw -split file.piton` - окремий SVG для кожної функції
- `./piton -draw -target=myFunc file.piton` - схема лише для однієї функції

## Korysni detali pro REPL

- вхід у REPL: `./piton`
- вихід з REPL: `exit`
- багаторядкові блоки після рядка з `:` підтримуються

## Zapusk prykladiv

```bash
./piton examples/slovnyk-demo.piton
./piton examples/spysok-demo.piton
./piton examples/vypadkovo-demo.piton
```

CLI зараз орієнтований на дві задачі:

- виконання `.piton` файлів
- генерація SVG-блок-схем
