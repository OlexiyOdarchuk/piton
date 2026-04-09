# Cookbook

Нижче - короткі рецепти, які часто потрібні при роботі з мовою.

## Proitysya po spysku

```piton
i = 0
poky i < dovzhyna(items):
    drukuvaty items[i]
    i = i + 1
```

## Vydalyty element zi spysku

```piton
s = [1, 2, 3, 4]
idx = 1
s = dodaty(s[:idx], s[idx + 1:])
drukuvaty s
```

## Dodaty pole do slovnyka

```piton
user = {"name": "Piton"}
user["role"] = "admin"
```

## Vydalyty klyuch zi slovnyka

```piton
user = {"name": "Piton", "role": "admin"}
delete(user, "role")
```

## Vypadkovyi element zi spysku

```piton
options = ["red", "green", "blue"]
drukuvaty vypadkovo(options)
```

## Okruglyty rezultat

```piton
pi = 22 / 7
drukuvaty zaokruhlennya(pi, 3)
```

## Rekursyvnyi factorial

```piton
functia fact(n):
    yaksho n <= 1:
        vernuty 1
    inackshe:
        vernuty n * fact(n - 1)

drukuvaty fact(5)
```

## Onovyty vkladene znachennya

```piton
profile = {"settings": {"theme": "light"}}
profile["settings"]["theme"] = "dark"
```
