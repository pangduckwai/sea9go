# sea9go
My golang utilities

## Packages
### `http\client`
a convinent wrapper of `http.Client`, with support of specifying TLS server certificates and/or mTLS certificates when creating the clients.

### `http\server`
a wrapper of `http.Server` with start/stop/timeout handling

### `io`
implements reading from and writing to streams with the following features:
- buffered read from io.Reader
- read from stdin
- stackable encoding/decoding filters

### `logger`
implements reusable loggers with prefix and labels, also includes a fast utility for determining number of digits of integers.

### `logger/metric`
implements conversion of integer values to a form with metric suffix.

### `rand`
a wrapper of fast pseudo random values generation.

### `traverse`
implements traversal of `map[string]any` structures (e.g. from json/yaml).

### `traverse/ordered`
implements traversal of `[]yaml.MapItem` structures from yaml while keeping items order of the file.

---

## Changelog
### v0.1.0
- First pre-release version

### v0.2.1
- Add http client and server
