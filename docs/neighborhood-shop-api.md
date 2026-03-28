# Neighborhood shop (maxalla do'koni) API

Base URL: `http://localhost:8081/api/v1`

Kod va jadval nomi inglizcha: **`neighborhood_shop`** / **`neighborhood_shops`**.

## Muhim

- Barcha endpointlar `Authorization: Bearer <token>` talab qiladi.
- Faqat `general` role ushbu endpointlardan foydalana oladi.
- Yangi yaratilganda **parol saqlanmaydi**. `password_setup_allowed` (default `true`) keyinroq parol o‘rnatishga ruxsat beradi.
- **Tahrirlashda** `password` ixtiyoriy yuborilsa — **hashlanib** saqlanadi; bo‘sh qoldirilsa parol o‘zgarmaydi.
- Login hozircha **yo‘q**; javobda `has_password` qaytariladi.

## Validatsiya

- `inn`: **ixtiyoriy**. Berilsa — faqat **9** yoki **12** raqam (`123456789` yoki `123456789012`). Bo‘sh yoki umuman yubormaslik mumkin.
- `phone`: `+998901234567` formatida, **unique** (faqat `neighborhood_shops` jadvali ichida).
- `region_id`, `district_id`, `mfy_id`: mavjud va **mos** bo‘lishi kerak (tuman shu viloyatga, MFY shu tumanga tegishli).
- `logo`: base64 matn (ixtiyoriy).
- `status`: `active` yoki `inactive`.

## 1) Yaratish

**POST** `/neighborhood_shops`

```json
{
  "name": "Do'kon nomi",
  "inn": "123456789",
  "region_id": 1,
  "district_id": 1,
  "mfy_id": 1,
  "phone": "+998901112233",
  "logo": "data:image/png;base64,...",
  "status": "active",
  "password_setup_allowed": true
}
```

`inn` va `logo` ixtiyoriy. `password_setup_allowed` berilmasa default `true`.

## 2) Ro‘yxat (pagination)

**GET** `/neighborhood_shops?page=1&limit=10`

- `page` — default `1`
- `limit` — default `10`, max `100`

Ro‘yxatda `logo` bo‘sh qaytariladi.

## 3) Bitta yozuv

**GET** `/neighborhood_shops/{id}`

## 4) Yangilash

**PUT** `/neighborhood_shops/{id}`

```json
{
  "name": "Yangilangan nom",
  "inn": "",
  "region_id": 1,
  "district_id": 1,
  "mfy_id": 1,
  "phone": "+998901112233",
  "logo": "",
  "status": "inactive",
  "password_setup_allowed": false,
  "password": "yangiParol123"
}
```

`password` ixtiyoriy.

## 5) Status

**PATCH** `/neighborhood_shops/{id}/status`

```json
{
  "status": "inactive"
}
```

## 6) O‘chirish

**DELETE** `/neighborhood_shops/{id}`

## 7) Status kodlar

- `200` — OK
- `201` — Yaratildi
- `400` — Validatsiya
- `401` / `403` — Auth / general emas
- `404` — Topilmadi
- `409` — Telefon band
- `500` — Server xatosi
