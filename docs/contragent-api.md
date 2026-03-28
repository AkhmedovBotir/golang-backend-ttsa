# Kontragent API

Base URL: `http://localhost:8081/api/v1`

## Muhim

- Barcha endpointlar `Authorization: Bearer <token>` talab qiladi.
- Faqat `general` role ushbu endpointlardan foydalana oladi.
- Yangi yaratilganda **parol saqlanmaydi** (`password` maydoni e'tiborsiz qoldiriladi). `password_setup_allowed` bilan keyinroq parol o‘rnatishga ruxsat beriladi.
- **Tahrirlashda** `password` ixtiyoriy yuborilsa, u **hashlanib** saqlanadi (bo‘sh qoldirilsa eski parol o‘zgarmaydi).
- Login hozircha **mavjud emas**; javobda `has_password` (parol bor-yo‘qligi) qaytariladi.

## Ma'lumotlar bazasi

- Jadval nomi: **`contragents`**.
- Agar ilgari `district_contragents` jadvalida ma'lumot bo'lsa, migratsiya yoki `ALTER TABLE ... RENAME TO contragents` qiling.

## Validatsiya

- `inn`: faqat **9** yoki **12** raqam (`123456789` yoki `123456789012`).
- `phone`: `+998901234567` formatida, **unique**.
- `region_id`, `district_id`, `mfy_id`: mavjud bo‘lishi va **mos kelishi** kerak (tuman shu viloyatga, MFY shu tumanga tegishli).
- `activity_type_id`: mavjud `contragent_type` yozuvining `id` si.
- `logo`: base64 matn (ixtiyoriy), juda katta bo‘lishi mumkin.
- `status`: `active` yoki `inactive`.

## 1) Yaratish

**POST** `/contragents`

```json
{
  "name": "MCHJ Namuna",
  "inn": "123456789",
  "region_id": 1,
  "district_id": 1,
  "mfy_id": 1,
  "phone": "+998901112233",
  "logo": "data:image/png;base64,...",
  "activity_type_id": 1,
  "status": "active",
  "password_setup_allowed": true
}
```

`password_setup_allowed` ixtiyoriy; berilmasa default `true`.

## 2) Ro‘yxat (pagination)

**GET** `/contragents?page=1&limit=10`

- `page` — default `1`
- `limit` — default `10`, max `100`

Ro‘yxatda `logo` bo‘sh qaytariladi (hajm kamaytirish uchun). To‘liq logo uchun bitta yozuvni oling.

## 3) Bitta yozuv

**GET** `/contragents/{id}`

## 4) Yangilash

**PUT** `/contragents/{id}`

```json
{
  "name": "MCHJ Namuna 2",
  "inn": "123456789012",
  "region_id": 1,
  "district_id": 1,
  "mfy_id": 1,
  "phone": "+998901112233",
  "logo": "",
  "activity_type_id": 2,
  "status": "active",
  "password_setup_allowed": false,
  "password": "yangiParol123"
}
```

`password` ixtiyoriy: faqat to‘ldirilganda yangilanadi.

## 5) Status

**PATCH** `/contragents/{id}/status`

```json
{
  "status": "inactive"
}
```

## 6) O‘chirish

**DELETE** `/contragents/{id}`

## 7) Javob maydonlari (qisqacha)

| Maydon | Tavsif |
|--------|--------|
| `has_password` | Parol bazada bor-yo‘qligi (login keyinroq) |
| `password_setup_allowed` | Keyin parol o‘rnatishga ruxsat |

## 8) Status kodlar

- `200` — OK
- `201` — Yaratildi
- `400` — Validatsiya xatosi
- `401` / `403` — Auth / general emas
- `404` — Topilmadi
- `409` — Telefon band
- `500` — Server xatosi
