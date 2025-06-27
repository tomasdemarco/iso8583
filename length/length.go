package length

import (
	"bufio"
	"fmt"
	"github.com/tomasdemarco/iso8583/prefix"
	"io"
)

type PackFunc func(prefixer prefix.Prefixer, lenMessage int) ([]byte, error)
type UnpackFunc func(r *bufio.Reader, prefixer prefix.Prefixer) (int, error)

func Pack(prefixer prefix.Prefixer, lenMessage int) ([]byte, error) {
	return prefixer.EncodeLength(lenMessage)
}

func Unpack(r *bufio.Reader, prefixer prefix.Prefixer) (int, error) {

	buf := make([]byte, prefixer.GetPackedLength())
	_, err := io.ReadFull(r, buf)
	if err != nil {
		if err != io.EOF {
			err = fmt.Errorf("reading length: %w", err)
		}

		return 0, err
	}

	result, err := prefixer.DecodeLength(buf, 0)
	if err != nil {
		return 0, err
	}

	return result, err
}
