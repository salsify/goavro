package goavro

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextualFromnative(t *testing.T) {
	source := []byte("\xDB\x81P\x12X[K\xEC\x96\x1AR\xC4\x1E\x1D\xE0W")
	buf := make([]byte, 0)

	buf, err := salsifyTextualFromNative(buf, source)

	assert.NoError(t, err)
	assert.Equal(t, buf, []byte("s-db815012-585b-4bec-961a-52c41e1de057"))
}
