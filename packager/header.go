package packager

import (
	"github.com/tomasdemarco/iso8583/header"
	"io"
)

type HeaderPackager interface {
	Pack(val header.Header) (valueRaw []byte, length int, err error)
	Unpack(r io.Reader) (val header.Header, length int, err error)
}
