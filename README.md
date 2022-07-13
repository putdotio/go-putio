[![Golang CI Linter](https://github.com/putdotio/go-putio/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/putdotio/go-putio/actions/workflows/golangci-lint.yml)
[![Golang Tests](https://github.com/putdotio/go-putio/actions/workflows/go-test.yml/badge.svg)](https://github.com/putdotio/go-putio/actions/workflows/go-test.yml)


# putio

putio is a Go client library for accessing the [Put.io API v2](https://api.put.io/v2/docs).

## Documentation

Available on [GoDoc](http://godoc.org/github.com/putdotio/go-putio)

## Install

```sh
go get github.com/putdotio/go-putio@latest
go get golang.org/x/oauth2@latest
```

## Usage

```go
package main

import (
        "fmt"
        "log"
        "context"

        "golang.org/x/oauth2"
        "github.com/putdotio/go-putio"
)

func main() {
    oauthToken := "<YOUR-TOKEN-HERE>"
    tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: oauthToken})
    oauthClient := oauth2.NewClient(context.TODO(), tokenSource)

    client := putio.NewClient(oauthClient)

    const rootDir = 0
    root, err := client.Files.Get(context.TODO(), rootDir)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(root.Name)
}
```
