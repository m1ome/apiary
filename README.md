# Apiary API - Golang library

[![Build Status](https://travis-ci.org/m1ome/apiary.svg?branch=master)](https://travis-ci.org/m1ome/apiary)
[![GoDoc](https://godoc.org/github.com/m1ome/apiary?status.svg)](https://godoc.org/github.com/m1ome/apiary)

# Description
This is a small golang library that will provide support for [Apiary](apiary.io) API.

# Installation
```
go get github.com/m1ome/apiary
```

# Usage
```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/m1ome/apiary"
)

func main() {
    token := os.Getenv("APIARY_TOKEN")

    api := NewApiary(ApiaryOptions{
        Token: Token,
    })

    response, err := api.Me()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("ID: %d\n", response.ID)
    fmt.Printf("Name: %s\n", response.Name)
    fmt.Printf("URL: %s\n", response.URL)
}

```