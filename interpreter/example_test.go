package interpreter

import "log"

func ExampleRun() {
	code := `
nekhay values = [1, 4, 10, 25, 0.5, 100]
nekhay terms = ["korin", "stupin", "loh10", "abs", "dovzhyna"]
nekhay currentValue = 0

functia main():
    drukuvaty "\n--- Piton Showcase ---"
    nekhay i = 0
    poky i < dovzhyna(values):
        currentValue = values[i]
        drukuvaty "\nЧисло:"
        drukuvaty currentValue
        drukuvaty "    Корінь:"
        drukuvaty korin currentValue
        drukuvaty "    Степінь 1.5:"
        drukuvaty currentValue stupin 1.5
        yaksho currentValue > 0:
            drukuvaty "    Log10:"
            drukuvaty loh10 currentValue
        inackshe:
            drukuvaty "    Log10: неприпустимо"
        drukuvaty "    Модуль від (число - 5):"
        drukuvaty abs (currentValue - 5)
        yaksho currentValue > 0:
            nekhay wave = korin currentValue / (currentValue + 1)
            drukuvaty "    Kosynus:"
            drukuvaty kosynus currentValue
            drukuvaty "    Arksyn (wave):"
            drukuvaty arksyn wave
        inackshe:
            drukuvaty "    Kosynus/Arksyn: вимагають додаткових даних"
        drukuvaty "    Класифікація:"
        drukuvaty klasify()
        i = i + 1
    kinets

    drukuvaty "\nСписок операцій:"
    nekhay cursor = 0
    poky cursor < dovzhyna(terms):
        drukuvaty terms[cursor]
        cursor = cursor + 1
    kinets

    nekhay first = values[0]
    nekhay last = values[dovzhyna(values) - 1]
    drukuvaty "\nКраї діапазону:"
    drukuvaty first
    drukuvaty last

    drukuvaty "\nСереднє арифметичне:"
    drukuvaty average()

    drukuvaty "\nМікро-аналіз розмітки термінів:"
    drukuvaty terms[dovzhyna(terms) - 2]

    drukuvaty "\nГрошовий сигнал:"
    drukuvaty memo()

    drukuvaty "\nДемонстрація передачі аргументів:"
    formatValue(first, "Початкове значення:")

    drukuvaty "\nФункція пошуку чисел Фібоначі (рекурсія):"
    nekhay i = 0
    poky i < 10:
        drukuvaty fib(i)
        i = i + 1
    kinets

main()

functia klasify():
    yaksho currentValue > 50:
        vernuty "Надвелике"
    inackshe yaksho currentValue > 10:
        vernuty "Велике"
    inackshe:
        vernuty "Компактне"

functia average():
    nekhay sum = 0
    nekhay i = 0
    poky i < dovzhyna(values):
        sum = sum + values[i]
        i = i + 1
    kinets
    vernuty sum / dovzhyna(values)

functia memo():
    nekhay summaryIndex = dovzhyna(terms) - 2
    vernuty korin (summaryIndex stupin 2) + loh10 (dovzhyna(values)) + abs (0 - currentValue)

functia formatValue(value, label):
    drukuvaty label
    drukuvaty value

functia fib(n):
    yaksho n < 2:
        vernuty n
    inackshe:
        vernuty fib(n - 1) + fib(n - 2)
`
	err := Run(code)
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// --- Piton Showcase ---
	//
	// Число:
	// 1
	//     Корінь:
	// 1
	//     Степінь 1.5:
	// 1
	//     Log10:
	// 0
	//     Модуль від (число - 5):
	// 4
	//     Kosynus:
	// 0.5403023058681398
	//     Arksyn (wave):
	// 0.5235987755982989
	//     Класифікація:
	// Компактне
	//
	// Число:
	// 4
	//     Корінь:
	// 2
	//     Степінь 1.5:
	// 8
	//     Log10:
	// 0.6020599913279624
	//     Модуль від (число - 5):
	// 1
	//     Kosynus:
	// -0.6536436208636119
	//     Arksyn (wave):
	// 0.41151684606748806
	//     Класифікація:
	// Компактне
	//
	// Число:
	// 10
	//     Корінь:
	// 3.1622776601683795
	//     Степінь 1.5:
	// 31.622776601683796
	//     Log10:
	// 1
	//     Модуль від (число - 5):
	// 5
	//     Kosynus:
	// -0.8390715290764524
	//     Arksyn (wave):
	// 0.29159450676335213
	//     Класифікація:
	// Компактне
	//
	// Число:
	// 25
	//     Корінь:
	// 5
	//     Степінь 1.5:
	// 125
	//     Log10:
	// 1.3979400086720375
	//     Модуль від (число - 5):
	// 20
	//     Kosynus:
	// 0.9912028118634736
	//     Arksyn (wave):
	// 0.19351319251078553
	//     Класифікація:
	// Компактне
	//
	// Число:
	// 0.5
	//     Корінь:
	// 0.7071067811865476
	//     Степінь 1.5:
	// 0.3535533905932738
	//     Log10:
	// -0.3010299956639812
	//     Модуль від (число - 5):
	// 4.5
	//     Kosynus:
	// 0.8775825618903728
	//     Arksyn (wave):
	// 0.4908826782893115
	//     Класифікація:
	// Компактне
	//
	// Число:
	// 100
	//     Корінь:
	// 10
	//     Степінь 1.5:
	// 1000.0000000000002
	//     Log10:
	// 2
	//     Модуль від (число - 5):
	// 95
	//     Kosynus:
	// 0.8623188722876839
	//     Arksyn (wave):
	// 0.09917238380592071
	//     Класифікація:
	// Компактне
	//
	// Список операцій:
	// korin
	// stupin
	// loh10
	// abs
	// dovzhyna
	//
	// Краї діапазону:
	// 1
	// 100
	//
	// Середнє арифметичне:
	// 23.416666666666668
	//
	// Мікро-аналіз розмітки термінів:
	// abs
	//
	// Грошовий сигнал:
	// 3.778151250383644
	//
	// Демонстрація передачі аргументів:
	// Початкове значення:
	// 1
	//
	// Функція пошуку чисел Фібоначі (рекурсія):
	// 0
	// 1
	// 1
	// 2
	// 3
	// 5
	// 8
	// 13
	// 21
	// 34
}
