package packager

import (
	"encoding/json"
	"fmt"
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

func LoadPackagers() (Packager, Packager, Packager, Packager) {
	absPath, _ := filepath.Abs("./iso8583/packager/iso87BVisaBase1Packager.json")
	jsonFile, err := os.Open(absPath)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var packagerVisa Packager
	json.Unmarshal(byteValue, &packagerVisa)
	defer jsonFile.Close()

	if packagerVisa.HeaderFile != "" {
		packagerVisa.Header = LoadHeaders(packagerVisa.HeaderFile)
	}

	for i, v := range packagerVisa.Fields {
		if v.SubFieldsFile != "" && v.SubFieldsFormat != "" {
			subFields := LoadSubfield(v.SubFieldsFile)
			fields := packagerVisa.Fields[i]
			fields.SubFields = subFields
			packagerVisa.Fields[i] = fields
		}
	}

	absPath, _ = filepath.Abs("./iso8583/packager/iso87EMasterPackager.json")
	jsonFile, err = os.Open(absPath)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ = ioutil.ReadAll(jsonFile)
	var packagerMaster Packager
	json.Unmarshal(byteValue, &packagerMaster)
	defer jsonFile.Close()

	if packagerMaster.HeaderFile != "" {
		packagerMaster.Header = LoadHeaders(packagerMaster.HeaderFile)
	}

	for i, v := range packagerMaster.Fields {
		if v.SubFieldsFile != "" && v.SubFieldsFormat != "" {
			subFields := LoadSubfield(v.SubFieldsFile)
			fields := packagerMaster.Fields[i]
			fields.SubFields = subFields
			packagerMaster.Fields[i] = fields
		}
	}

	absPath, _ = filepath.Abs("./iso8583/packager/iso87BPackager.json")
	jsonFile, err = os.Open(absPath)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ = ioutil.ReadAll(jsonFile)
	var packagerGp Packager
	json.Unmarshal(byteValue, &packagerGp)
	defer jsonFile.Close()

	if packagerGp.HeaderFile != "" {
		packagerGp.Header = LoadHeaders(packagerGp.HeaderFile)
	}

	for i, v := range packagerGp.Fields {
		if v.SubFieldsFile != "" && v.SubFieldsFormat != "" {
			subFields := LoadSubfield(v.SubFieldsFile)
			field := packagerGp.Fields[i]
			field.SubFields = subFields
			packagerGp.Fields[i] = field
		}
	}

	absPath, _ = filepath.Abs("./iso8583/packager/iso87BCabalPackager.json")
	jsonFile, err = os.Open(absPath)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ = ioutil.ReadAll(jsonFile)
	var packagerCabal Packager
	json.Unmarshal(byteValue, &packagerCabal)
	defer jsonFile.Close()

	if packagerCabal.HeaderFile != "" {
		packagerCabal.Header = LoadHeaders(packagerCabal.HeaderFile)
	}

	for i, v := range packagerCabal.Fields {
		if v.SubFieldsFile != "" && v.SubFieldsFormat != "" {
			subFields := LoadSubfield(v.SubFieldsFile)
			field := packagerCabal.Fields[i]
			field.SubFields = subFields
			packagerCabal.Fields[i] = field
		}
	}

	return packagerGp, packagerVisa, packagerMaster, packagerCabal
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
