# Contragent Region + Delivery Area API

Base URL: `http://localhost:8081/api/v1`

`Authorization: Bearer <contragent_token>` majburiy.

## 1) Region GET API (3 ta alohida)

### Viloyatlar
`GET /contragents/me/regions`

### Tumanlar
`GET /contragents/me/districts`

Ixtiyoriy filter:
`GET /contragents/me/districts?region_id=1`

### MFYlar
`GET /contragents/me/mfys`

Ixtiyoriy filter:
`GET /contragents/me/mfys?district_id=10`

## 2) Yetkazib berish hududlari (save/update)

Bir nechta viloyat va bir nechta tuman tanlanadi, saqlanadi va keyin qayta yangilanadi.

### Save/Update delivery areas
`PUT /contragents/me/delivery-areas`

```json
{
  "region_ids": [1, 3, 5],
  "district_ids": [11, 12, 31]
}
```

Qoidalar:
- `region_ids` bo'sh bo'lmasligi kerak
- `district_ids` bo'sh bo'lmasligi kerak
- har bir `district_id` tanlangan `region_ids` ichidagi viloyatlardan biriga tegishli bo'lishi shart
- bu endpoint eski tanlovni o'chirib, yangisini to'liq saqlaydi (replace)

### Get saved delivery areas
`GET /contragents/me/delivery-areas`

Namuna javob:
```json
{
  "message": "Yetkazib berish hududlari olindi",
  "data": {
    "region_ids": [1, 3, 5],
    "district_ids": [11, 12, 31]
  }
}
```

## 3) Status kodlar

- `200` - Muvaffaqiyatli
- `400` - So'rov noto'g'ri (`region_ids`, `district_ids`, hierarchy)
- `401` - Token yo'q yoki yaroqsiz
- `404` - Region yoki district topilmadi
- `500` - Server xatoligi
