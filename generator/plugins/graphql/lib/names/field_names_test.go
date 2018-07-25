package names

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterNotSupportedFieldNameCharacters(t *testing.T) {
	var tests = map[string]string{
		"-hello":  "hello",
		"0hello":  "hello",
		"h0el_lo": "h0el_lo",
	}
	for in, out := range tests {
		assert.Equal(t, out, FilterNotSupportedFieldNameCharacters(in))
	}
}
