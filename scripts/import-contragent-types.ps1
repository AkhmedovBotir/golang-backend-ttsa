param(
    [string]$JsonPath = "scripts/ttsa.contragenttypes.json"
)

go run ./cmd/import-contragent-types $JsonPath
