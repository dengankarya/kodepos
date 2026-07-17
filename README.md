# kodepos

Indonesian postal code search API — rewritten in Go.

> Go rewrite of [sooluh/kodepos](https://github.com/sooluh/kodepos) by Abu Masyail.
> Same API contract, same data, zero external dependencies.

## Why Go?

| | Node.js (original) | Go (this) |
|---|---|---|
| Image size | ~100MB (node:alpine) | ~10MB (`scratch`) |
| Dependencies | 6 npm packages | stdlib only |
| Cold start | ~200ms | ~5ms |
| Memory | ~80MB | ~30MB |

## Quick Start

```bash
git clone https://github.com/dengankarya/kodepos.git && cd kodepos
API_KEY=secret go run .
# → http://localhost:3000
```

## Docker

```bash
docker build -t kodepos .
docker run -p 3000:3000 -e API_KEY=your_secret kodepos
```

## Environment Variables

| Variable  | Default | Description                                     |
| --------- | ------- | ----------------------------------------------- |
| `PORT`    | `3000`  | Server listen port                              |
| `API_KEY` | (empty) | API key for `X-API-KEY` header. Empty = no auth |

## Authentication

All requests require the `X-API-KEY` header:

```
X-API-KEY: your_secret
```

If `API_KEY` is unset, authentication is disabled.

## API

### `GET /search`

Fuzzy search by place name.

| Param | Type   | Required | Description     |
| ----- | ------ | :------: | --------------- |
| `q`   | string |    ✅    | Search keywords |

```bash
curl -H "X-API-KEY: secret" "http://localhost:3000/search?q=danasari"
```

<details>
<summary>Response</summary>

```json
{
  "statusCode": 200,
  "code": "OK",
  "data": [
    {
      "code": 46386,
      "village": "Danasari",
      "district": "Cisaga",
      "regency": "Ciamis",
      "province": "Jawa Barat",
      "latitude": -7.3271342,
      "longitude": 108.4577572,
      "elevation": 110,
      "timezone": "WIB"
    }
  ]
}
```

</details>

---

### `GET /detect`

Find nearest postal code by GPS coordinates.

| Param       | Type   | Required | Description |
| ----------- | ------ | :------: | ----------- |
| `latitude`  | number |    ✅    | Latitude    |
| `longitude` | number |    ✅    | Longitude   |

```bash
curl -H "X-API-KEY: secret" "http://localhost:3000/detect?latitude=-6.547&longitude=107.398"
```

<details>
<summary>Response</summary>

```json
{
  "statusCode": 200,
  "code": "OK",
  "data": {
    "code": 41152,
    "village": "Kembangkuning",
    "district": "Jatiluhur",
    "regency": "Purwakarta",
    "province": "Jawa Barat",
    "latitude": -6.5495591,
    "longitude": 107.4121855,
    "elevation": 112,
    "timezone": "WIB",
    "distance": 1.589
  }
}
```

</details>

---

### `GET /`

- With `?q=` → redirects to `/search?q=...`
- Without → redirects to [github.com/dengankarya/kodepos](https://github.com/dengankarya/kodepos)

## Error Responses

```json
{ "statusCode": 400, "code": "BAD_REQUEST",    "message": "The 'q' parameter is required." }
{ "statusCode": 401, "code": "UNAUTHORIZED",    "message": "Invalid or missing API key." }
{ "statusCode": 404, "code": "NOT_FOUND",       "message": "This endpoint cannot be found." }
{ "statusCode": 500, "code": "INTERNAL_SERVER_ERROR", "message": "Please contact the developer." }
```

## Testing

```bash
go test ./... -v
go vet ./...
```

## Project Structure

```
├── main.go          # Server, routing, embedded data
├── model.go         # PostalCode, APIResponse types
├── search.go        # Inverted index + token search
├── detect.go        # Haversine distance + nearest detection
├── handler.go       # HTTP handlers (/, /search, /detect)
├── middleware.go     # API key auth + gzip compression
├── data/
│   └── kodepos.json # 83k+ postal codes (embedded at compile time)
└── Dockerfile       # Multi-stage build → scratch
```

## Credits

Original project by [sooluh/kodepos](https://github.com/sooluh/kodepos) (Abu Masyail).

## License

[Apache 2.0](https://github.com/dengankarya/kodepos/blob/main/LICENSE)
