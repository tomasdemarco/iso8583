package field

import (
	"errors"
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	pkgField "github.com/tomasdemarco/iso8583/packager/field"
	"github.com/tomasdemarco/iso8583/padding"
	"github.com/tomasdemarco/iso8583/prefix"
	"reflect"
	"regexp"
	"strings"
)

type Field interface{}

func Unpack(fieldPackager pkgField.Field, messageRaw []byte, position int) (*string, *int, error) {

	length, lengthPrefix, err := prefix.Unpack(fieldPackager.Prefix, messageRaw[position:])
	if err != nil {
		return nil, nil, err
	}

	if length == 0 {
		length = fieldPackager.Length
	}

	paddingLeft, paddingRight := padding.Unpack(fieldPackager.Padding, length)

	length += paddingLeft + paddingRight

	if fieldPackager.Encoding == encoding.Bcd ||
		fieldPackager.Encoding == encoding.Hex {
		length = length / 2
	}

	if len(messageRaw) < position+length+lengthPrefix {
		return nil, nil, errors.New("index out of range while trying to unpack")
	}

	value, err := encoding.Unpack(fieldPackager.Encoding, messageRaw[position+lengthPrefix:position+lengthPrefix+length])
	if err != nil {
		return nil, nil, err
	}

	value = value[paddingLeft : len(value)-paddingRight]

	match, _ := regexp.MatchString(fieldPackager.Pattern, value)
	if !match {
		err = errors.New("invalid format")
		return nil, nil, err
	}

	//if fieldPackager.SubFields != nil { //TODO ver como resolver subfields
	//	m.UnpackSubfields(field, value)
	//}

	length = length + lengthPrefix

	return &value, &length, nil
}

func Pack(fieldPackager pkgField.Field, value string) ([]byte, error) {

	length := len(value)
	if fieldPackager.Encoding == encoding.Binary {
		length = length / 2
	}

	fieldPrefix, err := prefix.Pack(fieldPackager.Prefix, length)
	if err != nil {
		return nil, err
	}

	padLeft, padRight := padding.Pack(fieldPackager.Padding, fieldPackager.Length, len(value))
	fieldEncode, err := encoding.Pack(fieldPackager.Encoding, padLeft+value+padRight)
	if err != nil {
		return nil, err
	}

	return append(fieldPrefix, fieldEncode...), nil
}

//func (f *Field) SetSubField(subFieldId string, value string) {
//
//	if f.SubFields == nil {
//		var subFields = make(map[string]string)
//		f.SubFields = subFields
//	}
//
//	f.SubFields[subFieldId] = value
//}
//
//func (f *Field) GetSubField(fieldId string, subFieldId string) (string, error) {
//	if value, ok := f.SubFields[subFieldId]; ok {
//		return value, nil
//	}
//
//	return "", errors.New(fmt.Sprintf(`the message does not contain with the id field "%s", the subfield with the id "%s"`, fieldId, subFieldId))
//}

func PackSubfield(fieldPkg pkgField.Field, value string) ([]byte, error) {

	length := len(value)
	if fieldPkg.Encoding == encoding.Binary {
		length = length / 2
	}

	fieldPrefix, err := prefix.Pack(fieldPkg.Prefix, length)
	if err != nil {
		return nil, err
	}

	padLeft, padRight := padding.Pack(fieldPkg.Padding, fieldPkg.Length, len(value))
	fieldEncode, err := encoding.Pack(fieldPkg.Encoding, padLeft+value+padRight)
	if err != nil {
		return nil, err
	}

	return append(fieldPrefix, fieldEncode...), nil
}

//func (f *Field) UnpackSubfields(subFields map[string]packager.SubField) error {
//	bitmapLength := subFields["00"].Length
//
//	f.SetSubField("00", f.Value[:subFields["00"].Length])
//
//	position := bitmapLength
//
//	firstBitmap, err := strconv.ParseUint(f.Value[:bitmapLength], 16, 32)
//	if err != nil {
//		return err
//	}
//
//	bitmapBinary := fmt.Sprintf("%0*b", bitmapLength, firstBitmap)
//
//	for i := 1; i < len(bitmapBinary); i++ {
//		if bitmapBinary[i-1:i] == "1" {
//			subfield := fmt.Sprintf("%02d", i)
//			f.SetSubField(subfield, f.Value[position:position+subFields[subfield].Length])
//			position += subFields[subfield].Length
//		}
//	}
//
//	return nil
//}

func ObjToFields[T any](obj T) (*map[string]Field, error) {
	var errorsArr []error

	objFields := reflect.ValueOf(&obj).Elem()

	isoFields := make(map[string]Field)
	for i := 0; i < objFields.NumField(); i++ {
		requiredField := objFields.Type().Field(i).Tag.Get("required")
		if strings.Contains(requiredField, "true") && objFields.Field(i).IsZero() {
			objField := objFields.Type().Field(i).Tag.Get("json")
			errorsArr = append(errorsArr, fmt.Errorf("%s is required", objField))
		}

		tagIsoField := objFields.Type().Field(i).Tag.Get("isoField")
		if tagIsoField != "" {
			if match, _ := regexp.MatchString("^[0-9]{1,3}$", tagIsoField); !match {
				objField := objFields.Type().Field(i).Tag.Get("json")
				errorsArr = append(errorsArr, fmt.Errorf("%s must match pattern ^[0-9]{1,3}$", objField))
			} else if objFields.Field(i).String() != "" {
				value := objFields.Field(i).String()

				isoFields[fmt.Sprintf("%03s", tagIsoField)] = Field(value)
			}
		}
	}

	return &isoFields, errors.Join(errorsArr...)
}
