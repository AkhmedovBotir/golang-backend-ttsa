# Backend (Go + PostgreSQL)

## Ishga tushirish

1. PostgreSQL bazani yarating/yoki mavjud DB nomini `.env` ga yozing
2. `.env` fayl qiymatlarini moslang
3. Dependency:
   - `go mod tidy`
4. General admin yaratish:
   - `powershell -ExecutionPolicy Bypass -File ./scripts/seed-general-admin.ps1`
5. API server:
   - `go run ./cmd/api`
6. Region/District/MFY import:
   - `powershell -ExecutionPolicy Bypass -File ./scripts/import-regions.ps1`
7. Contragent Type import:
   - `powershell -ExecutionPolicy Bypass -File ./scripts/import-contragent-types.ps1`
8. Category/Subcategory import:
   - `powershell -ExecutionPolicy Bypass -File ./scripts/import-categories.ps1`

## Modular monolit struktura

- `cmd/api` - server entrypoint
- `cmd/seed-general-admin` - general admin seed script
- `cmd/import-regions` - mongo formatdagi region JSON import script
- `cmd/import-contragent-types` - mongo formatdagi contragent type JSON import script
- `cmd/import-categories` - mongo formatdagi category/subcategory JSON import script
- `internal/config` - environment config
- `internal/platform/database` - postgres ulanishi
- `internal/platform/httpserver` - gin server setup
- `internal/pkg/response` - yagona response formati
- `internal/pkg/security` - parol hash va JWT
- `modules/admin` - admin + region/district/mfy + contragent type + kontragent + neighborhood shop + agent + punkt + manager + category/subcategory + product CRUD/moderation
- `modules/contragents` - contragent auth (sms code, login, profile, change-password)
- `docs` - API hujjatlari (`admin-crud-api.md`, `admin-product-api.md`, `regions-api.md`, `contragent-type-api.md`, `contragent-api.md`, `contragent-auth-api.md`, `contragent-category-api.md`, `contragent-region-delivery-api.md`, `contragent-product-api.md`, `neighborhood-shop-api.md`, `agent-api.md`, `punkt-api.md`, `manager-api.md`, `category-subcategory-api.md`)
"# golang-backend-ttsa" 
