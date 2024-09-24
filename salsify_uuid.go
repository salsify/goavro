package goavro

import (
	"encoding/hex"
	"fmt"
)

func salsifyTextualFromNative(buf []byte, d interface{}) ([]byte, error) {
	switch data := d.(type) {
	case []uint8:
		buf = append(buf, "s-"...)
		string_uuid := hex.EncodeToString(data)

		buf = append(buf, string_uuid[0:8]...)
		buf = append(buf, "-"...)

		buf = append(buf, string_uuid[8:12]...)
		buf = append(buf, "-"...)

		buf = append(buf, string_uuid[12:16]...)
		buf = append(buf, "-"...)

		buf = append(buf, string_uuid[16:20]...)
		buf = append(buf, "-"...)

		buf = append(buf, string_uuid[20:32]...)
		return buf, nil
	default:
		return nil, fmt.Errorf("cannot encode textual salsify uuid, expected []uint8 received: %T", data)
	}
}

func salsifyNativeFromBinary(buf []byte) (interface{}, []byte, error) {
	if buflen := uint(len(buf)); 16 > buflen {
		return nil, nil, fmt.Errorf("cannot decode binary fixed %q: schema 16 exceeds remaining buffer 16: %d > %d (short buffer)", "com.salsify.salsify_uuid_binary", 16, buflen)
	}
	return buf[:16], buf[16:], nil
}
