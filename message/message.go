package message

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tomasdemarco/iso8583/bitmap"
	"github.com/tomasdemarco/iso8583/packager"
	"regexp"
)

type Message struct {
	Packager *packager.Packager
	Length   int
	Header   interface{}
	Trailer  interface{}
	Bitmap   []string
	Fields   map[string]string
	TagsEmv  map[string]string
}

func NewMessage(packager *packager.Packager) *Message {
	return &Message{
		Packager: packager,
	}
}

func (m *Message) Unpack(messageRaw []byte) (err error) {
	if _, ok := m.Packager.Fields["000"]; !ok {
		err = errors.New("packager does not contain field 000")
		return err
	}

	value, length, err := m.Packager.Fields["000"].Unpack(messageRaw, 0)
	if err != nil {
		return errors.New(fmt.Sprintf("pack field 000: %v", err))
	}

	m.SetField("000", value)
	position := length

	match, _ := regexp.MatchString(m.Packager.Fields["000"].Pattern, m.Fields["000"])
	if !match {
		err = errors.New("invalid format in field 000")
		return err
	}

	if _, ok := m.Packager.Fields["001"]; !ok {
		err = errors.New("packager does not contain field 001")
		return err
	}

	if len(messageRaw) < position+m.Packager.Fields["001"].Length {
		err = errors.New("the message is too short to be unpacked")
		return err
	}

	lengthBitmap, sliceBitmap, err := bitmap.Unpack(m.Packager.Fields["001"], messageRaw, position)
	if err != nil {
		err = errors.New("could not get bitmap, " + err.Error())
		return err
	}

	m.Bitmap = sliceBitmap

	m.Packager.Fields["001"].Encoding.SetLength(*lengthBitmap)
	value, err = m.Packager.Fields["001"].Encoding.Decode(messageRaw[position:])
	if err != nil {
		return errors.New(fmt.Sprintf("unpack field 001: %v", err))
	}

	m.SetField("001", value)
	position += *lengthBitmap

	match, _ = regexp.MatchString(m.Packager.Fields["001"].Pattern, m.Fields["001"])
	if !match {
		err = errors.New("invalid format in field 001")
		return err
	}

	for _, fieldId := range m.Bitmap {
		if fieldId != "001" {
			value, length, err = m.Packager.Fields[fieldId].Unpack(messageRaw, position)
			if err != nil {
				return errors.New(fmt.Sprintf("unpack field %s: %v", fieldId, err))
			}

			m.SetField(fieldId, value)
			position += length
		}
	}

	return nil
}

func (m *Message) Pack() ([]byte, error) {

	bitmapSlice, bitmapString, err := bitmap.Pack(m.Fields)
	if err != nil {
		return nil, err
	}

	m.SetField("001", fmt.Sprintf("%x", bitmapString))

	m.Bitmap = append(bitmapSlice[:1], append([]string{"001"}, bitmapSlice[1:]...)...)

	msgPacked := new(bytes.Buffer)
	for _, k := range m.Bitmap {
		var value string

		//if m.Fields[k].Subfields != nil {
		//	subfields := m.Fields[k].Subfields
		//	val, err := subfields.Pack(m.Packager.Fields[k])
		//	if err != nil {
		//		return nil, err
		//	}
		//	value = val
		//} else {
		//}
		value = m.Fields[k]

		encodeField, plainField, errPack := m.Packager.Fields[k].Pack(value)
		if errPack != nil {
			return nil, errors.New(fmt.Sprintf("pack field %s: %v", k, errPack))
		}
		m.SetField(k, plainField)

		msgPacked.Write(encodeField)
	}

	return msgPacked.Bytes(), err
}

func (m *Message) SetField(fieldId string, value string) {
	if m.Fields == nil {
		var fields = make(map[string]string)
		m.Fields = fields
	}
	fieldAux := m.Fields[fieldId]
	fieldAux = value
	m.Fields[fieldId] = fieldAux
}

func (m *Message) GetField(fieldId string) (string, error) {
	if fld, ok := m.Fields[fieldId]; ok {
		return fld, nil
	}

	return "", errors.New(fmt.Sprintf(`the message does not contain the field with the id "%s"`, fieldId))
}

//func (m *Message) SetSubfield(fieldId string, subfieldId string, value string) {
//	fields := m.Fields
//	if fields == nil {
//		fields = make(map[string]field.Field)
//	}
//
//	if fields[fieldId].Subfields == nil {
//		var subfields = make(subfield.Subfields)
//		subfields[subfieldId] = value
//		fieldAux := m.Fields[fieldId]
//		fieldAux.Subfields = subfields
//		m.Fields[fieldId] = fieldAux
//	}
//}
//
//func (m *Message) SetSubfields(fieldId string, value subfield.Subfields) {
//	if m.Fields == nil {
//		var fields = make(map[string]field.Field)
//		m.Fields = fields
//	}
//	fieldAux := m.Fields[fieldId]
//	fieldAux.Subfields = value
//	m.Fields[fieldId] = fieldAux
//}
//
//func (m *Message) GetSubfields(fieldId string) (subfield.Subfields, error) {
//	if fld, ok := m.Fields[fieldId]; ok {
//		return fld.Subfields, nil
//	}
//
//	return nil, errors.New(fmt.Sprintf(`the message does not contain subfields in the field with the id "%s"`, fieldId))
//}

//func Build[T any](obj T) (*Message, error) {
//	var errorsArr []error
//
//	objFields := reflect.ValueOf(&obj).Elem()
//
//	isoFields := make(map[string]field.Field)
//	for i := 0; i < objFields.NumField(); i++ {
//		requiredField := objFields.Type().Field(i).Tag.Get("required")
//		if strings.Contains(requiredField, "true") && objFields.Field(i).IsZero() {
//			objField := objFields.Type().Field(i).Tag.Get("json")
//			errorsArr = append(errorsArr, fmt.Errorf("%s is required", objField))
//		}
//
//		tagIsoField := objFields.Type().Field(i).Tag.Get("isoField")
//		if tagIsoField != "" && objFields.Elem().Field(i).String() != "" {
//			if match, _ := regexp.MatchString("^[0-9]{1,3}$", tagIsoField); !match {
//				objField := objFields.Type().Field(i).Tag.Get("json")
//				errorsArr = append(errorsArr, fmt.Errorf("%s must match pattern ^[0-9]{1,3}$", objField))
//			} else {
//				value := objFields.Elem().Field(i).String()
//				isoFields[fmt.Sprintf("%03s", tagIsoField)] = value
//			}
//		}
//	}
//	return nil, errors.Join(errorsArr...)
//}
