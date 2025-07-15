// Package packager defines the structure for ISO 8583 packagers,
// which describe the format of messages and their fields.
package packager

import (
	"github.com/tomasdemarco/iso8583/packager/field"
	"github.com/tomasdemarco/iso8583/prefix"
)

// Packager represents the definition of an ISO 8583 message format.
// It contains a description, the message prefix, and a map of fields.
type Packager struct {
	Description string                 `json:"description"`
	Prefix      prefix.Prefixer        `json:"prefix"`
	Fields      map[int]field.Packager `json:"fields"`
}
