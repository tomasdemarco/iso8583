// Package packager defines the structure for ISO 8583 packagers,
// which describe the format of messages and their fields.
package packager

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/packager/field"
	"github.com/tomasdemarco/iso8583/padding"
	"github.com/tomasdemarco/iso8583/prefix"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// PackagerDto represents the structure of a packager as defined in a JSON file.
// It is used to deserialize the packager configuration.
type PackagerDto struct {
	Description string              `json:"description"`
	Prefix      prefix.Prefix       `json:"prefix"`
	Fields      map[string]FieldDto `json:"fields"`
}

// FieldDto represents the structure of a field as defined in a JSON file.
// It is used to deserialize the field configuration.
type FieldDto struct {
	Description string            `json:"description"`
	Type        field.Type        `json:"type"`
	Length      int               `json:"length"`
	Pattern     string            `json:"pattern"`
	Encoding    encoding.Encoding `json:"encoding"` // Encoding type (ASCII, BCD, etc.)
	Prefix      prefix.Prefix     `json:"prefix"`   // Length prefix configuration
	Padding     padding.Padding   `json:"padding"`  // Padding configuration
}

// LoadFromJson loads the packager configuration from a JSON file.
// It takes the directory path and the JSON file name.
// It returns a Packager instance and an error if loading fails.
func LoadFromJson(path, file string) (*Packager, error) {
	absPath, err := filepath.Abs(filepath.Join(path, file))
	if err != nil {
		return nil, err
	}

	jsonFile, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var pkgDto PackagerDto
	err = json.Unmarshal(byteValue, &pkgDto)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	var pkg Packager

	pkg.Description = pkgDto.Description
	pf, err := GetPrefixer(pkgDto.Prefix)
	if err != nil {
		return nil, err
	}

	pkg.Prefix = pf

	fields := make(map[int]field.Packager)

	for k, v := range pkgDto.Fields {
		fld, err := SetField(v)
		if err != nil {
			return nil, err
		}

		kNum, err := strconv.Atoi(k)
		if err != nil {
			return nil, fmt.Errorf("invalid field number: %s", k)
		}

		fields[kNum] = fld
	}

	pkg.Fields = fields

	return &pkg, err
}

// SetField converts a FieldDto (from JSON) to a field.Field instance.
// It initializes the encoding, prefix, padding, and regex pattern components.
// It returns a field.Field instance and an error if initialization fails.
func SetField(f FieldDto) (field.Packager, error) {
	length := f.Length
	if f.Encoding == encoding.Bcd {
		length = length / 2
	}

	enc, err := GetEncoder(f.Encoding, bcdPadLeft(f.Padding))
	if err != nil {
		return nil, err
	}

	//if f.Type == field.Bitmap {
	//	return field.NewBitmapField(f.Description, length, enc, true), nil
	//}

	pf, err := GetPrefixer(f.Prefix)
	if err != nil {
		return nil, err
	}

	pad, err := GetPadder(f.Padding)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(f.Pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern for field %w", err)
	}

	return field.NewField(f.Description, f.Type, length, re, enc, pf, pad), nil
}

// GetPrefixer creates a Prefixer instance based on the prefix.Prefix configuration.
// It returns the Prefixer interface and an error if the prefix type is invalid.
func GetPrefixer(pf prefix.Prefix) (prefix.Prefixer, error) {
	switch pf.Encoding {
	case encoding.Bcd:
		return prefix.NewBcdPrefixer(pf.Type.EnumIndex(), pf.Hex, pf.IsInclusive), nil
	case encoding.Ebcdic:
		return prefix.NewEbcdicPrefixer(pf.Type.EnumIndex(), pf.Hex, pf.IsInclusive), nil
	case encoding.Binary:
		return prefix.NewBinaryPrefixer(pf.Type.EnumIndex(), pf.IsInclusive), nil
	case encoding.Ascii:
		return prefix.NewAsciiPrefixer(pf.Type.EnumIndex(), pf.Hex, pf.IsInclusive), nil
	default:
		return prefix.NONE.FIXED, nil
	}
}

// GetEncoder creates an Encoder instance based on the encoding.Encoding type.
// It takes an encoding.Encoding value and a boolean for BCD left padding.
// It returns the Encoder interface and an error if the encoding type is invalid.
func GetEncoder(enc encoding.Encoding, bcdPadLeft bool) (encoding.Encoder, error) {
	switch enc {
	case encoding.Bcd:
		return encoding.NewBcdEncoder(bcdPadLeft), nil
	case encoding.Ebcdic:
		return &encoding.EBCDIC{}, nil
	case encoding.Binary:
		return &encoding.BINARY{}, nil
	case encoding.Ascii:
		return &encoding.ASCII{}, nil
	default:
		return nil, errors.New("invalid encoding")
	}
}

// GetPadder creates a Padder instance based on the padding.Padding configuration.
// It returns the Padder interface and an error if the padding type is invalid.
func GetPadder(p padding.Padding) (padding.Padder, error) {
	switch p.Type {
	case padding.Parity:
		switch p.Position {
		case padding.Left:
			return padding.NewParityPadder(true, p.Char), nil
		case padding.Right:
			return padding.NewParityPadder(false, p.Char), nil
		default:
			return nil, errors.New("invalid padding")
		}
	case padding.Fill:
		switch p.Position {
		case padding.Left:
			return padding.NewFillPadder(true, p.Char), nil
		case padding.Right:
			return padding.NewFillPadder(false, p.Char), nil
		default:
			return nil, errors.New("invalid padding")
		}
	default:
		return padding.NONE.NONE, nil
	}
}

// bcdPadLeft determines if left padding should be applied for BCD encoding.
// This is relevant for BCD fields with left parity padding.
func bcdPadLeft(p padding.Padding) bool {
	if p.Type == padding.Parity && p.Position == padding.Left {
		return true
	}
	return false
}
