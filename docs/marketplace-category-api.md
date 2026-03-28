# Marketplace Category + Subcategory GET API

Marketplace uchun kategoriya va subkategoriyalarni olish API (faqat GET).

Base URL: `http://localhost:8081/api/v1`

`Authorization` talab qilinmaydi.

Muhim:
- faqat `active` statusdagi kategoriya/subkategoriyalar qaytadi.

## 1) Kategoriyalar ro'yxati

- **Method:** `GET`
- **URL:** `/marketplace/categories?page=1&limit=10`

## 2) Bitta kategoriya

- **Method:** `GET`
- **URL:** `/marketplace/categories/{id}`

## 3) Subkategoriyalar ro'yxati

- **Method:** `GET`
- **URL:** `/marketplace/subcategories?page=1&limit=10`

Ixtiyoriy filter:
- `parent_id` (ma'lum kategoriya ostidagi subkategoriyalarni olish uchun)

Misol:
- `/marketplace/subcategories?page=1&limit=20&parent_id=3`

## 4) Bitta subkategoriya

- **Method:** `GET`
- **URL:** `/marketplace/subcategories/{id}`

## Status kodlar

- `200` - muvaffaqiyatli
- `400` - ID yoki `parent_id` noto'g'ri
- `404` - kategoriya/subkategoriya topilmadi
- `500` - server xatoligi

