param(
    [string]$JsonPath = "scripts/ttsa.categories.json"
)

go run ./cmd/import-categories $JsonPath
