package goavro

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

const (
	salsifyCompressedUuidLength = 16
	salsifyUuidLength           = 38
)

func salsifyCompressedUuidTextualFromNative(buf []byte, d interface{}) ([]byte, error) {
	switch data := d.(type) {
	case []uint8:
		buf = append(buf, "s-"...)
		stringUuid := hex.EncodeToString(data)

		buf = append(buf, stringUuid[0:8]...)
		buf = append(buf, '-')

		buf = append(buf, stringUuid[8:12]...)
		buf = append(buf, '-')

		buf = append(buf, stringUuid[12:16]...)
		buf = append(buf, '-')

		buf = append(buf, stringUuid[16:20]...)
		buf = append(buf, '-')

		buf = append(buf, stringUuid[20:32]...)
		return buf, nil
	default:
		return nil, fmt.Errorf("cannot encode textual salsify uuid, expected []uint8 received: %T", data)
	}
}

func salsifyCompressedUuidNativeFromTextual(buf []byte) (interface{}, []byte, error) {
	var datum interface{}
	var err error
	datum, buf, err = bytesNativeFromTextual(buf)
	if err != nil {
		return nil, buf, err
	}
	datumBytes := datum.([]byte)
	if count := uint(len(datumBytes)); count != salsifyUuidLength {
		return nil, nil, fmt.Errorf("cannot decode textual fixed salsify_uuid_binary: datum size ought to equal schema size: %d != %d", count, salsifyUuidLength)
	}

	datumBytes, ok := bytes.CutPrefix(datumBytes, []byte("s-"))
	if !ok {
		return nil, nil, fmt.Errorf("cannot decode textual fixed salsify_uuid_binary: datum ought to start with a leading 's-'")
	}
	datumBytes = bytes.ReplaceAll(datumBytes, []byte("-"), nil)
	intermediateBuffer := make([]byte, hex.DecodedLen(len(datumBytes)))
	hex.Decode(intermediateBuffer, datumBytes)

	return intermediateBuffer, buf, err
}

// Since there isn't any difference in how we'll handle "fixed" types went encoding/decoding
// from binary to native (or vice versa), these methods we're taken from the fixed type codec
// in fixed.go
func salsifyCompressedUuidNativeFromBinary(buf []byte) (interface{}, []byte, error) {
	if buflen := uint(len(buf)); salsifyCompressedUuidLength > buflen {
		return nil, nil, fmt.Errorf("cannot decode binary fixed salsify_uuid_binary: schema 16 exceeds remaining buffer size: %d > %d (short buffer)", buflen, salsifyCompressedUuidLength)
	}
	return buf[:16], buf[16:], nil
}

func salsifyCompressedUuidBinaryFromNative(buf []byte, datum interface{}) ([]byte, error) {
	var someBytes []byte
	switch d := datum.(type) {
	case []byte:
		someBytes = d
	case string:
		someBytes = []byte(d)
	default:
		return nil, fmt.Errorf("cannot encode binary fixed salsify_uuid_binary: expected []byte or string; received: %T", datum)
	}
	if count := uint(len(someBytes)); count != salsifyCompressedUuidLength {
		return nil, fmt.Errorf("cannot encode binary fixed salsify_uuid_binary: datum size ought to equal schema size: %d != %d", count, salsifyCompressedUuidLength)
	}
	return append(buf, someBytes...), nil

}
