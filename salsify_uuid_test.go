package goavro

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSalsifyCompressedUUIDTextualFromNative(t *testing.T) {
	source := []byte("\xDB\x81P\x12X[K\xEC\x96\x1AR\xC4\x1E\x1D\xE0W")
	buf := make([]byte, 0)

	buf, err := salsifyCompressedUUIDTextualFromNative(buf, source)

	assert.NoError(t, err)
	assert.Equal(t, buf, []byte("\"s-db815012-585b-4bec-961a-52c41e1de057\""))
}

func TestSalsifyCompressedUUIDNativeFromTextual(t *testing.T) {
	source := []byte("\"s-db815012-585b-4bec-961a-52c41e1de057\"")

	native, _, err := salsifyCompressedUUIDNativeFromTextual(source)

	assert.NoError(t, err)
	assert.Equal(t, native, []byte("\xDB\x81P\x12X[K\xEC\x96\x1AR\xC4\x1E\x1D\xE0W"))
}

func TestSalsifyCompressedEndToEnd(t *testing.T) {
	schema := `{"type":"fixed","name":"salsify_uuid_binary","namespace":"com.salsify","size":16}`
	rawBinary, err := os.ReadFile("fixtures/salsify_uuid_binary.avro")
	if err != nil {
		t.Fatal(err)
	}
	codec, err := NewCodec(schema)
	if err != nil {
		t.Fatal(err)
	}
	fromNative, _, err := codec.NativeFromBinary(rawBinary)
	if err != nil {
		t.Fatal(err)
	}
	fromTextual, err := codec.TextualFromNative(nil, fromNative)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fromTextual, []byte("\"s-e895a664-e7e5-43e6-9f5f-2320787ff569\""))

	toNative, _, err := codec.NativeFromTextual(fromTextual)
	if err != nil {
		t.Fatal(err)
	}
	toBinary, err := codec.BinaryFromNative(nil, toNative)
	if err != nil {
		t.Fatal(err)
	}

	// Go decoding adds a non-breaking space 0xA0 at the end.
	assert.Equal(t, toBinary[:16], rawBinary[:16])

	backUpNative, _, err := codec.NativeFromBinary(toBinary)
	if err != nil {
		t.Fatal(err)
	}
	backUpTextual, err := codec.textualFromNative(nil, backUpNative)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, backUpTextual, []byte("\"s-e895a664-e7e5-43e6-9f5f-2320787ff569\""))
}
