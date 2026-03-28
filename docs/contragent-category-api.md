# Contragent Category GET API

Base URL: `http://localhost:8081/api/v1`

Bu API `contragent` uchun `admin` dagi category/subcategory GET endpointlariga o'xshash read-only endpointlarni beradi.

## Muhim

- `Authorization: Bearer <contragent_token>` majburiy.
- Faqat GET endpointlar.
- Pagination: `page` default `1`, `limit` default `10`, max `100`.
- Subcategory listda `parent_id` filter ixtiyoriy.

## 1) Kategoriyalar (contragent)

### List
`GET /contragents/me/categories?page=1&limit=10`

### Get by id
`GET /contragents/me/categories/{id}`

## 2) Subkategoriyalar (contragent)

### List
`GET /contragents/me/subcategories?page=1&limit=10`

Filter bilan:
`GET /contragents/me/subcategories?parent_id=1&page=1&limit=20`

### Get by id
`GET /contragents/me/subcategories/{id}`

## 3) Status kodlar

- `200` - Muvaffaqiyatli
- `400` - So'rov noto'g'ri (`id` yoki `parent_id` xato)
- `401` - Token yo'q yoki yaroqsiz
- `404` - Kategoriya/Subkategoriya topilmadi
- `500` - Server xatoligi
