# Piton z Go

Piton можна використовувати як Go-бібліотеку.

## Vykonannya kodu

```go
package main

import (
	"fmt"
	"os"

	"github.com/OlexiyOdarchuk/piton/pkg/interpreter"
)

func main() {
	code := `
x = 10
y = 20
drukuvaty x + y
`

	if err := interpreter.Run(code, os.Stdout); err != nil {
		fmt.Println("Pomylka vikonannya:", err)
	}
}
```

## Generatsiya SVG z Go

```go
package main

import (
	"fmt"
	"os"

	"github.com/OlexiyOdarchuk/piton/pkg/visualizer"
)

func main() {
	code := `
functia main():
    drukuvaty "hello"

main()
`

	images, err := visualizer.Visualize(code, "", false)
	if err != nil {
		panic(err)
	}

	for name, svg := range images {
		_ = os.WriteFile(name, svg, 0644)
		fmt.Println("saved:", name)
	}
}
```

Це зручно для:

- IDE-утиліт
- навчальних платформ
- генерації схем у CI
- вбудовування Piton як DSL у Go-проєкт
