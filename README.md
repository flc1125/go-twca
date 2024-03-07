# Go TWCA

[![Go Version](https://badgen.net/github/release/flc1125/go-twca/stable)](https://github.com/flc1125/go-twca/releases)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-kratos-ecosystem)](https://pkg.go.dev/github.com/flc1125/go-twca)
[![codecov](https://codecov.io/gh/flc1125/go-twca/graph/badge.svg?token=QPTHZ5L9GT)](https://codecov.io/gh/flc1125/go-twca)
[![Go Report Card](https://goreportcard.com/badge/github.com/flc1125/go-twca)](https://goreportcard.com/report/github.com/flc1125/go-twca)
[![lint](https://github.com/flc1125/go-twca/actions/workflows/lint.yml/badge.svg)](https://github.com/flc1125/go-twca/actions/workflows/lint.yml)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

TWCA is a Go library for Taiwan Certificate Authority ([TWCA](https://www.twca.com.tw/)) services.

## MID

### Installation

```bash
go get github.com/flc1125/twca/mid
```

### Usage

```go
package main

import (
	"context"
	"time"

	"github.com/flc1125/go-twca/mid"
)

var ctx = context.Background()

func init() {
	time.Local = time.UTC
}

func main() {
	client := mid.New(&mid.Config{
		Addr:       "https://midonlinetest.twca.com.tw",
		BusinessNo: ".....",
		HashKeyNo:  ".....",
		HashKey:    ".....",
	})

	// MIDClause
	resp, err := client.MIDClause(ctx)
	if err != nil {
		panic(err)
	}
	_ = resp

	// ServerSideTransaction
	resp2, err := client.ServerSideTransaction(ctx, mid.ServerSideTransactionRequest{
		VerifyNo: mid.DefaultGenerator.Generate(),
		MemberNo: ".....",
		Action:   mid.ValidateMSISDNAdvanceAction,
		MIDInputParams: &mid.MIDInputParams{
			Msisdn:     ".....",
			Birthday:   nil,
			ClauseVer:  ".....",
			ClauseTime: ".....",
		},
	})
	if err != nil {
		panic(err)
	}
	_ = resp2

	// ServerSideVerifyResult
	resp3, err := client.ServerSideVerifyResult(ctx, mid.ServerSideVerifyResultRequest{
		VerifyNo: resp2.VerifyNo,
		MemberNo: resp2.OutputParams.MemberNo,
		Token:    resp2.OutputParams.Token,
	})
	if err != nil {
		panic(err)
	}
	_ = resp3

}
```

## License

[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
