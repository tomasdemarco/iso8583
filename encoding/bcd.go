package encoding

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/utils"
)

func BcdDecode(src []byte) (string, error) {
	return fmt.Sprintf("%x", src), nil
}

func BcdEncode(src string) []byte {
	return utils.Hex2Byte(src)
}
