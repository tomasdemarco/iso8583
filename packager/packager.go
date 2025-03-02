package packager

import (
	"encoding/json"
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/padding"
	"github.com/tomasdemarco/iso8583/prefix"
	"io"
	"os"
	"path/filepath"
)

type Packager struct {
	Name         string           `json:"name"`
	Prefix       prefix.Prefix    `json:"prefix"`
	HeaderLength int              `json:"headerLength"`
	Header       Header           `json:"header"`
	HeaderFile   string           `json:"headerFile"`
	Fields       map[string]Field `json:"fields"`
}

type Field struct {
	Name            string              `json:"name"`
	Type            string              `json:"type"`
	Length          int                 `json:"length"`
	Pattern         string              `json:"pattern"`
	Encoding        encoding.Encoding   `json:"encoding"`
	Prefix          prefix.Prefix       `json:"prefix"`
	Padding         padding.Padding     `json:"padding"`
	SubFieldsFile   string              `json:"subFieldsFile"`
	SubFieldsFormat string              `json:"subFieldsFormat"`
	SubFields       map[string]SubField `json:"subFields"`
}

type SubField struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Length   int               `json:"length"`
	Pattern  string            `json:"pattern"`
	Encoding encoding.Encoding `json:"encoding"`
	Prefix   prefix.Prefix     `json:"prefix"`
	Padding  padding.Padding   `json:"padding"`
}

type Header struct {
	Name         string                  `json:"name"`
	HeaderFields map[string]HeaderFields `json:"headerFields"`
}

type HeaderFields struct {
	Name           string `json:"name"`
	Length         int    `json:"length"`
	DefaultValue   string `json:"defaultValue"`
	InRequest      bool   `json:"inRequest"`
	InvertPrevious bool   `json:"invertPrevious"`
}

func LoadFromJson(path, file string) (pkg Packager, err error) {
	absPath, err := filepath.Abs(path + "/" + file)
	if err != nil {
		return pkg, err
	}

	jsonFile, err := os.Open(absPath)
	if err != nil {
		return pkg, err
	}

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return pkg, err
	}

	err = json.Unmarshal(byteValue, &pkg)
	if err != nil {
		return pkg, err
	}
	defer jsonFile.Close()

	if pkg.HeaderFile != "" {
		pkg.Header, err = loadHeader(path, pkg.HeaderFile)
		if err != nil {
			return pkg, err
		}
	}

	for i, v := range pkg.Fields {
		if v.SubFieldsFile != "" {
			subFields, err := loadSubfields(path, v.SubFieldsFile)
			if err != nil {
				return pkg, err
			}

			fields := pkg.Fields[i]
			fields.SubFields = subFields
			pkg.Fields[i] = fields
		}
	}

	return pkg, nil
}

func loadHeader(path, file string) (header Header, err error) {
	absPath, err := filepath.Abs(path + "/" + file)
	if err != nil {
		return header, err
	}

	jsonFile, err := os.Open(absPath)
	if err != nil {
		return header, err
	}

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return header, err
	}

	err = json.Unmarshal(byteValue, &header)
	if err != nil {
		return header, err
	}

	defer jsonFile.Close()

	return header, nil
}

func loadSubfields(path, file string) (subFields map[string]SubField, err error) {
	absPath, err := filepath.Abs(path + "/" + file)
	if err != nil {
		return subFields, err
	}

	jsonFile, err := os.Open(absPath)
	if err != nil {
		return subFields, err
	}

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return subFields, err
	}

	err = json.Unmarshal(byteValue, &subFields)
	if err != nil {
		return subFields, err
	}

	defer jsonFile.Close()

	return subFields, nil
}
