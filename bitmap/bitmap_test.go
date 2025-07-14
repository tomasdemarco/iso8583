package bitmap

import (
	"github.com/tomasdemarco/iso8583/packager"
	"reflect"
	"testing"
)

func TestUnpack(t *testing.T) {
	pkg, err := packager.LoadFromJson("../message", "test_packager.json")
	if err != nil {
		t.Fatalf("Error al cargar el packager de prueba: %v", err)
	}

	testCases := []struct {
		name           string
		packedBytes    []byte
		expectedFields []int
		expectedLength int
		expectError    bool
	}{
		{
			name:           "Primary Bitmap - Fields 2, 3, 5",
			packedBytes:    []byte{0x68, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expectedFields: []int{2, 3, 5},
			expectedLength: 8,
			expectError:    false,
		},
		{
			name:           "Primary and Secondary Bitmap - Fields 1, 2, 3, 5, 65",
			packedBytes:    []byte{0xE8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expectedFields: []int{1, 2, 3, 5, 65},
			expectedLength: 16,
			expectError:    false,
		},
		{
			name:           "Primary and Secondary Bitmap - Unpack Secondary Error",
			packedBytes:    []byte{0xE8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expectedFields: nil,
			expectedLength: 0,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bmap, _, err := Unpack(pkg.Fields[1], tc.packedBytes, 0)
			if (err != nil) != tc.expectError {
				t.Fatalf("Unpack() error = %v, wantErr %v", err, tc.expectError)
			}
			if !tc.expectError {
				if !reflect.DeepEqual(bmap.GetSliceString(), tc.expectedFields) {
					t.Errorf("Unpack() fields = %v, want %v", bmap.GetSliceString(), tc.expectedFields)
				}
				if len(bmap.ToBytes()) != tc.expectedLength {
					t.Errorf("Unpack() length = %d, want %d", len(bmap.ToBytes()), tc.expectedLength)
				}
			}
		})
	}
}
