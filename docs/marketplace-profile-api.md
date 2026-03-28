# Marketplace Profile + Avatar API

Marketplace foydalanuvchining o'z profilini boshqarish API.

Base URL: `http://localhost:8081/api/v1`

`Authorization: Bearer <marketplace_token>` majburiy.

## 1) O'z profilini olish

- **Method:** `GET`
- **URL:** `/marketplace/me/profile`

## 2) O'z profilini yangilash

- **Method:** `PUT`
- **URL:** `/marketplace/me/profile`
- **Request:**

```json
{
  "first_name": "Ali",
  "last_name": "Valiyev",
  "gender": "erkak",
  "region_id": 1,
  "district_id": 10,
  "mfy_id": 265,
  "birth_date": "1995-03-21"
}
```

Yangilanadigan maydonlar:
- `first_name`
- `last_name`
- `gender`
- `birth_date`
- `region_id`
- `district_id`
- `mfy_id`

## 3) Avatarni ko'rish

- **Method:** `GET`
- **URL:** `/marketplace/me/avatar`

Javob:
- `has_avatar` (`true/false`)
- `avatar` (base64 string, bo'sh bo'lishi mumkin)

## 4) Avatar yuklash / o'zgartirish

- **Method:** `PUT`
- **URL:** `/marketplace/me/avatar`
- **Request:**

```json
{
  "avatar": "data:image/png;base64,iVBORw0KGgoAAA..."
}
```

Eslatma:
- oddiy base64 ham qabul qilinadi;
- data URL format (`data:image/...;base64,`) ham qabul qilinadi.

## 5) Avatarni o'chirish

- **Method:** `DELETE`
- **URL:** `/marketplace/me/avatar`

## Status kodlar

- `200` - muvaffaqiyatli
- `400` - validation xato (`gender`, `birth_date`, location yoki base64 noto'g'ri)
- `401` - token yo'q/yaroqsiz
- `404` - user topilmadi
- `500` - server xatoligi

