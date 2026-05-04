package packager

import (
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/field"
	"github.com/tomasdemarco/iso8583/padding"
	"github.com/tomasdemarco/iso8583/prefix"
	"github.com/tomasdemarco/iso8583/utils"
	"regexp"
)

type FieldPackager interface {
	Pack(value string) ([]byte, string, error)
	Unpack(messageRaw []byte, position int) (string, int, error)
	NewDataField() field.Field
	Length() int
	Pattern() *regexp.Regexp
	Encoder() encoding.Encoder
	Prefixer() prefix.Prefixer
	Padder() padding.Padder
	Bitmap() *utils.BitSet
	SetBitmap(bmap *utils.BitSet)
	GetType() FieldType
}

// Field represents an ISO 8583 field's definition, including its data type,
// length, validation pattern, encoding, prefix, and padding rules.
type Field struct {
	Description      string
	Type             FieldType
	length           int
	encoder          encoding.Encoder
	prefixer         prefix.Prefixer
	padder           padding.Padder
	pattern          *regexp.Regexp
	dataFieldFactory func() field.Field
}

func NewField(
	description string,
	length int,
	enc encoding.Encoder,
	opts ...Option,
) *Field {
	fld := Field{
		Description: description,
		Type:        getFieldType(enc), // Default type, can be overridden by WithFieldType
		length:      length,
		encoder:     enc,
	}

	for _, opt := range opts {
		opt(&fld)
	}

	return &fld
}

func getFieldType(encoder encoding.Encoder) FieldType {
	switch encoder.(type) {
	case *encoding.BcdEncoder:
		return Numeric
	case *encoding.AsciiEncoder, *encoding.EbcdicEncoder:
		return String
	case *encoding.BinaryEncoder:
		return Binary
	}
	return String
}

type Option func(*Field)

func WithPrefix(prefixer prefix.Prefixer, opts ...prefix.Option) Option {
	return func(s *Field) {
		s.prefixer = prefixer

		for _, opt := range opts {
			opt(s.prefixer)
		}
	}
}

func WithPadding(padder padding.Padder, opts ...padding.Option) Option {
	return func(s *Field) {
		s.padder = padder

		if s.Type == Numeric {
			s.padder.SetChar("0")
		} else if s.Type == String {
			s.padder.SetChar(" ")
		}

		for _, opt := range opts {
			opt(s.padder)
		}
	}
}

func WithPattern(pattern *regexp.Regexp) Option {
	return func(s *Field) {
		s.pattern = pattern
	}
}

func WithFieldType(fieldType FieldType) Option {
	return func(s *Field) {
		s.Type = fieldType
	}
}

func WithCustomType(factory func() field.Field) Option {
	return func(f *Field) {
		f.dataFieldFactory = factory
	}
}

// NewDataField creates a new data field instance based on the packager's configuration.
// It uses the custom factory if provided, otherwise creates a default field type based on Field.Type.
func (f *Field) NewDataField() field.Field {
	if f.dataFieldFactory != nil {
		return f.dataFieldFactory()
	}

	switch f.Type {
	case Numeric, String:
		return &field.String{}
	case Binary, Bitmap:
		return &field.Bytes{}
	default:
		return &field.String{}
	}
}

// Unpack unpacks a field's value from a raw message byte slice.
func (f *Field) Unpack(messageRaw []byte, position int) (string, int, error) {
	var length, paddingLeft, paddingRight int
	var err error

	if f.Prefixer() != nil {
		length, err = f.Prefixer().DecodeLength(messageRaw, position)
		if err != nil {
			return "", 0, err
		}
		position += f.Prefixer().GetPackedLength()
	} else {
		length = f.Length()
	}

	if f.Padder() != nil {
		paddingLeft, paddingRight = f.Padder().Unpack(length)
	}

	length += paddingLeft + paddingRight
	f.Encoder().SetLength(length)

	if len(messageRaw) < position+length {
		return "", 0, ErrUnpackIndexOutOfRange
	}

	value, err := f.Encoder().Decode(messageRaw[position:])
	if err != nil {
		return "", 0, err
	}

	length = f.Encoder().GetLength()
	value = value[paddingLeft : len(value)-paddingRight]

	if f.Pattern() != nil && !f.Pattern().MatchString(value) {
		return "", 0, ErrInvalidFieldFormat
	}

	if f.Prefixer() != nil {
		length += f.Prefixer().GetPackedLength()
	}

	return value, length, nil
}

// Pack packs a field's string value into a byte slice according to its configuration.
func (f *Field) Pack(value string) ([]byte, string, error) {
	if f.Padder() != nil {
		padLeft, padRight, err := f.Padder().Pack(f.Length(), len(value), f.Encoder())
		if err != nil {
			return nil, "", err
		}
		value = padLeft + value + padRight
	}

	fieldEncode, err := f.Encoder().Encode(value)
	if err != nil {
		return nil, "", err
	}

	if f.Prefixer() != nil {
		fieldPrefix, err := f.Prefixer().EncodeLength(f.Encoder().GetLength())
		if err != nil {
			return nil, "", err
		}
		fieldEncode = append(fieldPrefix, fieldEncode...)
	}

	return fieldEncode, value, nil
}

func (f *Field) Length() int {
	return f.length
}

func (f *Field) Pattern() *regexp.Regexp {
	return f.pattern
}

func (f *Field) Encoder() encoding.Encoder {
	return f.encoder
}

func (f *Field) Padder() padding.Padder {
	return f.padder
}

func (f *Field) Prefixer() prefix.Prefixer {
	return f.prefixer
}

func (f *Field) Bitmap() *utils.BitSet {
	return nil
}

func (f *Field) SetBitmap(*utils.BitSet) {}

func (f *Field) GetType() FieldType {
	return f.Type
}
