# Contragent Type API hujjati

Base URL: `http://localhost:8081/api/v1`

## Muhim
- Barcha endpointlar `Authorization: Bearer <token>` talab qiladi.
- Faqat `general` role ushbu endpointlardan foydalana oladi.
- Har bir javobda `message` bo'ladi.

## 1) Import (Mongo JSON -> PostgreSQL)

JSON fayl: `scripts/ttsa.contragenttypes.json`

PowerShell:
```powershell
powershell -ExecutionPolicy Bypass -File ./scripts/import-contragent-types.ps1
```

Yoki custom path bilan:
```powershell
powershell -ExecutionPolicy Bypass -File ./scripts/import-contragent-types.ps1 -JsonPath "C:\path\to\file.json"
```

To'g'ridan-to'g'ri:
```powershell
go run ./cmd/import-contragent-types scripts/ttsa.contragenttypes.json
```

## 2) Contragent Type CRUD

### Create
`POST /contragent_type`
```json
{
  "name": "Elektronika va maishiy texnika",
  "icon": "Devices",
  "status": "active"
}
```

### List
`GET /contragent_type`

Query params:
- `page` (ixtiyoriy, default: `1`)
- `limit` (ixtiyoriy, default: `10`, max: `100`)

Namuna javob:
```json
{
  "message": "Kontragent turlari ro'yxati olindi",
  "data": {
    "items": [
      {
        "id": 1,
        "external_id": "69569da144cf7f7cf949e426",
        "name": "Elektronika va maishiy texnika",
        "icon": "Devices",
        "status": "active"
      }
    ],
    "total": 15,
    "page": 1,
    "limit": 10,
    "total_pages": 2
  }
}
```

### Get by id
`GET /contragent_type/{id}`

### Update
`PUT /contragent_type/{id}`
```json
{
  "name": "Elektronika",
  "icon": "Devices",
  "status": "inactive"
}
```

### Status update
`PATCH /contragent_type/{id}/status`
```json
{
  "status": "active"
}
```

### Delete
`DELETE /contragent_type/{id}`

## 3) Status qiymatlari
- `active`
- `inactive`

## 4) Status kodlar
- `200` - Muvaffaqiyatli
- `201` - Yaratildi
- `400` - So'rov noto'g'ri
- `401` - Token yo'q yoki yaroqsiz
- `403` - Ruxsat yo'q (general emas)
- `404` - Topilmadi
- `500` - Server xatoligi
