# sea9go
My golang utilities

## Packages
### `http/client`
a convinent wrapper of `http.Client`, with support of specifying TLS server certificates and/or mTLS certificates when creating clients.

### `http/server`
a wrapper of `http.Server` with start/stop/timeout handling

### `io`
implements reading from and writing to streams with the following features:
- buffered read from io.Reader
- stackable encoding/decoding filters

### `io/line`
read lines from io.Reader

### `io/prompt`
read from stdin

### `logger`
- implements reusable loggers with prefix and labels,
- a fast utility for determining number of digits of integers (reference: github.com/doloopwhile/go-fastlog).

### `logger/metric`
implements conversion of integer values to decimal values with metric suffix (reference: github.com/bmkessler/fastdiv/).

### `rand`
a wrapper of fast pseudo random values generation from "github.com/bytedance/gopkg/lang/fastrand".

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

### v0.2.2
- Add read line

### v0.2.3
- Refactor package `io`