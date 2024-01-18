package tplparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suifengpiao14/tormgenerator/parser/tplparser"
)

func TestIsNameChar(t *testing.T) {
	t.Run("string number", func(t *testing.T) {
		md5 := []byte(`md5`)
		c := byte(md5[2])
		yes := tplparser.IsNameChar(c)
		assert.Equal(t, yes, true)
	})

	t.Run("number", func(t *testing.T) {
		yes := tplparser.IsNameChar(5)
		assert.Equal(t, yes, true)
	})
}
