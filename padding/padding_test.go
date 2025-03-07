package padding

import (
	"github.com/tomasdemarco/iso8583/utils"
	"testing"
)

var (
	PaddingTypes      = []Type{Parity, Fill}
	PaddingPositions  = []Position{Right, Left}
	ValuesPaddingR    = []string{"0", "", "00000", ""}
	ValuesPaddingL    = []string{"", "0", "", "00000"}
	ValuesNumPaddingR = []int{1, 0, 0, 0}
	ValuesNumPaddingL = []int{0, 1, 0, 0}
)

// TestUnpackPadding calls message.UnpackPadding
func TestUnpackPadding(t *testing.T) {
	for pt, paddingType := range PaddingTypes {
		for pp, paddingPosition := range PaddingPositions {
			data := "1"
			expectedResultR := ValuesNumPaddingR[pp+(pt*2)]
			expectedResultL := ValuesNumPaddingL[pp+(pt*2)]

			padding := Padding{}
			padding.Type = paddingType
			padding.Position = paddingPosition
			padding.Char = utils.ByteFromString('0')

			resultL, resultR := Unpack(padding, len(data))

			if resultL != expectedResultL {
				t.Fatalf(`UnpackPadding(%s) PaddingType=%s - PaddingPosition=%s - PaddingChar=%x - ResultL "%d" does not match "%d"`, data, paddingType.String(), paddingPosition.String(), padding.Char, resultL, expectedResultL)
			}
			t.Logf(`UnpackPadding(%s) PaddingType=%s - PaddingPosition=%s - PaddingChar=%x - ResultL "%d" match "%d"`, data, paddingType.String(), paddingPosition.String(), padding.Char, resultL, expectedResultL)

			if resultR != expectedResultR {
				t.Fatalf(`UnpackPadding(%s) PaddingType=%s - PaddingPosition=%s - PaddingChar=%x - ResultR "%d" does not match "%d"`, data, paddingType.String(), paddingPosition.String(), padding.Char, resultR, expectedResultR)
			}
			t.Logf(`UnpackPadding(%s) PaddingType=%s - PaddingPosition=%s - PaddingChar=%x - ResultR "%d" match "%d"`, data, paddingType.String(), paddingPosition.String(), padding.Char, resultR, expectedResultR)
		}
	}
}

// TestPackPadding calls message.PackPadding
func TestPackPadding(t *testing.T) {
	for pt, paddingType := range PaddingTypes {
		for pp, paddingPosition := range PaddingPositions {
			data := "1"
			expectedResultR := ValuesPaddingR[pp+(pt*2)]
			expectedResultL := ValuesPaddingL[pp+(pt*2)]

			padding := Padding{}
			padding.Type = paddingType
			padding.Position = paddingPosition
			padding.Char = utils.ByteFromString('0')

			resultL, resultR := Pack(padding, 6, len(data))

			if resultL != expectedResultL {
				t.Fatalf(`PackPadding(%s) PaddingType=%s - PaddingPosition=%s - PaddingChar=%x - ResultL "%s" does not match "%s"`, data, paddingType.String(), paddingPosition.String(), padding.Char, resultL, expectedResultL)
			}
			t.Logf(`PackPadding(%s) PaddingType=%s - PaddingPosition=%s - PaddingChar=%x - ResultL "%s" match "%s"`, data, paddingType.String(), paddingPosition.String(), padding.Char, resultL, expectedResultL)

			if resultR != expectedResultR {
				t.Fatalf(`PackPadding(%s) PaddingType=%s - PaddingPosition=%s - PaddingChar=%x - ResultR "%s" does not match "%s"`, data, paddingType.String(), paddingPosition.String(), padding.Char, resultR, expectedResultR)
			}
			t.Logf(`PackPadding(%s) PaddingType=%s - PaddingPosition=%s - PaddingChar=%x - ResultR "%s" match "%s"`, data, paddingType.String(), paddingPosition.String(), padding.Char, resultR, expectedResultR)
		}
	}
}
