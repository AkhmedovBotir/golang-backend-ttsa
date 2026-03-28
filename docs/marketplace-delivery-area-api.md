# Marketplace Delivery Area API

Marketplace foydalanuvchisi o'zi uchun yetkazib berish hududlarini saqlashi uchun API.

Base URL: `/api/v1`

`Authorization: Bearer <marketplace_token>` majburiy.

## Model

Har bir hudud quyidagi maydonlardan iborat:
- `name` - hudud nomi (masalan: "Uy atrofi", "Ish joyi hududi")
- `region_id` - viloyat ID
- `district_id` - tuman ID
- `mfy_id` - MFY ID
- `is_default` - asosiy manzil (`true`/`false`)

Qoidalar:
- `name` bo'sh bo'lmasligi kerak
- `region_id`, `district_id`, `mfy_id` majburiy
- tanlangan `district_id` shu `region_id` ga tegishli bo'lishi kerak
- tanlangan `mfy_id` shu `district_id` ga tegishli bo'lishi kerak
- bitta foydalanuvchida faqat bitta `is_default=true` bo'ladi
- birinchi yaratilgan hudud avtomatik asosiy bo'ladi

## Endpointlar

### 1) Yangi hudud qo'shish

`POST /marketplace/me/delivery-areas`

### 2) O'zimning hududlar ro'yxati

`GET /marketplace/me/delivery-areas`

### 3) Mavjud hududni yangilash

`PUT /marketplace/me/delivery-areas/{id}`

### 4) Mavjud hududni o'chirish

`DELETE /marketplace/me/delivery-areas/{id}`

### 5) Asosiy manzil qilib belgilash

`PATCH /marketplace/me/delivery-areas/{id}/set-default`

## Foydalanish tartibi

1. Avval `marketplace` region API orqali viloyat/tuman/mfy tanlovlarini yuklang.
2. Foydalanuvchi hudud nomini yozsin va kerakli uchta ID ni tanlasin.
3. `POST` bilan saqlang.
4. Kerak bo'lsa `PATCH .../set-default` bilan asosiy manzilni almashtiring.
5. Keyin ro'yxatdan `PUT` bilan tahrir qiling yoki `DELETE` bilan o'chiring.

## Status kodlar

- `200` - muvaffaqiyatli (`GET`, `PUT`, `DELETE`)
- `201` - muvaffaqiyatli yaratildi (`POST`)
- `400` - noto'g'ri so'rov yoki hierarchy xato
- `401` - token yo'q/yaroqsiz
- `404` - hudud topilmadi
- `500` - server xatoligi
