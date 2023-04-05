package sessions

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

func init() {
	gob.Register(map[string]interface{}{})
}

func Encode(value StoreValue) ([]byte, error) {
	buf := bytes.Buffer{}
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	src := buf.Bytes()
	dst := make([]byte, base64.RawURLEncoding.EncodedLen(len(src)))
	base64.RawURLEncoding.Encode(dst, src)
	return dst, nil
}

func EncodeToString(value StoreValue) (string, error) {
	buf := bytes.Buffer{}
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf.Bytes()), nil
}

func Decode(src []byte, value *StoreValue) error {
	if len(src) == 0 {
		return nil
	}
	buf := make([]byte, base64.RawURLEncoding.DecodedLen(len(src)))
	if n, err := base64.RawURLEncoding.Decode(buf, src); err != nil {
		return err
	} else {
		buf = buf[:n]
	}
	reader := bytes.NewReader(buf)
	decoder := gob.NewDecoder(reader)
	if err := decoder.Decode(value); err != nil {
		return err
	}
	return nil
}

func DecodeString(src string, value *StoreValue) error {
	if src == "" {
		return nil
	}
	buf, err := base64.RawURLEncoding.DecodeString(src)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(buf)
	decoder := gob.NewDecoder(reader)
	if err = decoder.Decode(value); err != nil {
		return err
	}
	return nil
}
