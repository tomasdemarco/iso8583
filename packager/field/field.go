package field

import (
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/padding"
	"github.com/tomasdemarco/iso8583/prefix"
	"github.com/tomasdemarco/iso8583/utils"
	"regexp"
)

type Packager interface {
	Pack(value string) ([]byte, string, error)
	Unpack(messageRaw []byte, position int) (string, int, error)
	Length() int
	Pattern() *regexp.Regexp
	Encoder() encoding.Encoder
	Prefixer() prefix.Prefixer
	Padder() padding.Padder
	Bitmap() *utils.BitSet
	SetBitmap(bmap *utils.BitSet)
}

// Field represents an ISO 8583 field's definition, including its data type,
// length, validation pattern, encoding, prefix, and padding rules.
type Field struct {
	Description string
	Type        Type
	length      int
	pattern     *regexp.Regexp
	encoder     encoding.Encoder
	prefixer    prefix.Prefixer
	padder      padding.Padder
}

func NewField(
	description string,
	fieldType Type,
	length int,
	pattern *regexp.Regexp,
	encoding encoding.Encoder,
	prefix prefix.Prefixer,
	padding padding.Padder,
) Packager {
	return &Field{
		Description: description,
		Type:        fieldType,
		length:      length,
		pattern:     pattern,
		encoder:     encoding,
		prefixer:    prefix,
		padder:      padding,
	}
}

// Unpack unpacks a field's value from a raw message byte slice.
// It handles length prefixes, padding, and decoding based on the field's configuration.
// It returns the unpacked field value as a string, the total length consumed in bytes,
// and an error if unpacking fails.
func (f Field) Unpack(messageRaw []byte, position int) (string, int, error) {

	var length int
	var err error

	if f.Prefixer() != nil {
		length, err = f.Prefixer().DecodeLength(messageRaw, position)
		if err != nil {
			return "", 0, err
		}
	}

	if length == 0 {
		length = f.Length()
	} else {
		position += f.Prefixer().GetPackedLength()

		if _, ok := f.Encoder().(*encoding.BCD); ok {
			length = length / 2
		}
	}

	paddingLeft, paddingRight := f.Padder().DecodePad(length)

	length += paddingLeft + paddingRight

	f.Encoder().SetLength(length)

	if len(messageRaw) < position+length {
		return "", 0, ErrUnpackIndexOutOfRange
	}

	value, err := f.Encoder().Decode(messageRaw[position:])
	if err != nil {
		return "", 0, err
	}

	value = value[paddingLeft : len(value)-paddingRight]

	if !f.Pattern().MatchString(value) {
		return "", 0, ErrInvalidFieldFormat
	}

	return value, length + f.Prefixer().GetPackedLength(), nil
}

// Pack packs a field's string value into a byte slice according to its configuration.
// It applies padding, encodes the value, and prepends the length prefix if defined.
// It returns the packed field as a byte slice, the plain (padded) field string,
// and an error if packing fails.
func (f Field) Pack(value string) ([]byte, string, error) {
	padLeft, padRight, err := f.Padder().EncodePad(f.Length(), len(value), f.Encoder())
	if err != nil {
		return nil, "", err
	}

	paddedField := padLeft + value + padRight

	fieldEncode, err := f.Encoder().Encode(paddedField)
	if err != nil {
		return nil, "", err
	}
	length := len(fieldEncode)

	if _, ok := f.Encoder().(*encoding.BCD); ok {
		length = length * 2

		if f.Padder().Type() == padding.Parity && len(padLeft+padRight) > 0 {
			length -= 1
		}
	}

	fieldPrefix, err := f.Prefixer().EncodeLength(length)
	if err != nil {
		return nil, "", err
	}

	return append(fieldPrefix, fieldEncode...), paddedField, nil
}

func (f Field) Length() int {
	return f.length
}

func (f Field) Pattern() *regexp.Regexp {
	return f.pattern
}

func (f Field) Encoder() encoding.Encoder {
	return f.encoder
}

func (f Field) Padder() padding.Padder {
	return f.padder
}

func (f Field) Prefixer() prefix.Prefixer {
	return f.prefixer
}

func (f Field) Bitmap() *utils.BitSet {
	return nil
}

func (f Field) SetBitmap(*utils.BitSet) {
}
