package goavro

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

const (
	salsifyCompressedUUIDLength = 16
	salsifyUUIDLength           = 38
)

func salsifyCompressedUUIDTextualFromNative(buf []byte, d interface{}) ([]byte, error) {
	switch data := d.(type) {
	case []uint8:
		buf = append(buf, "\"s-"...)
		stringUUID := hex.EncodeToString(data)

		buf = append(buf, stringUUID[0:8]...)
		buf = append(buf, '-')

		buf = append(buf, stringUUID[8:12]...)
		buf = append(buf, '-')

		buf = append(buf, stringUUID[12:16]...)
		buf = append(buf, '-')

		buf = append(buf, stringUUID[16:20]...)
		buf = append(buf, '-')

		buf = append(buf, stringUUID[20:32]...)
		buf = append(buf, '"')
		return buf, nil
	default:
		return nil, fmt.Errorf("cannot encode textual salsify uuid, expected []uint8 received: %T", data)
	}
}

func salsifyCompressedUUIDNativeFromTextual(buf []byte) (interface{}, []byte, error) {
	var datum interface{}
	var err error
	datum, buf, err = bytesNativeFromTextual(buf)
	if err != nil {
		return nil, buf, err
	}
	datumBytes := datum.([]byte)
	if count := uint(len(datumBytes)); count != salsifyUUIDLength {
		return nil, nil, fmt.Errorf("cannot decode textual fixed salsify_uuid_binary: datum size ought to equal schema size: %d != %d", count, salsifyUUIDLength)
	}

	datumBytes, ok := bytes.CutPrefix(datumBytes, []byte("s-"))
	if !ok {
		return nil, nil, fmt.Errorf("cannot decode textual fixed salsify_uuid_binary: datum ought to start with a leading 's-'")
	}
	datumBytes = bytes.ReplaceAll(datumBytes, []byte("-"), nil)
	intermediateBuffer := make([]byte, hex.DecodedLen(len(datumBytes)))
	_, err = hex.Decode(intermediateBuffer, datumBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot decode textual fixed salsify_uuid_binary: could not decode hex representation to bytes")
	}

	return intermediateBuffer, buf, err
}

// Since there isn't any difference in how we'll handle "fixed" types went encoding/decoding
// from binary to native (or vice versa), these methods we're taken from the fixed type codec
// in fixed.go
func salsifyCompressedUUIDNativeFromBinary(buf []byte) (interface{}, []byte, error) {
	if buflen := uint(len(buf)); salsifyCompressedUUIDLength > buflen {
		return nil, nil, fmt.Errorf("cannot decode binary fixed salsify_uuid_binary: schema 16 exceeds remaining buffer size: %d > %d (short buffer)", buflen, salsifyCompressedUUIDLength)
	}
	return buf[:16], buf[16:], nil
}

func salsifyCompressedUUIDBinaryFromNative(buf []byte, datum interface{}) ([]byte, error) {
	var someBytes []byte
	switch d := datum.(type) {
	case []byte:
		someBytes = d
	case string:
		someBytes = []byte(d)
	default:
		return nil, fmt.Errorf("cannot encode binary fixed salsify_uuid_binary: expected []byte or string; received: %T", datum)
	}
	if count := uint(len(someBytes)); count != salsifyCompressedUUIDLength {
		return nil, fmt.Errorf("cannot encode binary fixed salsify_uuid_binary: datum size ought to equal schema size: %d != %d", count, salsifyCompressedUUIDLength)
	}
	return append(buf, someBytes...), nil

}
