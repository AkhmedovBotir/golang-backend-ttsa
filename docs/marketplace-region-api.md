# Marketplace Region API

Marketplace foydalanuvchilari uchun viloyat, tuman va MFY ma'lumotlarini olish API.

Base URL: `/api/v1`

Bu endpointlar uchun token talab qilinmaydi (public).

## 1) Viloyatlar

`GET /marketplace/regions`

Maqsad:
- tizimdagi barcha viloyatlarni olish.

## 2) Tumanlar

`GET /marketplace/districts`

Ixtiyoriy filter:
- `region_id` berilsa, faqat shu viloyatga tegishli tumanlar qaytadi.

Misol:
- `GET /marketplace/districts?region_id=1`

## 3) MFYlar

`GET /marketplace/mfys`

Ixtiyoriy filter:
- `district_id` berilsa, faqat shu tumanga tegishli MFYlar qaytadi.

Misol:
- `GET /marketplace/mfys?district_id=10`

## Foydalanish tartibi (UI uchun)

1. Avval viloyatlar ro'yxatini oling.
2. Foydalanuvchi viloyat tanlagach, shu viloyat bo'yicha tumanlarni chaqiring.
3. Foydalanuvchi tuman tanlagach, shu tuman bo'yicha MFYlarni chaqiring.

Shu ketma-ketlik registratsiya formasi va profil tahrirlash ekranlari uchun qulay.

## Status kodlar

- `200` - muvaffaqiyatli.
- `400` - `region_id` yoki `district_id` noto'g'ri formatda.
- `500` - server xatoligi.
