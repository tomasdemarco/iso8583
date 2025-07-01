package packager

import (
	"encoding/json"
	"errors"
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/packager/field"
	"github.com/tomasdemarco/iso8583/padding"
	"github.com/tomasdemarco/iso8583/prefix"
	"io"
	"os"
	"path/filepath"
)

type PackagerDto struct {
	Description string              `json:"description"`
	Prefix      prefix.Prefix       `json:"prefix"`
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

	pf.SetHex(pkgDto.Prefix.Hex)
	pf.SetIsInclusive(pkgDto.Prefix.IsInclusive)
	pkg.Prefix = pf

	fields := make(map[string]field.Field)

	for k, v := range pkgDto.Fields {
		fld, err := SetField(v)
		if err != nil {
			return nil, err
		}
		fields[k] = *fld
	}

	pkg.Fields = fields

	return &pkg, err
}

func SetField(f FieldDto) (*field.Field, error) {
	length := f.Length
	if f.Encoding == encoding.Binary || f.Encoding == encoding.Bcd {
		length = length / 2
	}

	enc, err := GetEncoder(f.Encoding)
	if err != nil {
		return nil, err
	}

	pf, err := GetPrefixer(f.Prefix)
	if err != nil {
		return nil, err
	}

	pf.SetHex(f.Prefix.Hex)
	pf.SetIsInclusive(f.Prefix.IsInclusive)

	pad, err := GetPadder(f.Padding)
	if err != nil {
		return nil, err
	}

	pad.SetChar(f.Padding.Char)

	return &field.Field{
		Description: f.Description,
		Type:        f.Type,
		Length:      length,
		Pattern:     f.Pattern,
		Encoding:    enc,
		Prefix:      pf,
		Padding:     pad,
	}, nil
}

func GetPrefixer(pf prefix.Prefix) (prefix.Prefixer, error) {
	switch pf.Encoding {
	case encoding.Bcd:
		return prefix.NewBcdPrefixer(pf.Type.EnumIndex(), pf.Hex, pf.IsInclusive), nil
	case encoding.Ebcdic:
		return prefix.NewEbcdicPrefixer(pf.Type.EnumIndex(), pf.Hex, pf.IsInclusive), nil
	case encoding.Binary:
		return prefix.NewBinaryPrefixer(pf.Type.EnumIndex(), pf.Hex, pf.IsInclusive), nil
	case encoding.Ascii:
		return prefix.NewAsciiPrefixer(pf.Type.EnumIndex(), pf.Hex, pf.IsInclusive), nil
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
