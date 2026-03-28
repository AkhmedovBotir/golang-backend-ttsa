# Admin Marketplace User API

Marketplace userlarni admin boshqaruvi uchun API.

Base URL: `http://localhost:8081/api/v1`

`Authorization: Bearer <admin_token>` majburiy, `general` role kerak.

Muhim:
- `CREATE` endpoint yo'q.
- Ya'ni admin faqat mavjud marketplace userlarni ko'radi, yangilaydi, statusini o'zgartiradi va o'chiradi.

## 1) List users

`GET /marketplace-users?page=1&limit=10`

Ixtiyoriy filterlar:
- `status` (`active`/`inactive`)
- `region_id`
- `district_id`
- `mfy_id`
- `phone`
- `q` (ism yoki familiya bo'yicha qidiruv)

## 2) Get by ID

`GET /marketplace-users/{id}`

## 3) Update user

`PUT /marketplace-users/{id}`

```json
{
  "first_name": "Ali",
  "last_name": "Valiyev",
  "gender": "erkak",
  "phone": "+998901234567",
  "region_id": 1,
  "district_id": 10,
  "mfy_id": 265,
  "birth_date": "1995-03-21",
  "status": "active"
}
```

## 4) Update status

`PATCH /marketplace-users/{id}/status`

```json
{
  "status": "inactive"
}
```

## 5) Delete user

`DELETE /marketplace-users/{id}`

## Status kodlar

- `200` - muvaffaqiyatli
- `400` - so'rov noto'g'ri
- `401` - token yo'q yoki yaroqsiz
- `403` - general admin emas
- `404` - user topilmadi
- `409` - telefon raqami allaqachon mavjud
- `500` - server xatoligi

