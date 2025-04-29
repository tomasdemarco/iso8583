package field

import (
	"errors"
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/padding"
	"github.com/tomasdemarco/iso8583/prefix"
	"regexp"
)

type Field struct {
	Description string           `json:"description"`
	Type        Type             `json:"type"`
	Length      int              `json:"length"`
	Pattern     string           `json:"pattern"`
	Encoding    encoding.Encoder `json:"encoding"`
	Prefix      prefix.Prefixer  `json:"prefix"`
	Padding     padding.Padder   `json:"padding"`
	//SubfieldsData subfield.SubfieldsData `json:"subFieldsData"`
}

func (f Field) Unpack(messageRaw []byte, position int) (string, int, error) {

	var length int
	var err error

	if f.Prefix != nil {
		length, err = f.Prefix.DecodeLength(messageRaw, position)
		if err != nil {
			return "", 0, err
		}
	}

	if length == 0 {
		length = f.Length
	} else {
		position += f.Prefix.GetPackedLength()

		if _, ok := f.Encoding.(*encoding.BCD); ok {
			length = length / 2
		}
	}

	f.Padding.SetEncoder(f.Encoding)

	paddingLeft, paddingRight := f.Padding.DecodePad(length)

	length += paddingLeft + paddingRight

	f.Encoding.SetLength(length)

	if len(messageRaw) < position+length {
		return "", 0, errors.New("index out of range while trying to unpack")
	}

	value, err := f.Encoding.Decode(messageRaw[position:])
	if err != nil {
		return "", 0, err
	}

	value = value[paddingLeft : len(value)-paddingRight]
	match, _ := regexp.MatchString(f.Pattern, value)
	if !match {
		err = errors.New("invalid format")
		return "", 0, err
	}

	//if fieldPackager.SubFields != nil { //TODO ver como resolver subfields
	//	m.UnpackSubfields(field, value)
	//}

	return value, length + f.Prefix.GetPackedLength(), nil
}

func (f Field) Pack(value string) ([]byte, error) {

	f.Padding.SetEncoder(f.Encoding)
	padLeft, padRight := f.Padding.EncodePad(f.Length, len(value))
	fieldEncode, err := f.Encoding.Encode(padLeft + value + padRight)
	if err != nil {
		return nil, err
	}

	length := len(fieldEncode)

	if _, ok := f.Encoding.(*encoding.BCD); ok {
		length = length * 2
	}

	fieldPrefix, err := f.Prefix.EncodeLength(length)
	if err != nil {
		return nil, err
	}

	return append(fieldPrefix, fieldEncode...), nil
}
