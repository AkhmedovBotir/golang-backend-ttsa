# Contragent Auth API

Base URL: `http://localhost:8081/api/v1`

Modul fayllari: `modules/contragents` (`domain`, `repository`, `service`, `handler`, `module.go`).

## Flow

1. `send-code` — telefon tekshiriladi, 5 xonali kod SMS orqali yuboriladi.
2. `verify-code` — kod to'g'ri/noto'g'ri yoki muddati o'tganligi tekshiriladi.
3. `set-password` — kod tasdiqlangandan keyin parol o'rnatiladi va token qaytadi (auto login).
4. Keyin `login` yoki token bilan `me/profile` / `me/change-password` / `me/logo` yangilash.

## Public endpointlar

### 1) SMS kod yuborish

**POST** `/contragents/auth/send-code`

```json
{
  "phone": "+998901112233"
}
```

Javob: `200` — `SMS kodi yuborildi`

### 2) SMS kodni tekshirish

**POST** `/contragents/auth/verify-code`

```json
{
  "phone": "+998901112233",
  "code": "12345"
}
```

Javob: `200` — `Kod tasdiqlandi`

### 3) Kodni qayta yuborish

**POST** `/contragents/auth/resend-code`

```json
{
  "phone": "+998901112233"
}
```

Javob: `200` — `SMS kodi qayta yuborildi`

### 4) Parol o'rnatish (kod tasdiqlangandan keyin)

**POST** `/contragents/auth/set-password`

```json
{
  "phone": "+998901112233",
  "password": "strongPass123"
}
```

Javobda token qaytadi:

```json
{
  "message": "Parol o'rnatildi",
  "data": {
    "token": "jwt-token",
    "contragent": {
      "id": 1,
      "name": "Test Contragent"
    }
  }
}
```

### 5) Login

**POST** `/contragents/auth/login`

```json
{
  "phone": "+998901112233",
  "password": "strongPass123"
}
```

## Protected endpointlar (Bearer token kerak)

Header:

`Authorization: Bearer <token>`

### 6) Profil

**GET** `/contragents/me/profile`

### 7) Parolni almashtirish

**POST** `/contragents/me/change-password`

```json
{
  "old_password": "strongPass123",
  "new_password": "newStrongPass456"
}
```

### 8) Logoni yangilash (base64)

**PATCH** `/contragents/me/logo`

```json
{
  "logo": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."
}
```

`logo` uchun oddiy base64 string ham qabul qilinadi:

```json
{
  "logo": "iVBORw0KGgoAAAANSUhEUgAA..."
}
```

Javobda `contragent` ichida:
- `logo` — yuborilgan base64
- `has_logo` — `true/false`

## Validatsiya va qoidalar

- `phone` formati: `+998901234567`
- SMS kod amal qilish muddati: **5 daqiqa**
- `set-password` faqat `verify-code` dan keyin ishlaydi
- `password` kamida 6 ta belgi
- `logo` base64 bo'lishi kerak (`data:image/...;base64,...` yoki plain base64)
- `send-code` / `resend-code` faqat:
  - kontragent mavjud bo'lsa
  - status `active` bo'lsa
  - parol o'rnatilmagan bo'lsa

## Status kodlar

- `200` — muvaffaqiyatli
- `400` — noto'g'ri so'rov yoki validatsiya
- `401` — login/token xato
- `403` — kontragent nofaol yoki parol o'rnatilgan
- `404` — kontragent yoki kod topilmadi
- `410` — SMS kod muddati o'tgan
- `500` — server xatosi
