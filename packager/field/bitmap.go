package field

//
//import (
//	"errors"
//	"github.com/tomasdemarco/iso8583/bitmap"
//	"github.com/tomasdemarco/iso8583/encoding"
//	"github.com/tomasdemarco/iso8583/padding"
//	"github.com/tomasdemarco/iso8583/prefix"
//	"github.com/tomasdemarco/iso8583/utils"
//	"regexp"
//)
//
//// Field represents an ISO 8583 field's definition, including its data type,
//// length, validation pattern, encoding, prefix, and padding rules.
//type BitmapField struct {
//	Description     string
//	length          int
//	pattern         *regexp.Regexp
//	encoder         encoding.Encoder
//	bitZeroExtended bool
//}
//
//func NewBitmapField(
//	description string,
//	length int,
//	encoding encoding.Encoder,
//	bitZeroExtended bool,
//) Packager {
//	return &BitmapField{
//		Description:     description,
//		length:          length,
//		encoder:         encoding,
//		bitZeroExtended: bitZeroExtended,
//	}
//}
//
//// Unpack unpacks a field's value from a raw message byte slice.
//// It handles length prefixes, padding, and decoding based on the field's configuration.
//// It returns the unpacked field value as a string, the total length consumed in bytes,
//// and an error if unpacking fails.
//func (f BitmapField) Unpack(messageRaw []byte, position int) (string, int, error) {
//	bitmapRaw, length, err := bitmap.Unpack(messageRaw, position, f.Length(), f.bitZeroExtended)
//	if err != nil {
//		return "", 0, err
//	}
//
//	f.Encoder().SetLength(length)
//
//	if len(messageRaw) < position+f.Length() {
//		return "", 0, errors.New("index out of range while trying to unpack")
//	}
//
//	value, err := f.Encoder().Decode(messageRaw[position:])
//	if err != nil {
//		return "", 0, err
//	}
//
//	return value, length, nil
//}
//
//// Pack packs a field's string value into a byte slice according to its configuration.
//// It applies padding, encodes the value, and prepends the length prefix if defined.
//// It returns the packed field as a byte slice, the plain (padded) field string,
//// and an error if packing fails.
//func (f BitmapField) Pack(value string) ([]byte, string, error) {
//	fieldEncode, err := f.Encoder().Encode(value)
//	if err != nil {
//		return nil, "", err
//	}
//	length := len(fieldEncode)
//
//	if _, ok := f.Encoder().(*encoding.BCD); ok {
//		length = length * 2
//	}
//
//	return fieldEncode, "", nil
//}
//
//func (f BitmapField) Length() int {
//	return f.length
//}
//
//func (f BitmapField) Pattern() *regexp.Regexp {
//	return f.pattern
//}
//
//func (f BitmapField) Encoder() encoding.Encoder {
//	return f.encoder
//}
//
//func (f BitmapField) Padder() padding.Padder {
//	return nil
//}
//
//func (f BitmapField) Prefixer() prefix.Prefixer {
//	return nil
//}
//
//func (f BitmapField) Bitmap() *utils.BitSet {
//	return nil
//}
//
//func (f BitmapField) SetBitmap(bmap *utils.BitSet) {
//	f.bitmap = bmap
//}
