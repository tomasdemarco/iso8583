package message

import (
	"errors"
	"strconv"
	"strings"
)

var tags = [16]string{"82", "84", "95", "9A", "9C", "5F2A", "9F02", "9F03", "9F10", "9F1A", "9F26", "9F27", "9F33", "9F34", "9F36", "9F37"}

func (m *Message) UnpackAndFilterTagsEmv() error {
	de55, err := m.GetField("055")
	if err == nil {
		var tagsValidate = make(map[string]string)
		tagsFilter := map[string]string{}
		for _, v := range tags {
			tagsFilter[v] = v
		}
		position := 0
		de55Filtered := ""
		for position < len(de55) {
			if strings.ToUpper(de55[position+1:position+2]) == "F" {
				tag := de55[position : position+4]
				position += 4
				length, err := strconv.ParseInt(de55[position:position+2], 16, 10)
				if err != nil {
					return err
				}
				if _, ok := tagsFilter[tag]; ok {
					tagsValidate[tag] = de55[position+2 : position+2+int(length*2)]
					de55Filtered += tag + de55[position:position+2] + de55[position+2:position+2+int(length*2)]
				}
				position += 2 + int(length*2)
			} else {
				tag := de55[position : position+2]
				position += 2
				length, err := strconv.ParseInt(de55[position:position+2], 16, 10)
				if err != nil {
					return err
				}
				if _, ok := tagsFilter[tag]; ok {
					de55Filtered += tag + de55[position:position+2] + de55[position+2:position+2+int(length*2)]
					tagsValidate[tag] = de55[position+2 : position+2+int(length*2)]
				}
				position += 2 + int(length*2)
			}
		}
		m.TagsEmv = tagsValidate
		for _, tag := range tags {
			if _, ok := m.TagsEmv[tag]; !ok {
				return errors.New("field 55 does not contain the tag '" + tag + "'")
			}
		}
		m.SetField("055", de55Filtered)
	} else {
		return errors.New("field 55 must be present in emv transactions")
	}
	return nil
}
