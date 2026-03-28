# Marketplace qidiruv va kontragentlar API

Marketplace uchun umumiy qidiruv hamda kontragentlar ro'yxati (joylashuv, yetkazib berish hududlari, ixtiyoriy mahsulotlar yoki kategoriya daraxti).

Base URL: `http://localhost:8081/api/v1`

`Authorization` talab qilinmaydi.

---

## 1) Umumiy qidiruv

Bitta so'rov bo'yicha **mahsulotlar**, **asosiy kategoriyalar**, **subkategoriyalar** va **kontragentlar** dan mos keluvchi yozuvlar qaytadi (har bir tur bo'yicha alohida ro'yxat).

- **Method:** `GET`
- **URL:** `/marketplace/search`

### Query parametrlar

| Parametr | Tavsif |
|----------|--------|
| `q` | Qidiruv qatori (bo'sh bo'lsa, barcha bo'limlar bo'sh qaytadi). |
| `limit_per_type` | Har bir tur uchun maksimal yozuv soni. Default: `10`, min: `1`, max: `50`. |
| `types` | Qaysi turlar qidirilsin — vergul bilan. Bo'sh yoki noto'g'ri qiymatda **hammasi** qidiriladi. Qiymatlar: `products`, `categories`, `subcategories`, `contragents`. |

Misol: faqat mahsulot va kontragent:  
`/marketplace/search?q=sut&types=products,contragents&limit_per_type=15`

### Qidiruv maydonlari

- **Mahsulotlar:** `name`, `description` (faqat `moderation_status = approved`, `status = active`).
- **Kategoriyalar:** `name`, `slug` (faqat asosiy, `active`).
- **Subkategoriyalar:** `name`, `slug` (`active`).
- **Kontragentlar:** `name`, `inn`, `phone` (faqat `status = active`).

### Javob tuzilishi (`data`)

`UnifiedSearchResponse`:

- `query` — qayta ishlangan qidiruv qatori
- `limit_per_type`
- `products` — `ProductOutput[]` (mahsulotlar ro'yxati API bilan bir xil maydonlar, shu jumladan `delivery_areas`)
- `categories` — `{ id, name, slug, image }`
- `subcategories` — `{ id, name, slug, image, parent_id }`
- `contragents` — `{ id, name, inn, phone, logo, region_id, district_id, mfy_id, delivery_areas }`  
  (`delivery_areas`: `{ region_ids, district_ids }`)

---

## 2) Kontragentlar ro'yxati (kengaytirilgan)

Faol kontragentlar sahifalangan ro'yxati. Har bir elementda **o'z joylashuvi** (`region_id`, `district_id`, `mfy_id`) va **yetkazib berish hududlari** (`delivery_areas`) bor. Ixtiyoriy ravishda ichida **mahsulotlar** yoki **kategoriya → subkategoriyalar** daraxti qo'shiladi.

- **Method:** `GET`
- **URL:** `/marketplace/contragents`

### Query parametrlar

| Parametr | Tavsif |
|----------|--------|
| `page` | Sahifa. Default: `1`. |
| `limit` | Sahifadagi yozuvlar. Default: `10`, max: `100`. |
| `q` | Ixtiyoriy: nom, STIR yoki telefon bo'yicha `ILIKE` filtr. |
| `include` | Vergul bilan: `products` va/yoki `categories` (yoki `subcategories` — `categories` bilan bir xil ishlaydi). Bo'sh bo'lsa, faqat asosiy kontragent maydonlari + joylashuv + `delivery_areas`. |
| `nested_limit` | `include=products` bo'lganda har bir kontragent uchun qaytariladigan mahsulotlar soni. Default: `30`, max: `100`. |

Misollar:

- Faqat ro'yxat: `/marketplace/contragents?page=1&limit=20`
- Qidiruv + mahsulotlar: `/marketplace/contragents?q=choy&include=products&nested_limit=10`
- Kategoriya tuzilmasi: `/marketplace/contragents?include=categories`  
  (faqat tasdiqlangan va faol mahsulotlari bo'lgan kategoriya/subkategoriya juftliklari)

### Javob (`data`)

```json
{
  "items": [ /* ContragentBrowseItem */ ],
  "total": 0,
  "page": 1,
  "limit": 10,
  "total_pages": 1
}
```

Har bir `items[]` elementi:

- Asosiy: `id`, `name`, `inn`, `phone`, `logo`, `region_id`, `district_id`, `mfy_id`, `activity_type_id`, `status`, `delivery_areas`
- Agar `include=products`: `products` — `ProductOutput[]` (kamida `nested_limit` tagacha)
- Agar `include=categories`: `category_branches` — `[{ "category": { ... }, "subcategories": [ ... ] }]`  
  har bir filialda bitta asosiy kategoriya va unga tegishli (mahsulotda uchragan) subkategoriyalar.

---

## Status kodlar

- `200` — muvaffaqiyatli
- `500` — server xatoligi
