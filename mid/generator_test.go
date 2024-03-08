package mid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {
	assert.IsType(t, "string", DefaultGenerator.Generate())
}
