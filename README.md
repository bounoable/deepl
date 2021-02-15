# DeepL Pro API client

[![PkgGoDev](https://pkg.go.dev/badge/github.com/bounoable/deepl)](https://pkg.go.dev/github.com/bounoable/deepl)

Client library for the [**DeepL Pro API**](https://deepl.com).

## Installation

```sh
go get github.com/bounoable/deepl
```

## Usage

See the [examples](./example_test.go).

```go
import (
  "github.com/bounoable/deepl"
)

client := deepl.New("your-auth-key")

translated, sourceLang, err := client.Translate(
  context.TODO(),
  "Hello, world",
  deepl.Chinese,
)
if err != nil {
  log.Fatal(err)
}

log.Println(fmt.Sprintf("source language: %s", sourceLang))
log.Println(translated)
```

## Testing

You can test the library against the real DeepL API by running the following command.

**CAUTION: Runnning these tests will add to your usage and therefore will be billed!**

```sh
make e2e-test authKey=YOUR_AUTH_KEY
```

## License

[MIT](./LICENSE)
