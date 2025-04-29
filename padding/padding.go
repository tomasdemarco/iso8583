package padding

import "github.com/tomasdemarco/iso8583/utils"

type Padding struct {
	Type     Type                 `json:"type"`
	Position Position             `json:"position"`
	Char     utils.ByteFromString `json:"char"`
}
