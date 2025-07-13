package bitmap

import (
	"github.com/tomasdemarco/iso8583/packager"
	"reflect"
	"testing"

	"sort"
)

func TestPack(t *testing.T) {
	testCases := []struct {
		name          string
		fieldNumbers  []int
		expectedBytes []byte
		expectError   bool
	}{
		{
			name:          "Primary Bitmap - Fields 2, 3, 5",
			fieldNumbers:  []int{2, 3, 5},
			expectedBytes: []byte{0x68, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expectError:   false,
		},
		{
			name:          "Primary + Secondary Bitmap - Fields 1, 65, 128",
			fieldNumbers:  []int{1, 65, 128},
			expectedBytes: []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualBytes, err := Pack(tc.fieldNumbers)
			if (err != nil) != tc.expectError {
				t.Fatalf("Pack() error = %v, wantErr %v", err, tc.expectError)
			}
			if !tc.expectError && !reflect.DeepEqual(actualBytes.ToBytes(), tc.expectedBytes) {
				t.Errorf("Pack() = %x, want %x", actualBytes.ToBytes(), tc.expectedBytes)
			}
		})
	}
}

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

func TestBitmapRoundTrip(t *testing.T) {
	pkg, err := packager.LoadFromJson("../message", "test_packager.json")
	if err != nil {
		t.Fatalf("Error al cargar el packager de prueba: %v", err)
	}

	testCases := []struct {
		name          string
		fieldNumbers  []int
		initialLength int // Initial length to pass to Unpack
	}{
		{
			name:          "Primary Bitmap Round Trip",
			fieldNumbers:  []int{2, 3, 5},
			initialLength: 8,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Pack the field numbers
			bmap, err := Pack(tc.fieldNumbers)
			if err != nil {
				t.Fatalf("Pack() failed: %v", err)
			}

			// Unpack the bytes
			bmap2, _, err := Unpack(pkg.Fields[1], bmap.ToBytes(), 0)
			if err != nil {
				t.Fatalf("Unpack() failed: %v", err)
			}

			// Sort both slices for consistent comparison
			sort.Ints(tc.fieldNumbers)

			if !reflect.DeepEqual(bmap2.GetSliceString(), tc.fieldNumbers) {
				t.Errorf("Round trip failed.\nExpected: %v\nGot:      %v", tc.fieldNumbers, bmap2.GetSliceString())
			}
		})
	}
}
