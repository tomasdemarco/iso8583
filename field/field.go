package field

import (
	"errors"
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/packager"
	"github.com/tomasdemarco/iso8583/padding"
	"github.com/tomasdemarco/iso8583/prefix"
	"regexp"
	"sort"
	"strconv"
)

type Field struct {
	Value     string
	SubFields map[string]string
}

func Unpack(fieldPackager packager.Field, messageRaw string, position int, field string) (*string, *int, error) {

	if len(messageRaw) < position+(fieldPackager.Prefix.Type.EnumIndex()*1) {
		err := errors.New("index out of range while trying to unpack prefix field " + field)
		return nil, nil, err
	}

	length, lengthPrefix, err := prefix.Unpack(fieldPackager.Prefix, messageRaw[position:])
	if err != nil {
		return nil, nil, err
	}

	if length == 0 {
		length = fieldPackager.Length
	}

	if len(messageRaw) < position+length+lengthPrefix {
		err = errors.New("index out of range while trying to unpack field " + field)
		return nil, nil, err
	}

	paddingRight, paddingLeft := padding.Unpack(fieldPackager.Padding, length)
	fmt.Println(paddingRight, paddingLeft)
	value, doubleLength, err := encoding.Unpack(fieldPackager.Encoding, messageRaw, field, position+lengthPrefix+paddingLeft, length)
	if err != nil {
		return nil, nil, err
	}

	match, _ := regexp.MatchString(fieldPackager.Pattern, value)
	if !match {
		err = errors.New("invalid format in field " + field)
		return nil, nil, err
	}

	//if fieldPackager.SubFields != nil { //TODO ver como resolver subfields
	//	m.UnpackSubfields(field, value)
	//}

	length = length*doubleLength + lengthPrefix + paddingRight + paddingLeft

	return &value, &length, nil
}

func Pack(fieldPackager packager.Field, value string) (string, error) {
	fieldPrefix, err := prefix.Pack(fieldPackager.Prefix, len(value))
	if err != nil {
		return "", err
	}
	padRight, padLeft := padding.Pack(fieldPackager.Padding, fieldPackager.Length, len(value))

	fieldEncode := encoding.Pack(fieldPackager.Encoding, padLeft+value+padRight)

	return fieldPrefix + fieldEncode, nil
}

func (f *Field) SetSubField(subFieldId string, value string) {

	if f.SubFields == nil {
		var subFields = make(map[string]string)
		f.SubFields = subFields
	}

	f.SubFields[subFieldId] = value
}

func (f *Field) GetSubField(fieldId string, subFieldId string) (string, error) {
	if value, ok := f.SubFields[subFieldId]; ok {
		return value, nil
	}

	return "", errors.New(fmt.Sprintf(`the message does not contain with the id field "%s", the subfield with the id "%s"`, fieldId, subFieldId))
}

func (f *Field) UnpackSubfields(subFields map[string]packager.SubField) error {
	bitmapLength := subFields["00"].Length

	f.SetSubField("00", f.Value[:subFields["00"].Length])

	position := bitmapLength

	firstBitmap, err := strconv.ParseUint(f.Value[:bitmapLength], 16, 32)
	if err != nil {
		return err
	}

	bitmapBinary := fmt.Sprintf("%0*b", bitmapLength, firstBitmap)

	for i := 1; i < len(bitmapBinary); i++ {
		if bitmapBinary[i-1:i] == "1" {
			subfield := fmt.Sprintf("%02d", i)
			f.SetSubField(subfield, f.Value[position:position+subFields[subfield].Length])
			position += subFields[subfield].Length
		}
	}

	return nil
}

func (f *Field) PackSubfields(subFields map[string]packager.SubField) string {
	var fieldResult string

	keys := make([]string, 0, len(f.SubFields))
	for k := range f.SubFields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i := range keys {
		f.SetSubField(keys[i], f.SubFields[keys[i]])
		fieldResult += encoding.PackSubField(subFields[keys[i]].Encoding, f.SubFields[keys[i]])
	}

	return fieldResult
}
