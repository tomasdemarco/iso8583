package length

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/prefix"
	"io"
)

type PackFunc func(prefixValue prefix.Prefix, lenMessage int) ([]byte, error)
type UnpackFunc func(r io.Reader, prefixValue prefix.Prefix) (int, error)

func Pack(prefixValue prefix.Prefix, lenMessage int) ([]byte, error) {
	return prefix.Pack(prefixValue, lenMessage)
}

func Unpack(r io.Reader, prefixValue prefix.Prefix) (int, error) {

	prefixLength := prefix.GetPrefixLen(prefixValue.Type, prefixValue.Encoding)

	buf := make([]byte, prefixLength)
	_, err := r.Read(buf)
	if err != nil {
		if err != io.EOF {
			err = fmt.Errorf("reading length: %w", err)
		}

		return 0, err
	}

	result, _, err := prefix.Unpack(prefixValue, buf)
	if err != nil {
		return 0, err
	}

	return result, err
}
