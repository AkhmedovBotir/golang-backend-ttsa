# Category va Subcategory API

Base URL: `http://localhost:8081/api/v1`

Manba JSON: `scripts/ttsa.categories.json` (parent `null` => category, parent bor => subcategory).

## 0) Import (Mongo JSON -> PostgreSQL)

PowerShell:
```powershell
powershell -ExecutionPolicy Bypass -File ./scripts/import-categories.ps1
```

Custom path bilan:
```powershell
powershell -ExecutionPolicy Bypass -File ./scripts/import-categories.ps1 -JsonPath "C:\path\to\ttsa.categories.json"
```

To'g'ridan-to'g'ri:
```powershell
go run ./cmd/import-categories scripts/ttsa.categories.json
```

Import logikasi:
- 1-pass: `parent=null` bo'lgan category lar.
- 2-pass: `parent` bor bo'lgan subcategory lar.
- `external_id` bo'yicha upsert (mavjud bo'lsa update, bo'lmasa create).

## Muhim

- Barcha endpointlar `Authorization: Bearer <token>` talab qiladi.
- Faqat `general` role ushbu endpointlardan foydalana oladi.
- `status` faqat: `active`, `inactive`.
- `slug` global unique (`categories` jadvali bo'yicha).
- `image` base64 string bo'lishi mumkin (`data:image/...;base64,...` yoki oddiy base64).

Model maydonlari:
- `name`
- `slug`
- `image`
- `censored` (`true/false`)
- `status`
- `parent_id` (`null` bo'lsa category, qiymat bo'lsa subcategory)

## 1) Category API (alohida)

### Create category
`POST /categories`

```json
{
  "name": "Turizm va turizm xizmatlari",
  "slug": "turizm-va-turizm-xizmatlari",
  "image": "data:image/png;base64,...",
  "censored": false,
  "status": "active"
}
```

### List categories
`GET /categories?page=1&limit=10`

Faqat `parent_id = null` lar qaytadi.

### Get category by id
`GET /categories/{id}`

### Update category
`PUT /categories/{id}`

```json
{
  "name": "Turizm",
  "slug": "turizm",
  "image": "",
  "censored": false,
  "status": "inactive"
}
```

### Category status update
`PATCH /categories/{id}/status`

```json
{
  "status": "active"
}
```

### Delete category
`DELETE /categories/{id}`

## 2) Subcategory API (alohida)

### Create subcategory
`POST /subcategories`

```json
{
  "name": "Turistik agentlar",
  "slug": "turistik-agentlar",
  "image": "",
  "censored": false,
  "parent_id": 1,
  "status": "active"
}
```

`parent_id` category (`parent_id=null`) bo'lgan yozuvga tegishli bo'lishi shart.

### List subcategories
`GET /subcategories?page=1&limit=10`

Ixtiyoriy filter:
- `parent_id` (masalan `GET /subcategories?parent_id=1&page=1&limit=20`)

### Get subcategory by id
`GET /subcategories/{id}`

### Update subcategory
`PUT /subcategories/{id}`

```json
{
  "name": "Aviabilet sotuvchilar",
  "slug": "aviabilet-sotuvchilar",
  "image": "",
  "censored": false,
  "parent_id": 1,
  "status": "active"
}
```

### Subcategory status update
`PATCH /subcategories/{id}/status`

```json
{
  "status": "inactive"
}
```

### Delete subcategory
`DELETE /subcategories/{id}`

## 3) Status kodlar

- `200` - Muvaffaqiyatli
- `201` - Yaratildi
- `400` - So'rov noto'g'ri (slug, parent_id, status, required field)
- `401` - Token yo'q yoki yaroqsiz
- `403` - Ruxsat yo'q (general emas)
- `404` - Topilmadi
- `500` - Server xatoligi
