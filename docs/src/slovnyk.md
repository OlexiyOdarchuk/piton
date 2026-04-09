# Slovnyk

Словник задається літералом через `{}`. Ключі мають бути рядками.

```piton
profil = {"imya": "Mavka", "vik": 128}
drukuvaty profil
```

## Chytannya

```piton
drukuvaty profil["imya"]
```

Якщо ключа немає, інтерпретатор надрукує помилку.

## Onovlennya i dodavannya

```piton
profil["vik"] = 129
profil["rol"] = "moderator"
```

## Vydalennya

Для видалення ключа використовується `delete(slovnyk, "key")`.

```piton
profil = {"imya": "Mavka", "vik": 128, "rol": "moderator"}
delete(profil, "rol")
drukuvaty profil
```

## Dovzhyna

```piton
drukuvaty dovzhyna(profil)
```

## Typovyi pryklad

```piton
profil = {"imya": "Mavka", "vik": 128, "mista": ["Lviv", "Kyiv"]}
profil["rol"] = "moderator"
profil["mista"][0] = "Odesa"
delete(profil, "vik")
drukuvaty profil
```

## Vkladenist

```piton
profil = {
    "imya": "Mavka",
    "mista": ["Lviv", "Kyiv"],
    "nalashtuvannya": {"tema": "temna"}
}

drukuvaty profil["nalashtuvannya"]["tema"]
```
