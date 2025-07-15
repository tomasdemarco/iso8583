package subfield

import "github.com/tomasdemarco/iso8583/packager/field"

type Packager interface {
	field.Field
	Pack(map[string]string) (string, error)
	Unpack(string) (map[string]string, error)
}
