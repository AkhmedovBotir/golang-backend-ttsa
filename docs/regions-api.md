# Region / District / MFY API hujjati

Base URL: `http://localhost:8081/api/v1`

## Muhim
- Barcha endpointlar `Authorization: Bearer <token>` talab qiladi.
- Faqat `general` role ushbu CRUD endpointlardan foydalana oladi.
- Har bir javobda `message` bo'ladi.

## 1) Import (Mongo JSON -> PostgreSQL)

JSON fayl: `scripts/ttsa.regions.json`

PowerShell:
```powershell
powershell -ExecutionPolicy Bypass -File ./scripts/import-regions.ps1
```

Yoki custom path bilan:
```powershell
powershell -ExecutionPolicy Bypass -File ./scripts/import-regions.ps1 -JsonPath "C:\path\to\file.json"
```

To'g'ridan-to'g'ri:
```powershell
go run ./cmd/import-regions scripts/ttsa.regions.json
```

## 2) Viloyat (Region) CRUD

### Create region
`POST /regions`
```json
{
  "name": "Sirdaryo",
  "code": "Sirdaryo",
  "status": "active"
}
```

### List regions
`GET /regions`

### Get region by id
`GET /regions/{id}`

### Update region
`PUT /regions/{id}`
```json
{
  "name": "Sirdaryo viloyati",
  "code": "Sirdaryo",
  "status": "active"
}
```

### Delete region
`DELETE /regions/{id}`

### Region status toggle
`PATCH /regions/{id}/status`
```json
{
  "status": "inactive"
}
```

## 3) Tuman (District) CRUD

### Create district
`POST /districts`
```json
{
  "region_id": 1,
  "name": "Guliston tumani",
  "code": "Guliston",
  "status": "active"
}
```

### List districts
`GET /districts`

### Get district by id
`GET /districts/{id}`

### Update district
`PUT /districts/{id}`
```json
{
  "region_id": 1,
  "name": "Yangilangan tuman",
  "code": "Yangilangan",
  "status": "inactive"
}
```

### Delete district
`DELETE /districts/{id}`

### District status toggle
`PATCH /districts/{id}/status`
```json
{
  "status": "active"
}
```

## 4) MFY CRUD

### Create MFY
`POST /mfys`
```json
{
  "district_id": 1,
  "name": "Do'stlik MFY",
  "code": "Dostlik",
  "status": "active"
}
```

### List MFY
`GET /mfys`

### Get MFY by id
`GET /mfys/{id}`

### Update MFY
`PUT /mfys/{id}`
```json
{
  "district_id": 1,
  "name": "Do'stlik MFY 2",
  "code": "Dostlik-2",
  "status": "active"
}
```

### Delete MFY
`DELETE /mfys/{id}`

### MFY status toggle
`PATCH /mfys/{id}/status`
```json
{
  "status": "inactive"
}
```

## 5) Status kodlar
- `200` - Muvaffaqiyatli
- `201` - Yaratildi
- `400` - So'rov noto'g'ri
- `401` - Token yo'q yoki yaroqsiz
- `403` - Ruxsat yo'q (general emas)
- `404` - Topilmadi
- `500` - Server xatoligi
