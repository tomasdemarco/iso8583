// Package message provides functionalities for packing and unpacking
// ISO 8583 messages.
package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tomasdemarco/iso8583/bitmap"
	"github.com/tomasdemarco/iso8583/header"
	"github.com/tomasdemarco/iso8583/packager"
	"github.com/tomasdemarco/iso8583/utils"
	"log"
	"strconv"
)

// Message represents an ISO 8583 message, containing its structure,
// fields, and the associated packager.
type Message struct {
	Packager *packager.Packager
	Length   int
	Header   header.Header
	Trailer  interface{}
	Bitmap   *utils.BitSet
	fields   map[int]Field // Cambiado a minúscula para ser interno
	// Almacena funciones de fábrica para StructField, permitiendo la creación dinámica.
	registeredStructFields map[int]func() Field
	TagsEmv                map[string]string
}

// NewMessage creates and returns a new Message instance
// initialized with the provided packager.
func NewMessage(packager *packager.Packager) *Message {
	return &Message{
		Packager:               packager,
		Bitmap:                 utils.NewBitSet(64, 128),
		fields:                 make(map[int]Field),
		registeredStructFields: make(map[int]func() Field),
	}
}

// SetField sets the value of a specific field in the message.
func (m *Message) SetField(id int, v any) {
	f, ok := m.fields[id]
	if !ok {
		f = m.createField(id)
	}

	f.Set(v)
	m.fields[id] = f
	m.Bitmap.Set(id)
}

// GetField retrieves the internal Field object for a given ID.
// This is primarily for internal use or advanced scenarios.
func (m *Message) GetField(id int) (Field, bool) {
	f, ok := m.fields[id]
	return f, ok
}

func (m *Message) Field(id int) *FieldAccessor {
	f, ok := m.fields[id] // Asumiendo que m.fields es tu mapa interno
	if !ok {
		return &FieldAccessor{
			err: fmt.Errorf("campo %d: %w", id, ErrNotFoundInMessage),
		}
	}
	return &FieldAccessor{field: f}
}

// String retrieves the value of a field as a string.
func (m *Message) String(id int) (string, error) {
	f, ok := m.GetField(id)
	if !ok {
		return "", fmt.Errorf("campo %d: %w", id, ErrNotFoundInMessage)
	}
	if sf, ok := f.Get().(*StringField); ok {
		return sf.Get().(string), nil
	}
	return "", fmt.Errorf("campo %d no es de tipo StringField", id)
}

// Int retrieves the value of a field as an int.
func (m *Message) Int(id int) (int, error) {
	f, ok := m.GetField(id)
	if !ok {
		return 0, fmt.Errorf("campo %d: %w", id, ErrNotFoundInMessage)
	}
	if ifld, ok := f.Get().(*IntField); ok {
		return ifld.Get().(int), nil
	}
	return 0, fmt.Errorf("campo %d no es de tipo IntField", id)
}

// Bytes retrieves the value of a field as a byte slice.
func (m *Message) Bytes(id int) ([]byte, error) {
	f, ok := m.GetField(id)
	if !ok {
		return nil, fmt.Errorf("campo %d: %w", id, ErrNotFoundInMessage)
	}
	if bf, ok := f.Get().(*BytesField); ok {
		return bf.Get().([]byte), nil
	}
	return nil, fmt.Errorf("campo %d no es de tipo BytesField", id)
}

// GetStruct retrieves the value of a field as a custom struct.
// T must be the struct type (e.g., CustomerInfo).
func GetStruct[T any](m *Message, id int) (T, error) {
	var result T // result es un valor cero del tipo T

	f, ok := m.GetField(id)
	if !ok {
		return result, fmt.Errorf("campo %d: %w", id, ErrNotFoundInMessage)
	}

	// Intentamos hacer un type assertion a *StructField[T]
	if sf, ok := f.Get().(*StructField[T]); ok {
		return sf.Get().(T), nil
	}

	return result, fmt.Errorf("campo %d no es de tipo StructField[%T]", id, result)
}

// RegisterField permite registrar manualmente una implementación de Field para un ID.
func (m *Message) RegisterField(id int, f Field) {
	m.fields[id] = f
}

// RegisterStructField registra un tipo de struct personalizado para un ID de campo.
// Esto permite que el mensaje sepa cómo crear e instanciar StructField[T]
// cuando se necesite para ese campo.
func RegisterStructField[T any](m *Message, id int) {
	m.registeredStructFields[id] = func() Field {
		return &StructField[T]{}
	}
}

// createField es una función interna para instanciar el tipo de Field correcto
// basado en la configuración del packager o en los tipos registrados.
func (m *Message) createField(id int) Field {
	if factory, ok := m.registeredStructFields[id]; ok {
		f := factory()
		m.fields[id] = f
		return f
	}

	var f Field
	f = &StringField{}
	fieldSpec, ok := m.Packager.Fields[id]
	if ok {
		switch fieldSpec.GetType() {
		case packager.Numeric:
			// Usamos StringField para NUMERIC por defecto para preservar ceros iniciales
			f = &StringField{}
		case packager.String: // Usar pkgfield.String para Alpha y AlphaNumeric
			f = &StringField{}
		case packager.Binary, packager.Bitmap:
			f = &BytesField{}
		default:
			f = &StringField{}
		}
	}

	m.fields[id] = f
	return f
}

// Unpack unpacks a byte slice of an ISO 8583 message
// into the Message structure, populating its fields.
// It returns an error if unpacking fails.
func (m *Message) Unpack(messageRaw []byte) (err error) {
	lengthMti, err := m.unpackMti(messageRaw)
	if err != nil {
		return err
	}

	lengthBitmap, err := m.unpackBitmap(messageRaw, lengthMti)
	if err != nil {
		return err
	}

	err = m.unpackFields(messageRaw, lengthMti+lengthBitmap)

	return err
}

// Pack packs the message fields into an ISO 8583 byte slice.
// It calculates the bitmap and encodes each field according to the packager's configuration.
// It returns the packed message as a byte slice, and an error if packing fails.
func (m *Message) Pack() ([]byte, error) {
	msgPacked := new(bytes.Buffer)
	encodeField, err := m.packMti()
	if err != nil {
		return nil, err
	}
	msgPacked.Write(encodeField)

	encodeField, err = m.packBitmap()
	if err != nil {
		return nil, err
	}
	msgPacked.Write(encodeField)

	encodeField, err = m.packFields()
	if err != nil {
		return nil, err
	}
	msgPacked.Write(encodeField)

	return msgPacked.Bytes(), nil
}

func (m *Message) Log() string {
	fieldsToLog := make(map[string]interface{})

	// Incluir MTI (campo 0) si está presente
	if field, ok := m.fields[0]; ok {
		fieldsToLog["0"] = field.Get()
	}

	// Incluir Bitmap (campo 1) si está presente
	if field, ok := m.fields[1]; ok {
		fieldsToLog["1"], _ = field.String()
	}

	if m.Bitmap != nil {
		for _, id := range m.Bitmap.GetSliceString() {
			// Los campos 0 y 1 ya se manejan explícitamente
			if id != 0 && id != 1 {
				if field, ok := m.fields[id]; ok {
					var err error
					fieldsToLog[strconv.Itoa(id)], err = field.Log()
					if err != nil {
						log.Println(err) //TODO ver si devolver el error
					}
				}
			}
		}
	}

	jsonBytes, err := json.Marshal(fieldsToLog)
	if err != nil {
		return fmt.Sprintf("{\"error\": \"no se pudo convertir el log a JSON: %v\"}", err)
	}

	return string(jsonBytes)
}

func (m *Message) packMti() ([]byte, error) {
	if fldPKg, ok := m.Packager.Fields[0]; ok {
		fld, err := m.Field(0).String()
		if err != nil {
			return nil, fmt.Errorf("pack mti: %w", err)
		}

		encodeField, _, err := fldPKg.Pack(fld)
		if err != nil {
			return nil, fmt.Errorf("pack mti: %w", err)
		}

		return encodeField, nil
	}

	return nil, ErrMTINotFoundInPackager
}

func (m *Message) packBitmap() ([]byte, error) {
	if fldPKg, ok := m.Packager.Fields[1]; ok {
		// El bitmap se maneja como BytesField, su String() devuelve hex
		if len(m.Bitmap.ToBytes()) > fldPKg.Length() {
			// Si el bitmap es secundario, SetField(1, ...) lo actualizará
			m.SetField(1, m.Bitmap.ToBytes())
		}

		encodeField, _, errPack := fldPKg.Pack(m.Bitmap.ToString())
		if errPack != nil {
			return nil, fmt.Errorf("pack bitmap: %w", errPack)
		}

		return encodeField, nil
	}
	return nil, ErrBitmapNotFoundInPackager
}

func (m *Message) packFields() ([]byte, error) {
	fieldsPacked := new(bytes.Buffer)

	for _, k := range m.Bitmap.GetSliceString() {
		if k != 0 && k != 1 { // MTI y Bitmap ya empaquetados
			if fldPkg, ok := m.Packager.Fields[k]; ok {
				fld, err := m.Field(k).String()
				if err != nil {
					return nil, fmt.Errorf("pack field %d: %w", k, err)
				}
				encodeField, _, errPack := fldPkg.Pack(fld)
				if errPack != nil {
					return nil, fmt.Errorf("pack field %d: %w", k, errPack)
				}

				// Actualizamos el campo con la versión "plain" si el packager la modifica (ej. padding)
				// m.SetField(k, plainField) // Esto podría ser problemático si plainField no es el tipo original
				fieldsPacked.Write(encodeField)
			} else {
				return nil, fmt.Errorf("field %d: %w", k, ErrNotFoundInPackager)
			}
		}
	}

	return fieldsPacked.Bytes(), nil
}

// unpackMti unpacks the Message FieldType Indicator (MTI) from the message.
// This is an internal helper method.
func (m *Message) unpackMti(messageRaw []byte) (int, error) {
	if fldPkg, ok := m.Packager.Fields[0]; ok {
		value, length, err := fldPkg.Unpack(messageRaw, 0)
		if err != nil {
			return 0, fmt.Errorf("unpack MTI: %w", err)
		}

		// Usamos SetField con el valor string, que creará un StringField por defecto
		m.SetField(0, value)

		return length, nil
	}

	return 0, ErrMTINotFoundInPackager
}

// unpackBitmap unpacks the bitmap from the message.
// This is an internal helper method.
func (m *Message) unpackBitmap(messageRaw []byte, offset int) (int, error) {
	if fldPkg, ok := m.Packager.Fields[1]; ok {
		bMap, length, err := bitmap.Unpack(fldPkg, messageRaw, offset)
		if err != nil {
			return 0, fmt.Errorf("unpack bitmap: %w", err)
		}

		m.Bitmap = bMap
		// Usamos SetField con los bytes del bitmap, que creará un BytesField por defecto
		m.SetField(1, bMap.ToBytes())

		return length, nil
	}

	return 0, ErrBitmapNotFoundInPackager
}

func (m *Message) unpackFields(messageRaw []byte, position int) error {
	for _, fieldId := range m.Bitmap.GetSliceString() {
		if fieldId != 0 && fieldId != 1 { // MTI y Bitmap ya desempaquetados
			if fldPkg, ok := m.Packager.Fields[fieldId]; ok {
				value, length, err := fldPkg.Unpack(messageRaw, position)
				if err != nil {
					return fmt.Errorf("unpack field %d: %w", fieldId, err)
				}

				// Obtenemos el Field existente o creamos uno por defecto
				f, ok := m.fields[fieldId]
				if !ok {
					f = m.createField(fieldId)
				}
				// Pasamos los bytes (convertidos a string) al SetBytes del Field
				f.SetBytes([]byte(value))
				position += length
			} else {
				return fmt.Errorf("field %d: %w", fieldId, ErrNotFoundInPackager)
			}
		}
	}
	return nil
}
