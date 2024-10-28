package packager

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Packager struct {
	Name           string                    `json:"name"`
	PrefixLength   int                       `json:"prefixLength"`
	PrefixEncoding int                       `json:"prefixEncoding"`
	HeaderLength   int                       `json:"headerLength"`
	Header         Header                    `json:"header"`
	HeaderFile     string                    `json:"headerFile"`
	Fields         map[string]FieldsPackager `json:"fields"`
}

type FieldsPackager struct {
	Type            string              `json:"type"`
	Length          int                 `json:"length"`
	Pattern         string              `json:"pattern"`
	Name            string              `json:"name"`
	Encoding        string              `json:"encoding"`
	Prefix          string              `json:"prefix"`
	PrefixEncoding  string              `json:"prefixEncoding"`
	Padding         Padding             `json:"padding"`
	SubFieldsFile   string              `json:"subFieldsFile"`
	SubFieldsFormat string              `json:"subFieldsFormat"`
	SubFields       map[string]SubField `json:"subFields"`
}

type SubField struct {
	Type           string  `json:"type"`
	Length         int     `json:"length"`
	Pattern        string  `json:"pattern"`
	Name           string  `json:"name"`
	Encoding       string  `json:"encoding"`
	Prefix         string  `json:"prefix"`
	PrefixEncoding string  `json:"prefixEncoding"`
	Padding        Padding `json:"padding"`
}

type Padding struct {
	Type     string `json:"type"`
	Position string `json:"position"`
	Pad      string `json:"pad"`
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

func LoadPackagers(file string) Packager {
	absPath, _ := filepath.Abs("./iso8583/packager/" + file)
	jsonFile, err := os.Open(absPath)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := io.ReadAll(jsonFile)

	var pkg Packager
	json.Unmarshal(byteValue, &pkg)
	defer jsonFile.Close()

	if pkg.HeaderFile != "" {
		pkg.Header = LoadHeaders(pkg.HeaderFile)
	}

	for i, v := range pkg.Fields {
		if v.SubFieldsFile != "" && v.SubFieldsFormat != "" {
			subFields := LoadSubfield(v.SubFieldsFile)
			fields := pkg.Fields[i]
			fields.SubFields = subFields
			pkg.Fields[i] = fields
		}
	}

	return pkg
}

func LoadHeaders(headerFile string) Header {
	absPath, _ := filepath.Abs("./iso8583/packager/" + headerFile)
	jsonFile, err := os.Open(absPath)
	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var header Header
	json.Unmarshal(byteValue, &header)
	defer jsonFile.Close()

	return header
}

func LoadSubfield(subFieldsFile string) (subFields map[string]SubField) {
	absPath, _ := filepath.Abs("./iso8583/packager/" + subFieldsFile)
	jsonFile, err := os.Open(absPath)
	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &subFields)
	defer jsonFile.Close()

	return subFields
}
