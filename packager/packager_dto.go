package packager

import (
	"encoding/json"
	"errors"
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/packager/field"
	"github.com/tomasdemarco/iso8583/padding"
	"github.com/tomasdemarco/iso8583/prefix"
	"github.com/tomasdemarco/iso8583/suffix"
	"io"
	"os"
	"path/filepath"
)

type PackagerDto struct {
	Description string              `json:"description"`
	Prefix      prefix.Prefix       `json:"prefix"`
	Suffix      suffix.Suffix       `json:"suffix"`
	Fields      map[string]FieldDto `json:"fields"`
}

type FieldDto struct {
	Description string            `json:"description"`
	Type        field.Type        `json:"type"`
	Length      int               `json:"length"`
	Pattern     string            `json:"pattern"`
	Encoding    encoding.Encoding `json:"encoding"`
	Prefix      prefix.Prefix     `json:"prefix"`
	Padding     padding.Padding   `json:"padding"`
	//SubfieldsData subfield.SubfieldsData `json:"subFieldsData"`
}

func LoadFromJsonV2(path, file string) (*Packager, error) {
	absPath, err := filepath.Abs(path + "/" + file)
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

	fields := make(map[string]field.Field)

	for k, v := range pkgDto.Fields {
		enc, err := GetEncoder(v.Encoding)
		if err != nil {
			return nil, err
		}

		pf, err := GetPrefixer(v.Prefix)
		if err != nil {
			return nil, err
		}

		if v.Prefix.Hex {
			pf.SetHex()
		}

		length := v.Length
		if v.Encoding == encoding.Binary || v.Encoding == encoding.Bcd {
			length = length / 2
		}

		pad, err := GetPadder(v.Padding)
		if err != nil {
			return nil, err
		}

		pad.SetChar(v.Padding.Char)

		fields[k] = field.Field{
			Description: v.Description,
			Type:        v.Type,
			Length:      length,
			Pattern:     v.Pattern,
			Encoding:    enc,
			Prefix:      pf,
			Padding:     pad,
		}
	}

	pkg.Fields = fields

	return &pkg, err
}

func GetPrefixer(pf prefix.Prefix) (prefix.Prefixer, error) {
	switch pf.Encoding {
	case encoding.Bcd:
		switch pf.Type {
		case prefix.L:
			return prefix.BCD.L, nil
		case prefix.LL:
			return prefix.BCD.LL, nil
		case prefix.LLL:
			return prefix.BCD.LLL, nil
		case prefix.LLLL:
			return prefix.BCD.LLLL, nil
		case prefix.LLLLL:
			return prefix.BCD.LLLLL, nil
		case prefix.LLLLLL:
			return prefix.BCD.LLLLLL, nil
		default:
			return nil, errors.New("invalid prefix")
		}
	case encoding.Ebcdic:
		switch pf.Type {
		case prefix.L:
			return prefix.EBCDIC.L, nil
		case prefix.LL:
			return prefix.EBCDIC.LL, nil
		case prefix.LLL:
			return prefix.EBCDIC.LLL, nil
		case prefix.LLLL:
			return prefix.EBCDIC.LLLL, nil
		case prefix.LLLLL:
			return prefix.EBCDIC.LLLLL, nil
		case prefix.LLLLLL:
			return prefix.EBCDIC.LLLLLL, nil
		default:
			return nil, errors.New("invalid prefix")
		}
	case encoding.Binary:
		switch pf.Type {
		case prefix.LL:
			return prefix.BINARY.B, nil
		case prefix.LLLL:
			return prefix.BINARY.BB, nil
		default:
			return nil, errors.New("invalid prefix")
		}
	case encoding.Ascii:
		switch pf.Type {
		case prefix.L:
			return prefix.ASCII.L, nil
		case prefix.LL:
			return prefix.ASCII.LL, nil
		case prefix.LLL:
			return prefix.ASCII.LLL, nil
		case prefix.LLLL:
			return prefix.ASCII.LLLL, nil
		case prefix.LLLLL:
			return prefix.ASCII.LLLLL, nil
		case prefix.LLLLLL:
			return prefix.ASCII.LLLLLL, nil
		default:
			return nil, errors.New("invalid prefix")
		}
	default:
		return prefix.NONE.Fixed, nil
	}
}

func GetEncoder(enc encoding.Encoding) (encoding.Encoder, error) {
	switch enc {
	case encoding.Bcd:
		return &encoding.BCD{}, nil
	case encoding.Ebcdic:
		return &encoding.EBCDIC{}, nil
	case encoding.Binary:
		return &encoding.BINARY{}, nil
	case encoding.Ascii:
		return &encoding.ASCII{}, nil
	case encoding.Hex:
		return &encoding.HEX{}, nil
	default:
		return nil, errors.New("invalid encoding")
	}
}

func GetPadder(p padding.Padding) (padding.Padder, error) {
	switch p.Type {
	case padding.Parity:
		switch p.Position {
		case padding.Left:
			return padding.PARITY.LEFT, nil
		case padding.Right:
			return padding.PARITY.RIGHT, nil
		default:
			return nil, errors.New("invalid padding")
		}
	case padding.Fill:
		switch p.Position {
		case padding.Left:
			return padding.FILL.LEFT, nil
		case padding.Right:
			return padding.FILL.RIGHT, nil
		default:
			return nil, errors.New("invalid padding")
		}
	default:
		return padding.NONE.NONE, nil
	}
}
