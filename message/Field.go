package message

import "fmt"

type FieldBase interface {
	SetBytes([]byte)
	GetBytes() []byte
	GetString() string
	GetSubfields() map[string]Field
	PackSubfields(map[string]Field) error
	UnpackSubfields() map[string]Field
}

type Field struct {
	value     []byte
	subfields map[string]Field
}

func (f *Field) SetBytes(value []byte) {
	f.value = value
}

func (f *Field) GetBytes() []byte {
	return f.value
}

func (f *Field) GetString() string {
	return fmt.Sprintf("%X", f.value)
}

func (f *Field) GetSubfields() map[string]Field {
	return f.subfields
}

func (f *Field) PackSubfields(subfields map[string]Field) error {
	f.subfields = subfields
	return nil
}

func (f *Field) UnpackSubfields() map[string]Field {
	return f.subfields
}
