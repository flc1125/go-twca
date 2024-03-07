package mid

import (
	"fmt"
	"math/rand"
	"time"
)

type VerifyNoGenerator interface {
	Generate() string
}

var DefaultVerifyNoGenerator VerifyNoGenerator = &defaultVerifyNoGenerator{}

type defaultVerifyNoGenerator struct{}

func (g *defaultVerifyNoGenerator) Generate() string {
	return time.Now().Format("20060102150405") +
		fmt.Sprintf("%06d", rand.Intn(1000000))
}
