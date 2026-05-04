package message

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/field"
	"github.com/tomasdemarco/iso8583/packager"
	"github.com/tomasdemarco/iso8583/prefix"
)

// --- Test Setup ---

type CustomerInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *CustomerInfo) Pack() (string, error) {
	return fmt.Sprintf("%s|%s", c.ID, c.Name), nil
}

func (c *CustomerInfo) Unpack(data string) error {
	parts := strings.SplitN(data, "|", 2)
	if len(parts) != 2 {
		return fmt.Errorf("formato inválido para CustomerInfo: se esperaban 2 partes, se obtuvieron %d", len(parts))
	}
	c.ID = parts[0]
	c.Name = parts[1]
	return nil
}

func (c *CustomerInfo) Log() (interface{}, error) {
	return c, nil
}

var _ field.CustomPacker = (*CustomerInfo)(nil)

// --- Pruebas ---

func TestNewMessage(t *testing.T) {
	pkg, err := packager.LoadFromJson("../packager/json_examples", "iso87BPackager.json")
	if err != nil {
		t.Fatalf("Error al cargar el packager: %v", err)
	}
	msg := NewMessage(pkg)

	if msg == nil {
		t.Fatal("NewMessage devolvió nil")
	}
	if msg.Packager == nil {
		t.Error("FieldPackager no inicializado")
	}
	if msg.Bitmap == nil {
		t.Error("Bitmap no inicializado")
	}
	if msg.fields == nil {
		t.Error("fields map no inicializado")
	}
}

func TestSetAndGetBasicFields(t *testing.T) {
	// Create a programmatic packager to test specific types
	progPkg := &packager.Packager{
		Description: "Programmatic Packager",
		Fields:      make(map[int]packager.FieldPackager),
	}
	progPkg.Fields[3] = packager.NewField("Processing Code", 6, encoding.NewAsciiEncoder(), packager.WithFieldType(packager.Numeric))
	progPkg.Fields[11] = packager.NewField("STAN", 6, encoding.NewAsciiEncoder(), packager.WithFieldType(packager.Numeric), packager.WithCustomType(func() field.Field { return &field.Int{} }))
	progPkg.Fields[126] = packager.NewField("Reserved", 999, encoding.NewBinaryEncoder(), packager.WithFieldType(packager.Binary), packager.WithCustomType(func() field.Field { return &field.Bytes{} }))

	msg := NewMessage(progPkg)

	// Test StringField
	if err := msg.SetFieldString(3, "000000"); err != nil {
		t.Fatalf("SetFieldString(3, ...) falló: %v", err)
	}
	valStr, err := msg.GetFieldString(3)
	if err != nil {
		t.Fatalf("GetFieldString(3) falló: %v", err)
	}
	if valStr != "000000" {
		t.Errorf("GetFieldString(3) esperado \"000000\", obtenido \"%s\"", valStr)
	}

	// Test IntField
	if err := msg.SetFieldInt(11, 123456); err != nil {
		t.Fatalf("SetFieldInt(11, ...) falló: %v", err)
	}
	valInt, err := msg.GetFieldInt(11)
	if err != nil {
		t.Fatalf("GetFieldInt(11) falló: %v", err)
	}
	if valInt != 123456 {
		t.Errorf("GetFieldInt(11) esperado 123456, obtenido %d", valInt)
	}

	// Test BytesField
	testBytes := []byte{0x01, 0x02, 0x03}
	if err := msg.SetFieldBytes(126, testBytes); err != nil {
		t.Fatalf("SetFieldBytes(126, ...) falló: %v", err)
	}
	valBytes, err := msg.GetFieldBytes(126)
	if err != nil {
		t.Fatalf("GetFieldBytes(126) falló: %v", err)
	}
	if !reflect.DeepEqual(valBytes, testBytes) {
		t.Errorf("GetFieldBytes(126) esperado %v, obtenido %v", testBytes, valBytes)
	}

	// Verificar que el bitmap se actualizó
	if !msg.Bitmap.Get(3) || !msg.Bitmap.Get(11) || !msg.Bitmap.Get(126) {
		t.Error("Bitmap no se actualizó correctamente para los campos establecidos")
	}
}

func TestSetAndGetStructField(t *testing.T) {
	pkg, err := packager.LoadFromJson("../packager/json_examples", "iso87BPackager.json")
	if err != nil {
		t.Fatalf("Error al cargar el packager: %v", err)
	}

	// Create a new field for the custom type and add it to the packager
	customField := packager.NewField(
		"Custom Data",
		999,
		encoding.NewAsciiEncoder(),
		packager.WithPrefix(prefix.NewAsciiPrefixer(int(prefix.LL), false, false)),
		packager.WithCustomType(func() field.Field {
			return &field.Struct[CustomerInfo]{}
		}),
	)
	pkg.Fields[48] = customField

	msg := NewMessage(pkg)

	originalInfo := CustomerInfo{ID: "C123", Name: "Juan Perez"}
	if err := SetFieldStruct(msg, 48, originalInfo); err != nil {
		t.Fatalf("SetFieldStruct(48, ...) falló: %v", err)
	}

	retrievedInfo, err := GetFieldStruct[CustomerInfo](msg, 48)
	if err != nil {
		t.Fatalf("GetFieldStruct(48) falló: %v", err)
	}

	if !reflect.DeepEqual(retrievedInfo, originalInfo) {
		t.Errorf("Struct recuperado no coincide. Esperado: %+v, Obtenido: %+v", originalInfo, retrievedInfo)
	}
}

func TestPackUnpackFullMessage(t *testing.T) {
	pkg, err := packager.LoadFromJson("../packager/json_examples", "iso87BPackager.json")
	if err != nil {
		t.Fatalf("Error al cargar el packager: %v", err)
	}

	// Programmatically add custom types to the packager
	pkg.Fields[11] = packager.NewField("STAN", 6, encoding.NewAsciiEncoder(), packager.WithFieldType(packager.Numeric), packager.WithCustomType(func() field.Field { return &field.Int{} }))
	pkg.Fields[48] = packager.NewField("Custom", 999, encoding.NewAsciiEncoder(), packager.WithPrefix(prefix.NewAsciiPrefixer(int(prefix.LLL), false, false)), packager.WithCustomType(func() field.Field { return &field.Struct[CustomerInfo]{} }))
	pkg.Fields[62] = packager.NewField("Bytes", 999, encoding.NewBinaryEncoder(), packager.WithPrefix(prefix.NewAsciiPrefixer(int(prefix.LLL), false, false)), packager.WithCustomType(func() field.Field { return &field.Bytes{} }))

	msgToPack := NewMessage(pkg)

	// Setear campos
	_ = msgToPack.SetFieldString(0, "0200")
	_ = msgToPack.SetFieldString(3, "000000")
	_ = msgToPack.SetFieldInt(11, 123456)
	originalInfo := CustomerInfo{ID: "C987", Name: "Maria Garcia"}
	_ = SetFieldStruct(msgToPack, 48, originalInfo)
	testBytes := []byte("aabbcc")
	_ = msgToPack.SetFieldBytes(62, testBytes)

	// Empaquetar
	packedData, err := msgToPack.Pack()
	if err != nil {
		t.Fatalf("Parse() falló: %v", err)
	}
	t.Logf("Mensaje Empaquetado: %X", packedData)

	// Desempaquetar
	msgToUnpack := NewMessage(pkg)
	if err := msgToUnpack.Unpack(packedData); err != nil {
		t.Fatalf("Unparse() falló: %v", err)
	}

	// Verificar campos
	mti, _ := msgToUnpack.GetFieldString(0)
	if mti != "0200" {
		t.Errorf("MTI esperado \"0200\", obtenido \"%s\"", mti)
	}

	procCode, _ := msgToUnpack.GetFieldString(3)
	if procCode != "000000" {
		t.Errorf("Processing Code esperado \"000000\", obtenido \"%s\"", procCode)
	}

	stan, _ := msgToUnpack.GetFieldInt(11)
	if stan != 123456 {
		t.Errorf("STAN esperado 123456, obtenido %d", stan)
	}

	retrievedInfo, _ := GetFieldStruct[CustomerInfo](msgToUnpack, 48)
	if !reflect.DeepEqual(retrievedInfo, originalInfo) {
		t.Errorf("Struct recuperado no coincide. Esperado: %+v, Obtenido: %+v", originalInfo, retrievedInfo)
	}

	retrievedBytes, _ := msgToUnpack.GetFieldBytes(62)
	if !reflect.DeepEqual(retrievedBytes, testBytes) {
		t.Errorf("Bytes(62) esperado %v, obtenido %v", testBytes, retrievedBytes)
	}
}

func TestLogMethod(t *testing.T) {
	pkg, err := packager.LoadFromJson("../packager/json_examples", "iso87BPackager.json")
	if err != nil {
		t.Fatalf("Error al cargar el packager: %v", err)
	}
	pkg.Fields[11] = packager.NewField("STAN", 6, encoding.NewAsciiEncoder(), packager.WithFieldType(packager.Numeric), packager.WithCustomType(func() field.Field { return &field.Int{} }))
	pkg.Fields[48] = packager.NewField("Custom", 999, encoding.NewAsciiEncoder(), packager.WithPrefix(prefix.NewAsciiPrefixer(int(prefix.LLL), false, false)), packager.WithCustomType(func() field.Field { return &field.Struct[CustomerInfo]{} }))
	pkg.Fields[62] = packager.NewField("Bytes", 999, encoding.NewAsciiEncoder(), packager.WithPrefix(prefix.NewAsciiPrefixer(int(prefix.LLL), false, false)), packager.WithCustomType(func() field.Field { return &field.Bytes{} }))

	msg := NewMessage(pkg)

	_ = msg.SetFieldString(0, "0200")
	_ = msg.SetFieldString(3, "000000")
	_ = msg.SetFieldInt(11, 123456)
	_ = SetFieldStruct(msg, 48, CustomerInfo{ID: "C123", Name: "LogMsg Test"})
	_ = msg.SetFieldBytes(62, []byte("aabbcc"))

	logOutput := msg.LogMsg()
	t.Logf("LogMsg Output:\n%s", logOutput)

	if !strings.Contains(logOutput, "\"0\":\"0200\"") {
		t.Error("LogMsg no contiene MTI")
	}
	if !strings.Contains(logOutput, "\"3\":\"000000\"") {
		t.Error("LogMsg no contiene Processing Code")
	}
	if !strings.Contains(logOutput, "\"11\":\"123456\"") {
		t.Error("LogMsg no contiene STAN")
	}
	if !strings.Contains(logOutput, "\"id\":\"C123\"") {
		t.Error("LogMsg no contiene CustomerInfo serializado correctamente")
	}
	if !strings.Contains(logOutput, "\"62\":\"616162626363\"") {
		t.Error("LogMsg no contiene BytesField serializado correctamente")
	}
}

func TestTypeSafety(t *testing.T) {
	pkg := &packager.Packager{
		Fields: make(map[int]packager.FieldPackager),
	}
	// Field 3 is a String
	pkg.Fields[3] = packager.NewField("String Field", 10, encoding.NewAsciiEncoder())
	// Field 11 is an Int
	pkg.Fields[11] = packager.NewField("Int Field", 6, encoding.NewAsciiEncoder(), packager.WithCustomType(func() field.Field { return &field.Int{} }))

	msg := NewMessage(pkg)

	// Try to set int on a string field
	err := msg.SetFieldInt(3, 123)
	if err == nil {
		t.Errorf("SetFieldInt(3) debería haber fallado para un StringField, pero no lo hizo")
	} else if !strings.Contains(err.Error(), "invalid type") {
		t.Errorf("SetFieldInt(3) falló con un error inesperado: %v", err)
	}

	// Try to set string on an int field
	err = msg.SetFieldString(11, "not-an-int")
	if err == nil {
		t.Errorf("SetFieldString(11) debería haber fallado para un IntField, pero no lo hizo")
	} else if !strings.Contains(err.Error(), "invalid type") {
		t.Errorf("SetFieldString(11) falló con un error inesperado: %v", err)
	}
}
