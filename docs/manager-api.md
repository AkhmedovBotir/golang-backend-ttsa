# Manager API

Base URL: `http://localhost:8081/api/v1`

Modul fayllari: `modules/admin` ichida `manager_` prefiks (`manager.go`, `manager_repository.go`, `manager_service.go`, `manager_http.go`).

## Muhim

- Barcha endpointlar `Authorization: Bearer <token>` talab qiladi.
- Faqat `general` role ushbu endpointlardan foydalana oladi.
- **Yaratishda** parol bazaga yozilmaydi (`password` e'tiborsiz). `password_setup_allowed` (default `true`) keyinroq parol qo‘yishga ruxsat beradi.
- **Tahrirlashda** `password` ixtiyoriy yuborilsa — **bcrypt** bilan hashlanib saqlanadi; bo‘sh qoldirilsa parol o‘zgarmaydi.
- Login hozircha **alohida endpoint sifatida yo‘q**; javobda `has_password` bor.

## Maydonlar (JSON)

| Maydon | Tavsif |
|--------|--------|
| `name` | Menejer nomi |
| `phone` | `+998901234567` formatida, **unique** (`managers` ichida) |
| `viloyat_id` | Viloyat (`regions` jadvalidagi `id`) |
| `status` | `active` yoki `inactive` |
| `password_setup_allowed` | Ixtiyoriy, default `true` |
| `password` | Faqat **PUT** da ixtiyoriy |

## 1) Yaratish

**POST** `/managers`

```json
{
  "name": "Menejer 1",
  "phone": "+998901112233",
  "viloyat_id": 1,
  "status": "active",
  "password_setup_allowed": true
}
```

## 2) Ro‘yxat (pagination)

**GET** `/managers?page=1&limit=10`

- `page` — default `1`
- `limit` — default `10`, max `100`

## 3) Bitta menejer

**GET** `/managers/{id}`

## 4) Yangilash

**PUT** `/managers/{id}`

```json
{
  "name": "Menejer 1",
  "phone": "+998901112233",
  "viloyat_id": 1,
  "status": "inactive",
  "password_setup_allowed": false,
  "password": "yangiParol123"
}
```

## 5) Status

**PATCH** `/managers/{id}/status`

```json
{
  "status": "inactive"
}
```

## 6) O‘chirish

**DELETE** `/managers/{id}`

## 7) Status kodlar

- `200` — OK
- `201` — Yaratildi
- `400` — Validatsiya
- `401` / `403` — Auth / general emas
- `404` — Topilmadi
- `409` — Telefon band
- `500` — Server xatosi
