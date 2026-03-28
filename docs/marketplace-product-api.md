# Marketplace Product GET API

Marketplace uchun kontragentlar mahsulotlarini ko'rish API.

Base URL: `/api/v1`

Muhim: `Authorization` token shart emas (public).

## 1) Mahsulotlar ro'yxati

`GET /marketplace/products`

Ixtiyoriy filterlar:
- `page` (default: `1`)
- `limit` (default: `10`, max: `100`)
- `category_id`
- `subcategory_id`
- `contragent_id`
- `status` (`active` yoki `inactive`)
- `q` (nom yoki description bo'yicha qidiruv)

Javobda `delivery_areas` maydoni bo'ladi:
- `delivery_areas.region_ids`
- `delivery_areas.district_ids`

Ya'ni har bir mahsulotda shu mahsulotni sotayotgan kontragentning yetkazib berish hududlari ko'rsatiladi.

## 2) Mahsulotni bitta ID bilan olish

`GET /marketplace/products/{id}`

Javobda ham `delivery_areas` maydoni qaytadi.

## Status kodlar

- `200` - muvaffaqiyatli
- `400` - noto'g'ri so'rov (parametrlar)
- `404` - mahsulot topilmadi
- `500` - server xatoligi

