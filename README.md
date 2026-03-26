# sea9go
My golang utilities

## Packages
### `http\client`
- convinent wrapper of `http.Client`

### `http\server`
- wrapper of `http.Server` with start/stop/timeout handling

### `io`
- buffered read from io.Reader
- read from stdin
- stackable encoding/decoding filters

### `logger`
- build reusable loggers with prefix
- fast utility for determining number of digits of integers

### `logger/metric`
- convert values to metric suffix forms

### `rand`
- fast pseudo random values generation

### `traverse`
- traverse `map[string]any` structures (e.g. from json/yaml)

### `traverse/ordered`
- traverse `map[string]any` structures from yaml while keeping the order

---

## Changelog
### v0.1.0
- First pre-release version

### v0.2.0
- Add http client and server
