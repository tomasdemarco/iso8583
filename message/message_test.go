package message

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/tomasdemarco/iso8583/packager"
)

// --- Test Setup ---

// CustomerInfo: Un struct personalizado que implementa CustomPacker
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

// --- Pruebas ---

func TestNewMessage(t *testing.T) {
	pkg, err := packager.LoadFromJson("../packager", "iso87BPackager.json")
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
	if msg.registeredStructFields == nil {
		t.Error("registeredStructFields map no inicializado")
	}
}

func TestSetAndGetBasicFields(t *testing.T) {
	pkg, err := packager.LoadFromJson("../packager", "iso87BPackager.json")
	if err != nil {
		t.Fatalf("Error al cargar el packager: %v", err)
	}
	msg := NewMessage(pkg)

	// Test StringField (por defecto para campos numéricos y alfanuméricos)
	msg.SetField(3, "000000")

	valStr, err := msg.String(3)
	if err != nil {
		t.Fatalf("String(3) falló: %v", err)
	}
	if valStr != "000000" {
		t.Errorf("String(3) esperado \"000000\", obtenido \"%s\"", valStr)
	}

	// Test IntField (registrado manualmente)
	msg.RegisterField(11, &IntField{})
	msg.SetField(11, 123456)

	valInt, err := msg.Int(11)
	if err != nil {
		t.Fatalf("Int(11) falló: %v", err)
	}
	if valInt != 123456 {
		t.Errorf("Int(11) esperado 123456, obtenido %d", valInt)
	}

	// Test BytesField (ahora definido en el packager)
	testBytes := []byte{0x01, 0x02, 0x03}
	msg.SetField(126, testBytes)

	valBytes, err := msg.Bytes(126)
	if err != nil {
		t.Fatalf("Bytes(126) falló: %v", err)
	}
	if !reflect.DeepEqual(valBytes, testBytes) {
		t.Errorf("Bytes(126) esperado %v, obtenido %v", testBytes, valBytes)
	}

	// Verificar que el bitmap se actualizó
	if !msg.Bitmap.Get(3) || !msg.Bitmap.Get(11) || !msg.Bitmap.Get(126) {
		t.Error("Bitmap no se actualizó correctamente para los campos establecidos")
	}
}

func TestSetAndGetStructField(t *testing.T) {
	pkg, err := packager.LoadFromJson("../packager", "iso87BPackager.json")
	if err != nil {
		t.Fatalf("Error al cargar el packager: %v", err)
	}
	msg := NewMessage(pkg)

	// Registrar el StructField para el campo 48
	RegisterStructField[CustomerInfo](msg, 48)

	originalInfo := CustomerInfo{ID: "C123", Name: "Juan Perez"}
	msg.SetField(48, originalInfo)

	retrievedInfo, err := GetStruct[CustomerInfo](msg, 48)
	if err != nil {
		t.Fatalf("GetStruct(48) falló: %v", err)
	}

	if !reflect.DeepEqual(retrievedInfo, originalInfo) {
		t.Errorf("Struct recuperado no coincide. Esperado: %+v, Obtenido: %+v", originalInfo, retrievedInfo)
	}
}

func TestPackUnpackFullMessage(t *testing.T) {
	pkg, err := packager.LoadFromJson("../packager", "iso87BPackager.json")
	if err != nil {
		t.Fatalf("Error al cargar el packager: %v", err)
	}
	msgToPack := NewMessage(pkg)

	// Registrar StructField para el campo 48
	RegisterStructField[CustomerInfo](msgToPack, 48)

	// Setear campos de diferentes tipos
	msgToPack.SetField(0, "0200")
	msgToPack.SetField(3, "000000")

	msgToPack.RegisterField(11, &IntField{}) // STAN
	msgToPack.SetField(11, 123456)

	originalInfo := CustomerInfo{ID: "C987", Name: "Maria Garcia"}
	msgToPack.SetField(48, originalInfo)

	testBytes := "aabbcc"
	msgToPack.RegisterField(62, &StringField{}) // Registrar BytesField para campo 62
	msgToPack.SetField(62, testBytes)

	// Empaquetar
	packedData, err := msgToPack.Pack()
	if err != nil {
		t.Fatalf("Pack() falló: %v", err)
	}
	t.Logf("Mensaje Empaquetado: %X", packedData)

	// Desempaquetar
	msgToUnpack := NewMessage(pkg)
	RegisterStructField[CustomerInfo](msgToUnpack, 48) // Registrar también en el mensaje de desempaque
	if err := msgToUnpack.Unpack(packedData); err != nil {
		t.Fatalf("Unpack() falló: %v", err)
	}

	// Verificar campos
	mti, err := msgToUnpack.String(0)
	if err != nil {
		t.Fatalf("String(0) falló: %v", err)
	}
	if mti != "0200" {
		t.Errorf("MTI esperado \"0200\", obtenido \"%s\"", mti)
	}

	procCode, err := msgToUnpack.String(3)
	if err != nil {
		t.Fatalf("String(3) falló: %v", err)
	}
	if procCode != "000000" {
		t.Errorf("Processing Code esperado \"000000\", obtenido \"%s\"", procCode)
	}

	stan, err := msgToUnpack.String(11)
	if err != nil {
		t.Fatalf("Int(11) falló: %v", err)
	}
	if stan != "123456" {
		t.Errorf("STAN esperado 123456, obtenido %s", stan)
	}

	retrievedInfo, err := GetStruct[CustomerInfo](msgToUnpack, 48)
	if err != nil {
		t.Fatalf("GetStruct(48) falló: %v", err)
	}
	if !reflect.DeepEqual(retrievedInfo, originalInfo) {
		t.Errorf("Struct recuperado no coincide. Esperado: %+v, Obtenido: %+v", originalInfo, retrievedInfo)
	}

	retrievedBytes, err := msgToUnpack.String(62)
	if err != nil {
		t.Fatalf("Bytes(62) falló: %v", err)
	}
	if !reflect.DeepEqual(retrievedBytes, testBytes) {
		t.Errorf("Bytes(62) esperado %v, obtenido %v", testBytes, retrievedBytes)
	}
}

func TestLogMethod(t *testing.T) {
	pkg, err := packager.LoadFromJson("../packager", "iso87BPackager.json")
	if err != nil {
		t.Fatalf("Error al cargar el packager: %v", err)
	}
	msg := NewMessage(pkg)

	RegisterStructField[CustomerInfo](msg, 48)

	msg.SetField(0, "0200")
	msg.SetField(3, "000000")
	msg.RegisterField(11, &IntField{})
	msg.SetField(11, 123456)
	msg.SetField(48, CustomerInfo{ID: "C123", Name: "Log Test"})
	msg.SetField(62, "aabbcc")

	logOutput := msg.Log()
	t.Logf("Log Output:\n%s", logOutput)

	// Verificar que el JSON contiene los campos esperados
	if !strings.Contains(logOutput, "\"0\": \"0200\"") {
		t.Error("Log no contiene MTI")
	}
	if !strings.Contains(logOutput, "\"3\": \"000000\"") {
		t.Error("Log no contiene Processing Code")
	}
	if !strings.Contains(logOutput, "\"11\": 123456") {
		t.Error("Log no contiene STAN")
	}
	if !strings.Contains(logOutput, "\"48\": {\n    \"id\": \"C123\",\n    \"name\": \"Log Test\"\n  }") {
		t.Error("Log no contiene CustomerInfo serializado correctamente")
	}
	if !strings.Contains(logOutput, "\"62\": \"aabbcc\"") {
		t.Error("Log no contiene BytesField serializado correctamente")
	}
}
