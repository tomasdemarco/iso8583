package field

import (
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/padding"
	"github.com/tomasdemarco/iso8583/prefix"
)

type Field struct {
	Description string            `json:"description"`
	Type        Type              `json:"type"`
	Length      int               `json:"length"`
	Pattern     string            `json:"pattern"`
	Encoding    encoding.Encoding `json:"encoding"`
	Prefix      prefix.Prefix     `json:"prefix"`
	Padding     padding.Padding   `json:"padding"`
	//SubfieldsData subfield.SubfieldsData `json:"subFieldsData"`
}
