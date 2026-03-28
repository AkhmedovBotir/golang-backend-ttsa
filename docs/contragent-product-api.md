# Contragent Product CRUD API

Base URL: `http://localhost:8081/api/v1`

`Authorization: Bearer <contragent_token>` majburiy.

## Muhim qoidalar

- `description` delta format string qabul qiladi.
- `images` minimum `1`, maksimum `5`, base64 format.
- `unit` faqat: `dona`, `litr`, `kg`.
- `kpi_bonus_percent`: `0..100`.
- `product_code` tizim tomonidan global ketma-ketlikda avtomatik beriladi.
- `moderation_status` yaratishda har doim `pending`.
- `rejection_reason` yaratishda bo'sh bo'ladi.
- `subcategory_id` tanlangan `category_id`ga tegishli bo'lishi shart.

## 1) Create product

`POST /contragents/me/products`

```json
{
  "name": "Suv 1L",
  "description": "{\"ops\":[{\"insert\":\"Mahsulot tavsifi\"}]}",
  "price": 12000,
  "original_price": 10000,
  "images": [
    "data:image/png;base64,iVBORw0KGgoAAA...",
    "data:image/png;base64,iVBORw0KGgoBBB..."
  ],
  "category_id": 1,
  "subcategory_id": 10,
  "quantity": 250,
  "unit": "litr",
  "unit_size": "1 litr",
  "status": "active",
  "kpi_bonus_percent": 20
}
```

## 2) List products

`GET /contragents/me/products?page=1&limit=10`

## 3) Get by id

`GET /contragents/me/products/{id}`

## 4) Update product

`PUT /contragents/me/products/{id}`

Yaratish payloadi bilan bir xil format.

Eslatma: update qilinganda mahsulot qayta moderatsiyaga tushadi (`moderation_status = pending`).

## 5) Update status

`PATCH /contragents/me/products/{id}/status`

```json
{
  "status": "inactive"
}
```

## 6) Delete product

`DELETE /contragents/me/products/{id}`

## 7) Javobdagi qo'shimcha maydonlar

- `product_code` (global ketma-ket kod)
- `moderation_status` (`pending/approved/rejected`)
- `rejection_reason`
- `kpi_bonus_amount` = `(price - original_price) * kpi_bonus_percent / 100` (manfiy bo'lmaydi)

## 8) Status kodlar

- `200` - Muvaffaqiyatli
- `201` - Yaratildi
- `400` - So'rov noto'g'ri (validation)
- `401` - Token yo'q yoki yaroqsiz
- `404` - Mahsulot topilmadi
- `500` - Server xatoligi
