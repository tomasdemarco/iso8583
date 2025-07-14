package emv

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func Unpack(value string, tags ...string) (map[string]string, error) {
	var tagsValidate = make(map[string]string)

	tagsFilter := map[string]string{}
	for _, v := range tags {
		tagsFilter[v] = v
	}

	position := 0
	de55Filtered := ""

	for position < len(value) {
		tagIdLength := 2
		if strings.ToUpper(value[position+1:position+2]) == "F" {
			tagIdLength = 4
		}

		tag := value[position : position+tagIdLength]
		position += tagIdLength

		length, err := strconv.ParseInt(value[position:position+2], 16, 10)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidTagLength, err)
		}

		if _, ok := tagsFilter[tag]; ok || tags == nil {
			tagsValidate[tag] = value[position+2 : position+2+int(length*2)]
			de55Filtered += tag + value[position:position+2] + value[position+2:position+2+int(length*2)]
		}

		position += 2 + int(length*2)
	}

	for _, tag := range tags {
		if _, ok := tagsValidate[tag]; !ok {
			return nil, fmt.Errorf("%w: '%s'", ErrTagNotFound, tag)
		}
	}

	return tagsValidate, nil
}

func Pack(tags map[string]string) string {
	var result strings.Builder

	sortedTags := make([]string, 0, len(tags))
	for tag := range tags {
		sortedTags = append(sortedTags, tag)
	}

	sort.Strings(sortedTags)

	for _, tag := range sortedTags {
		result.WriteString(fmt.Sprintf("%s%02d%s", tag, len(tags[tag]), tags[tag]))
	}

	return result.String()
}
