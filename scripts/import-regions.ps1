param(
    [string]$JsonPath = "scripts/ttsa.regions.json"
)

go run ./cmd/import-regions $JsonPath
