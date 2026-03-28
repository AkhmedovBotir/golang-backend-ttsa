# Punkt API

Base URL: `http://localhost:8081/api/v1`

Modul fayllari: `modules/admin` ichida `punkt_` prefiks (`punkt.go`, `punkt_repository.go`, `punkt_service.go`, `punkt_http.go`).

**Agent** API dan farqi: punkt faqat **viloyat** va **tuman** bilan bog‘lanadi; **`mfy_id` yo‘q**.

## Muhim

- Barcha endpointlar `Authorization: Bearer <token>` talab qiladi.
- Faqat `general` role ushbu endpointlardan foydalana oladi.
- **Yaratishda** parol bazaga yozilmaydi (`password` e’tiborsiz). `password_setup_allowed` (default `true`) keyinroq parol qo‘yishga ruxsat beradi.
- **Tahrirlashda** `password` ixtiyoriy yuborilsa — **bcrypt** bilan hashlanib saqlanadi; bo‘sh qoldirilsa parol o‘zgarmaydi.
- Login hozircha **alohida endpoint sifatida yo‘q**; javobda `has_password` bor.

## So‘rov maydonlari (JSON)

| Maydon | Tavsif |
|--------|--------|
| `name` | Punkt nomi |
| `viloyat_id` | Viloyat (`regions` jadvalidagi `id`) |
| `tuman_id` | Tuman (`districts` jadvalidagi `id`) |
| `phone` | `+998901234567` formatida, **unique** (`punkts` ichida) |
| `status` | `active` yoki `inactive` |
| `password_setup_allowed` | Ixtiyoriy, default `true` |
| `password` | Faqat **PUT** da ixtiyoriy |

Viloyat–tuman mosligi tekshiriladi: tuman shu viloyatga tegishli bo‘lishi kerak.

## Javob obyekti (GET ro‘yxat elementi / GET bitta / POST / PUT)

Parol qaytarilmaydi. Qolgan maydonlar:

| Maydon | Tavsif |
|--------|--------|
| `id` | Punkt `id` |
| `name` | Nom |
| `viloyat_id` | Viloyat `id` |
| `tuman_id` | Tuman `id` |
| `phone` | Telefon |
| `status` | `active` / `inactive` |
| `password_setup_allowed` | Parol o‘rnatishga ruxsat |
| `has_password` | Bazada parol hash bor-yo‘qligi |
| `created_at`, `updated_at` | Vaqt belgilari |

**Eslatma:** javobda `mfy_id` **bo‘lmaydi**.

## 1) Yaratish

**POST** `/punkts`

```json
{
  "name": "1-son punkt",
  "viloyat_id": 1,
  "tuman_id": 1,
  "phone": "+998901112233",
  "status": "active",
  "password_setup_allowed": true
}
```

## 2) Ro‘yxat (pagination)

**GET** `/punkts?page=1&limit=10`

- `page` — default `1`
- `limit` — default `10`, max `100`

## 3) Bitta punkt

**GET** `/punkts/{id}`

## 4) Yangilash

**PUT** `/punkts/{id}`

```json
{
  "name": "1-son punkt",
  "viloyat_id": 1,
  "tuman_id": 1,
  "phone": "+998901112233",
  "status": "active",
  "password_setup_allowed": false,
  "password": "yangiParol123"
}
```

## 5) Status

**PATCH** `/punkts/{id}/status`

```json
{
  "status": "inactive"
}
```

## 6) O‘chirish

**DELETE** `/punkts/{id}`

## 7) Ma'lumotlar bazasi

Jadval: **`punkts`**. Ustunlar: `name`, `region_id` (API: `viloyat_id`), `district_id` (API: `tuman_id`), `phone`, `status`, `password`, `password_setup_allowed`, vaqt maydonlari. **`mfy_id` jadvalda bo‘lmasligi kerak.**

Agar eski versiyadan `mfy_id` qolgan bo‘lsa (PostgreSQL):

```sql
ALTER TABLE punkts DROP COLUMN IF EXISTS mfy_id;
```

## 8) Status kodlar

- `200` — OK
- `201` — Yaratildi
- `400` — Validatsiya
- `401` / `403` — Auth / general emas
- `404` — Topilmadi
- `409` — Telefon band
- `500` — Server xatosi
