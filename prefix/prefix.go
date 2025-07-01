package prefix

import (
	enc "github.com/tomasdemarco/iso8583/encoding"
)

type Prefix struct {
	Type        Type         `json:"type"`
	Encoding    enc.Encoding `json:"encoding"`
	Hex         bool         `json:"hex"`
	IsInclusive bool         `json:"isInclusive"`
}
