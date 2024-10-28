package message

import (
	"github.com/tomasdemarco/iso8583/packager"
	"testing"
)

var (
	PaddingTypes      = []string{"PARITY", "FIXED"}
	PaddingPositions  = []string{"RIGHT", "LEFT"}
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

			message := Message{}
			fieldsPackager := packager.FieldsPackager{}
			fieldsPackager.Length = 6
			fieldsPackager.Padding.Type = paddingType
			fieldsPackager.Padding.Position = paddingPosition
			fieldsPackager.Padding.Pad = "0"

			fields := make(map[string]packager.FieldsPackager)
			fields["011"] = fieldsPackager

			pkg := packager.Packager{}
			pkg.Fields = fields
			message.Packager = &pkg

			resultR, resultL := message.UnpackPadding("011")

			if resultR != expectedResultR {
				t.Fatalf(`UnpackPadding(%s) PaddingType=%s - PaddingPosition=%s - PaddingPad=%s - ResultR "%d" does not match "%d"`, data, paddingType, paddingPosition, "0", resultR, expectedResultR)
			}
			t.Logf(`UnpackPadding=%s PaddingType=%-6s - PaddingPosition=%s - PaddingPad=%s - ResultR "%d" match "%d"`, data, paddingType, paddingPosition, "0", resultR, expectedResultR)

			if resultL != expectedResultL {
				t.Fatalf(`UnpackPadding(%s) PaddingType=%s - PaddingPosition=%s - PaddingPad=%s - ResultL "%d" does not match "%d"`, data, paddingType, paddingPosition, "0", resultL, expectedResultL)
			}
			t.Logf(`UnpackPadding=%s PaddingType=%-6s - PaddingPosition=%s - PaddingPad=%s - ResultL "%d" match "%d"`, data, paddingType, paddingPosition, "0", resultL, expectedResultL)
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

			message := Message{}
			fieldsPackager := packager.FieldsPackager{}
			fieldsPackager.Length = 6
			fieldsPackager.Padding.Type = paddingType
			fieldsPackager.Padding.Position = paddingPosition
			fieldsPackager.Padding.Pad = "0"

			fields := make(map[string]packager.FieldsPackager)
			fields["011"] = fieldsPackager

			pkg := packager.Packager{}
			pkg.Fields = fields
			message.Packager = &pkg

			message.SetField("011", data)

			resultR, resultL := message.PackPadding("011")

			if resultR != expectedResultR {
				t.Fatalf(`PackPadding(%s) PaddingType=%s - PaddingPosition=%s - PaddingPad=%s - ResultR "%s" does not match "%s"`, data, paddingType, paddingPosition, "0", resultR, expectedResultR)
			}
			t.Logf(`PackPadding=%s PaddingType=%-6s - PaddingPosition=%s - PaddingPad=%s - ResultR "%s" match "%s"`, data, paddingType, paddingPosition, "0", resultR, expectedResultR)

			if resultL != expectedResultL {
				t.Fatalf(`PackPadding(%s) PaddingType=%s - PaddingPosition=%s - PaddingPad=%s - ResultL "%s" does not match "%s"`, data, paddingType, paddingPosition, "0", resultL, expectedResultL)
			}
			t.Logf(`PackPadding=%s PaddingType=%-6s - PaddingPosition=%s - PaddingPad=%s - ResultL "%s" match "%s"`, data, paddingType, paddingPosition, "0", resultL, expectedResultL)
		}
	}
}
