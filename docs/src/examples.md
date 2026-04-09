# Pryklady

У репозиторії вже є готові сценарії, які зручно використовувати як живу документацію:

- `examples/hello.piton` - мінімальний старт
- `examples/session-tracker.piton` - базовий потік програми
- `examples/spysok-demo.piton` - робота зі списками
- `examples/dodaty-demo.piton` - злиття списків і зрізи
- `examples/slovnyk-demo.piton` - словники, вкладені значення, `delete()`
- `examples/vypadkovo-demo.piton` - випадкові значення
- `examples/vykorystaty-demo.piton` - модулі та імпорти
- `examples/chas-demo.piton` - `chas()`, `zatrymka()`, округлення

## Rekomendovanyi maršrut

1. `examples/hello.piton`
2. `examples/spysok-demo.piton`
3. `examples/slovnyk-demo.piton`
4. `examples/vykorystaty-demo.piton`
5. `examples/session-tracker.piton`

## Koryst ne tilky yak demo

Ці файли зручні ще й як тестовий набір для ручної перевірки:

- запуск виконання
- генерація SVG
- перевірка імпортів
- перевірка колекцій та builtins
- швидка звірка документації з реальною поведінкою мови

Для швидкої перевірки будь-якого прикладу:

```bash
./piton examples/slovnyk-demo.piton
```
