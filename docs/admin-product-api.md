# Admin Product CRUD + Moderation API

Base URL: `http://localhost:8081/api/v1`

`Authorization: Bearer <admin_token>` majburiy, `general` role kerak.

## 1) CRUD endpointlar

### Create
`POST /products`

```json
{
  "contragent_id": 1,
  "name": "Suv 1L",
  "description": "{\"ops\":[{\"insert\":\"Tavsif\"}]}",
  "price": 12000,
  "original_price": 10000,
  "images": ["data:image/png;base64,..."],
  "category_id": 1,
  "subcategory_id": 10,
  "quantity": 100,
  "unit": "litr",
  "unit_size": "1 litr",
  "status": "active",
  "kpi_bonus_percent": 20
}
```

Yaratishda `moderation_status` avtomatik `pending`.

### List
`GET /products?page=1&limit=10`

Ixtiyoriy filter:
- `contragent_id`
- `moderation_status` (`pending`, `approved`, `rejected`)

### Get by id
`GET /products/{id}`

### Update
`PUT /products/{id}`

Create payload bilan bir xil format.

### Status update
`PATCH /products/{id}/status`

```json
{
  "status": "inactive"
}
```

### Delete
`DELETE /products/{id}`

## 2) Moderatsiya endpointlar

### Tasdiqlash
`PATCH /products/{id}/approve`

Natija: `moderation_status = approved`, `rejection_reason = ""`

### Bekor qilish (rad etish)
`PATCH /products/{id}/reject`

```json
{
  "rejection_reason": "Noto'g'ri rasm va tavsif"
}
```

Natija: `moderation_status = rejected`

## 3) Qoidalar

- `images`: min `1`, max `5`, base64.
- `unit`: `dona`, `litr`, `kg`.
- `kpi_bonus_percent`: `0..100`.
- `subcategory_id` tanlangan `category_id`ga tegishli bo'lishi shart.
- `product_code` global ketma-ketlikda avtomatik beriladi.

## 4) Status kodlar

- `200` - Muvaffaqiyatli
- `201` - Yaratildi
- `400` - So'rov noto'g'ri (validation)
- `401` - Token yo'q yoki yaroqsiz
- `403` - Faqat general admin
- `404` - Topilmadi
- `500` - Server xatoligi
