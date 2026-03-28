# Admin CRUD API hujjati

Base URL: `http://localhost:8081/api/v1`

## 1) Login

**POST** `/admins/login`

Body:
```json
{
  "username": "general_admin",
  "password": "12345678"
}
```

Muvaffaqiyatli javob (`200`):
```json
{
  "message": "Muvaffaqiyatli login qilindi",
  "data": {
    "token": "jwt_token",
    "admin": {
      "id": 1,
      "name": "Super Admin",
      "role": "general",
      "phone": "+998901234567",
      "username": "general_admin",
      "status": "active"
    }
  }
}
```

Xatoliklar:
- `400` - So'rov formati noto'g'ri
- `401` - Username yoki parol noto'g'ri
- `403` - Admin faol emas
- `500` - Serverda xatolik

## 2) Admin yaratish (faqat general)

**POST** `/admins`

Header:
- `Authorization: Bearer <token>`

Body:
```json
{
  "name": "Ali Valiyev",
  "role": "admin",
  "phone": "+998901112233",
  "username": "ali_admin",
  "password": "12345678",
  "status": "active"
}
```

Xatoliklar:
- `400` - Validation xatolari
- `401` - Token yo'q/yaroqsiz
- `403` - Faqat general adminga ruxsat
- `409` - `phone` yoki `username` band
- `500` - Server xatosi

## 3) Barcha adminlar ro'yxati (faqat general)

**GET** `/admins`

Header:
- `Authorization: Bearer <token>`

Query params:
- `page` (ixtiyoriy, default: `1`)
- `limit` (ixtiyoriy, default: `10`, max: `100`)

Javob:
- `200` - Adminlar ro'yxati
- `401`, `403`, `500`

Namuna javob:
```json
{
  "message": "Adminlar ro'yxati olindi",
  "data": {
    "items": [
      {
        "id": 1,
        "name": "Super Admin",
        "role": "general",
        "phone": "+998901234567",
        "username": "general_admin",
        "status": "active"
      }
    ],
    "total": 25,
    "page": 1,
    "limit": 10,
    "total_pages": 3
  }
}
```

## 4) Bitta adminni olish (faqat general)

**GET** `/admins/{id}`

Javob:
- `200` - Admin topildi
- `400` - ID noto'g'ri
- `404` - Admin topilmadi
- `401`, `403`, `500`

## 5) Admin yangilash (faqat general)

**PUT** `/admins/{id}`

Body:
```json
{
  "name": "Ali Valiyev",
  "role": "admin",
  "phone": "+998901112233",
  "username": "ali_admin",
  "password": "newpassword123",
  "status": "inactive"
}
```

Eslatma:
- `password` bo'sh yuborilsa eski parol saqlanadi.

Xatoliklar:
- `400`, `401`, `403`, `404`, `409`, `500`

## 6) Admin o'chirish (faqat general)

**DELETE** `/admins/{id}`

Javob:
- `200` - O'chirildi
- `400`, `401`, `403`, `404`, `500`

## 7) Admin statusini o'zgartirish (faqat general)

**PATCH** `/admins/{id}/status`

Body:
```json
{
  "status": "inactive"
}
```

`status` qiymati:
- `active`
- `inactive`

Javob:
- `200` - Status yangilandi
- `400`, `401`, `403`, `404`, `500`

## Umumiy response formati

Har doim quyidagi format qaytadi:
```json
{
  "message": "o'zbekcha xabar",
  "data": {},
  "error": "xatolik tafsiloti"
}
```

`data` yoki `error` holatga qarab bo'sh bo'lishi mumkin.
