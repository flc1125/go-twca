package mid

import (
	"fmt"
	"math/rand"
	"time"
)

type Generator interface {
	Generate() string
}

var DefaultGenerator Generator = &defaultGenerator{}

type defaultGenerator struct{}

func (g *defaultGenerator) Generate() string {
	return time.Now().Format("20060102150405") +
		fmt.Sprintf("%06d", rand.Intn(1000000)) //nolint:gomnd
}
