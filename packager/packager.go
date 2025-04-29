package packager

import (
	"encoding/json"
	"github.com/tomasdemarco/iso8583/packager/field"
	"github.com/tomasdemarco/iso8583/prefix"
	"github.com/tomasdemarco/iso8583/suffix"
	"io"
	"os"
	"path/filepath"
)

type Packager struct {
	Description string                 `json:"description"`
	Prefix      prefix.Prefixer        `json:"prefix"`
	Suffix      suffix.Suffix          `json:"suffix"`
	Fields      map[string]field.Field `json:"fields"`
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

	//for i, v := range pkg.Fields {
	//	if v.SubfieldsData.File != "" {
	//		subFields, err := loadSubfields(path, v.SubfieldsData.File)
	//		if err != nil {
	//			return pkg, err
	//		}
	//
	//		fields := pkg.Fields[i]
	//		fields.SubfieldsData.SubFields = subFields
	//		pkg.Fields[i] = fields
	//	}
	//}

	return pkg, nil
}

//func loadSubfields(path, file string) (subFields map[string]Field, err error) {
//	absPath, err := filepath.Abs(path + "/" + file)
//	if err != nil {
//		return subFields, err
//	}
//
//	jsonFile, err := os.Open(absPath)
//	if err != nil {
//		return subFields, err
//	}
//
//	byteValue, err := io.ReadAll(jsonFile)
//	if err != nil {
//		return subFields, err
//	}
//
//	err = json.Unmarshal(byteValue, &subFields)
//	if err != nil {
//		return subFields, err
//	}
//
//	defer jsonFile.Close()
//
//	return subFields, nil
//}
