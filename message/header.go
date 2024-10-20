package message

import (
	"errors"
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/packager"
	"sort"
	"strconv"
)

func UnpackHeader(messageRaw string, packager *packager.Packager) (header map[string]string, err error) {
	header = make(map[string]string)
	position := packager.PrefixLength

	for i, v := range packager.Header.HeaderFields {
		if len(messageRaw) > position+v.Length {
			header[i] = messageRaw[position : position+v.Length]
			position += v.Length
		} else {
			return nil, errors.New("cannot get header, message length is short")
		}
	}

	return header, err
}

func (m *Message) PackHeader(packager *packager.Packager) (header string) {
	for i := range packager.Header.HeaderFields {
		if packager.Header.HeaderFields[i].DefaultValue != "" {
			m.Header[i] = packager.Header.HeaderFields[i].DefaultValue
		}
	}
	keys := make([]string, 0, len(m.Header))
	for k := range m.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, v := range keys {
		if packager.Header.HeaderFields[v].InRequest {
			header += m.Header[v]
		}
	}
	return header
}

func (m *Message) PackHeaderResponse(packager *packager.Packager) (header string) {
	for i := 1; i < len(m.Header); i++ {
		if packager.Header.HeaderFields[fmt.Sprintf("%02d", i)].InvertPrevious {
			aux := m.Header[fmt.Sprintf("%02d", i)]
			m.Header[fmt.Sprintf("%02d", i)] = m.Header[fmt.Sprintf("%02d", i-1)]
			m.Header[fmt.Sprintf("%02d", i-1)] = aux
		}
	}
	for i, v := range m.Header {
		if packager.Header.HeaderFields[i].InRequest {
			header += v
		}
	}
	return header
}

func UnpackLength(messageRaw []byte) (length int, err error) {
	length64, err := strconv.ParseInt(fmt.Sprintf("%x", messageRaw[:2]), 16, 64)
	length = int(length64)
	return length, err
}

func PackLength(messageRaw string, headerLength int) (lengthHex string) {
	lengthHex = encoding.HexEncode(fmt.Sprintf("%d", (len(messageRaw)+headerLength)/2))
	return fmt.Sprintf("%04s", lengthHex)
}
