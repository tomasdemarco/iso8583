package message

import (
	"errors"
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/field"
	"sort"
	"strconv"
)

func (m *Message) UnpackSubfields(field string, value string) {
	bitmapLength := m.Packager.Fields[field].SubFields["00"].Length
	m.SetSubField(field, "00", value[:m.Packager.Fields[field].SubFields["00"].Length])
	position := bitmapLength
	firstBitmap, err := strconv.ParseUint(value[:bitmapLength], 16, 32)
	if err != nil {
		return
	}
	bitmapBinary := fmt.Sprintf("%0*b", bitmapLength, firstBitmap)
	for i := 1; i < len(bitmapBinary); i++ {
		if bitmapBinary[i-1:i] == "1" {
			subfield := fmt.Sprintf("%02d", i)
			m.SetSubField(field, subfield, value[position:position+m.Packager.Fields[field].SubFields[subfield].Length])
			position += m.Packager.Fields[field].SubFields[subfield].Length
		}
	}
}

func (m *Message) PackSubfields(fieldsAux map[string]field.Field, field string) {
	var fieldResult string

	keys := make([]string, 0, len(fieldsAux[field].SubFields))
	for k := range fieldsAux[field].SubFields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i := range keys {
		m.SetSubField(field, keys[i], fieldsAux[field].SubFields[keys[i]])
		fieldResult += encoding.PackSubField(m.Packager.Fields[field].SubFields[keys[i]].Encoding, fieldsAux[field].SubFields[keys[i]])
	}

	m.SetField(field, fieldResult)
}

func (m *Message) SetSubField(fieldId string, subField string, value string) {

	if m.Fields == nil {
		var fields = make(map[string]field.Field)
		m.Fields = fields
	}
	fieldAux := m.Fields[fieldId]
	fieldAux.Value = ""
	m.Fields[fieldId] = fieldAux
	fieldAux = m.Fields[fieldId]
	if m.Fields[fieldId].SubFields == nil {
		var subFields = make(map[string]string)
		fieldAux.SubFields = subFields
	}
	fieldAux.SubFields[subField] = value
	m.Fields[fieldId] = fieldAux
}

func (m *Message) GetSubField(field string, subField string) (value string, err error) {
	if value, ok := m.Fields[field].SubFields[subField]; ok {
		return value, nil
	}
	err = errors.New("the message does not contain with the id field'" + field + "', the subfield with the id '" + subField + "'")
	return value, err
}
