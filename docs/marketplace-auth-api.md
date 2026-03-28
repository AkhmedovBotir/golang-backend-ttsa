# Marketplace Auth API

Marketplace foydalanuvchilari uchun SMS-code asosidagi auth.

Base URL: `http://localhost:8081/api/v1`

## Biznes mantiq (strategik ketma-ketlik)

Asosiy oqim quyidagicha:

1. Foydalanuvchi telefon raqam kiritadi.
2. Tizim shu raqamga SMS kod yuboradi.
3. Foydalanuvchi SMS kodni kiritib tasdiqlaydi.
4. Tizim telefon raqamga qarab 2 xil yo'lga ajratadi:
   - **Raqam mavjud bo'lsa (profil bor):**
     - profil preview qaytadi;
     - foydalanuvchi o'ziniki ekanini tekshiradi;
     - to'g'ri bo'lsa darhol login bo'ladi (token qaytadi);
     - noto'g'ri bo'lsa boshqa raqam bilan qayta urinadi.
   - **Raqam mavjud bo'lmasa:**
     - registratsiya ma'lumotlari so'raladi;
     - registratsiya tugashi bilan avtomatik login bo'ladi (token qaytadi).

Muhim:
- Marketplace auth'da parol yo'q.
- Faqat SMS kod orqali tasdiqlash ishlatiladi.

## Postman uslubida endpointlar

### 1) Kirishni boshlash (telefon + SMS yuborish)

- **Method:** `POST`
- **URL:** `http://localhost:8081/api/v1/marketplace/auth/entry`
- **Request (JSON):**

```json
{
  "phone": "+998901234567"
}
```

- **Response (200) - agar profil mavjud bo'lsa:**

```json
{
  "message": "sms kod yuborildi",
  "data": {
    "flow": "login",
    "exists": true,
    "sent": true,
    "profile": {
      "id": 12,
      "phone": "+998901234567",
      "first_name": "Ali",
      "last_name": "Valiyev",
      "gender": "erkak",
      "region_id": 1,
      "district_id": 10,
      "mfy_id": 265,
      "birth_date": "1995-03-21",
      "status": "active",
      "created_at": "2026-03-20T08:00:00Z",
      "updated_at": "2026-03-20T08:00:00Z"
    }
  }
}
```

- **Response (200) - agar profil mavjud bo'lmasa:**

```json
{
  "message": "sms kod yuborildi",
  "data": {
    "flow": "register",
    "exists": false,
    "sent": true
  }
}
```

### 2) SMS kodni tasdiqlash va keyingi bosqich

- **Method:** `POST`
- **URL:** `http://localhost:8081/api/v1/marketplace/auth/entry/verify`
- **Request (JSON):**

```json
{
  "phone": "+998901234567",
  "code": "12345"
}
```

- **Response (200) - login holati (darhol token):**

```json
{
  "message": "muvaffaqiyatli login",
  "data": {
    "flow": "login",
    "verified": true,
    "needs_registration": false,
    "token": "jwt_token_here",
    "profile": {
      "id": 12,
      "phone": "+998901234567",
      "first_name": "Ali",
      "last_name": "Valiyev",
      "gender": "erkak",
      "region_id": 1,
      "district_id": 10,
      "mfy_id": 265,
      "birth_date": "1995-03-21",
      "status": "active",
      "created_at": "2026-03-20T08:00:00Z",
      "updated_at": "2026-03-20T08:00:00Z"
    }
  }
}
```

- **Response (200) - register holati (registratsiya kerak):**

```json
{
  "message": "kod tasdiqlandi, registratsiyani yakunlang",
  "data": {
    "flow": "register",
    "verified": true,
    "needs_registration": true
  }
}
```

### 3) Registratsiyani yakunlash (avtomatik login bilan)

- **Method:** `POST`
- **URL:** `http://localhost:8081/api/v1/marketplace/auth/register`
- **Request (JSON):**

```json
{
  "phone": "+998901234567",
  "first_name": "Ali",
  "last_name": "Valiyev",
  "gender": "erkak",
  "region_id": 1,
  "district_id": 10,
  "mfy_id": 265,
  "birth_date": "1995-03-21"
}
```

- **Response (201):**

```json
{
  "message": "ro'yxatdan o'tildi",
  "data": {
    "token": "jwt_token_here",
    "profile": {
      "id": 99,
      "phone": "+998901234567",
      "first_name": "Ali",
      "last_name": "Valiyev",
      "gender": "erkak",
      "region_id": 1,
      "district_id": 10,
      "mfy_id": 265,
      "birth_date": "1995-03-21",
      "status": "active",
      "created_at": "2026-03-26T10:00:00Z",
      "updated_at": "2026-03-26T10:00:00Z"
    }
  }
}
```

### 4) Login bo'lgandan keyin profil olish

- **Method:** `GET`
- **URL:** `http://localhost:8081/api/v1/marketplace/me/profile`
- **Header:** `Authorization: Bearer <token>`
- **Response (200):**

```json
{
  "message": "ok",
  "data": {
    "id": 99,
    "phone": "+998901234567",
    "first_name": "Ali",
    "last_name": "Valiyev",
    "gender": "erkak",
    "region_id": 1,
    "district_id": 10,
    "mfy_id": 265,
    "birth_date": "1995-03-21",
    "status": "active",
    "created_at": "2026-03-26T10:00:00Z",
    "updated_at": "2026-03-26T10:00:00Z"
  }
}
```

## Validatsiya qoidalari

- Telefon formati: `+998901234567`
- Gender: `erkak` yoki `ayol`
- `birth_date`: `YYYY-MM-DD`
- `region_id`, `district_id`, `mfy_id` o'zaro mos bo'lishi shart
- SMS kod amal qilish vaqti: 5 daqiqa

## Status kodlar

- `200` - muvaffaqiyatli
- `201` - muvaffaqiyatli yaratildi (register)
- `400` - so'rov xato / validation xato
- `403` - kod muddati o'tgan yoki oldin tasdiqlanmagan
- `404` - kerakli ma'lumot topilmadi
- `409` - bu telefon bo'yicha profil allaqachon mavjud
- `500` - server xatoligi
