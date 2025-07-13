package utils

import (
	"reflect"
	"sort"
	"testing"
)

func TestBitSet(t *testing.T) {
	// Test NewBitSet and GetSize
	t.Run("NewBitSet and GetSize", func(t *testing.T) {
		bs := NewBitSet(8, 8) // 8 bits, max 8 bits
		if bs.GetSize() != 8 {
			t.Errorf("Expected size 8 bits, got %d bits", bs.GetSize())
		}
		if len(bs.ToBytes()) != 1 { // 8 bits = 1 byte
			t.Errorf("Expected bytes slice length 1, got %d", len(bs.ToBytes()))
		}

		bs = NewBitSet(16, 16) // 16 bits, max 16 bits
		if bs.GetSize() != 16 {
			t.Errorf("Expected size 16 bits, got %d bits", bs.GetSize())
		}
		if len(bs.ToBytes()) != 2 { // 16 bits = 2 bytes
			t.Errorf("Expected bytes slice length 2, got %d", len(bs.ToBytes()))
		}

		bs = NewBitSet(64, 64) // 64 bits, max 64 bits
		if bs.GetSize() != 64 {
			t.Errorf("Expected size 64 bits, got %d bits", bs.GetSize())
		}
		if len(bs.ToBytes()) != 8 { // 64 bits = 8 bytes
			t.Errorf("Expected bytes slice length 8, got %d", len(bs.ToBytes()))
		}

		bs = NewBitSet(128, 128) // 128 bits, max 128 bits
		if bs.GetSize() != 128 {
			t.Errorf("Expected size 128 bits, got %d bits", bs.GetSize())
		}
		if len(bs.ToBytes()) != 16 { // 128 bits = 16 bytes
			t.Errorf("Expected bytes slice length 16, got %d", len(bs.ToBytes()))
		}
	})

	// Test Set and Get
	t.Run("Set and Get", func(t *testing.T) {
		bs := NewBitSet(16, 128) // 16 bits initial, max 128 bits

		// Set and get first bit (ISO field 1)
		bs.Set(1)
		if !bs.Get(1) {
			t.Errorf("Bit 1 should be set")
		}
		if bs.Get(2) {
			t.Errorf("Bit 2 should not be set")
		}

		// Set and get a middle bit
		bs.Set(8) // Last bit of first byte
		if !bs.Get(8) {
			t.Errorf("Bit 8 should be set")
		}

		bs.Set(9) // First bit of second byte
		if !bs.Get(9) {
			t.Errorf("Bit 9 should be set")
		}

		// Set and get last bit
		bs.Set(16) // Last bit of second byte
		if !bs.Get(16) {
			t.Errorf("Bit 16 should be set")
		}

		// Test out of bounds with auto-expansion
		initialSize := bs.GetSize()
		bs.Set(initialSize + 1) // Set the bit at the current size + 1 (which will trigger expansion)
		if !bs.Get(initialSize + 1) {
			t.Errorf("Bit at initial size %d should be set after auto-expansion", initialSize+1)
		}
		if bs.GetSize() <= initialSize {
			t.Errorf("BitSet size should have increased after setting bit %d", initialSize+1)
		}
		// Test setting a bit far beyond current size
		bs.Set(128) // Set bit 128 (should expand to 128 bits)
		if !bs.Get(128) {
			t.Errorf("Bit 128 should be set after auto-expansion")
		}
		if bs.GetSize() < 128 {
			t.Errorf("BitSet size should be at least 128 bits, got %d", bs.GetSize())
		}

		// Test setting a bit beyond maxSize
		bs = NewBitSet(8, 16) // 8 bits initial, max 16 bits
		bs.Set(17)            // Try to set bit 17, which is > maxSize (16)
		if bs.Get(17) {
			t.Errorf("Bit 17 should NOT be set as it exceeds maxSize")
		}
		if bs.GetSize() > 16 {
			t.Errorf("BitSet size should NOT exceed maxSize (16), got %d", bs.GetSize())
		}
	})

	// Test ToBytes
	t.Run("ToBytes", func(t *testing.T) {
		bs := NewBitSet(8, 8)    // 8 bits
		bs.Set(1)                // Bit 1 (MSB of first byte)
		bs.Set(8)                // Bit 8 (LSB of first byte)
		expected := []byte{0x81} // 10000001
		actual := bs.ToBytes()
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("ToBytes() for 0x81: Expected %x, got %x", expected, actual)
		}

		bs = NewBitSet(16, 16) // 16 bits
		bs.Set(1)              // Bit 1 (MSB of first byte)
		bs.Set(9)              // Bit 9 (MSB of second byte)
		expected = []byte{0x80, 0x80}
		actual = bs.ToBytes()
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("ToBytes() for 0x8080: Expected %x, got %x", expected, actual)
		}

		bs = NewBitSet(64, 64)                                            // 64 bits
		bs.Set(1)                                                         // Bit 1 (ISO field 1)
		bs.Set(2)                                                         // Bit 2 (ISO field 2)
		bs.Set(64)                                                        // Bit 64 (ISO field 64)
		expected = []byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01} // 1100...0001
		actual = bs.ToBytes()
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("ToBytes() for primary bitmap: Expected %x, got %x", expected, actual)
		}

		bs = NewBitSet(128, 128) // 128 bits
		bs.Set(1)                // Bit 1 (ISO field 1 - secondary bitmap indicator)
		bs.Set(65)               // Bit 65 (ISO field 65)
		bs.Set(128)              // Bit 128 (ISO field 128)
		expected = []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
		actual = bs.ToBytes()
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("ToBytes() for secondary bitmap: Expected %x, got %x", expected, actual)
		}
	})

	// Test Round-Trip (Set -> ToBytes -> NewBitSet -> Get)
	t.Run("Round-Trip", func(t *testing.T) {
		bitsToSet := []int{1, 8, 9, 16, 64, 65, 128}
		originalBs := NewBitSet(128, 128) // Use 128 bits for round-trip
		for _, bit := range bitsToSet {
			originalBs.Set(bit)
		}

		packedBytes := originalBs.ToBytes()

		unpackedBs := NewBitSet(len(packedBytes)*8, len(packedBytes)*8) // Create new BitSet from packed bytes, convert bytes to bits
		copy(unpackedBs.bytes, packedBytes)

		for i := 1; i <= originalBs.GetSize(); i++ { // Iterate up to original bit size
			if originalBs.Get(i) != unpackedBs.Get(i) {
				t.Errorf("Round-trip failed at bit %d. Original: %v, Unpacked: %v", i, originalBs.Get(i), unpackedBs.Get(i))
			}
		}
	})

	// Test Concatenate
	t.Run("Concatenate", func(t *testing.T) {
		bs1 := NewBitSet(8, 128) // 8 bits, max 128
		bs1.Set(1)
		bs1.Set(5)

		bs2 := NewBitSet(8, 128) // 8 bits, max 128
		bs2.Set(2)
		bs2.Set(7)

		// Concatenate bs1 and bs2
		concatenatedBs := bs1.Concatenate(bs2)

		// Expected bits: 1, 5 (from bs1) and 8+2=10, 8+7=15 (from bs2)
		expectedBits := []int{1, 5, 10, 15}
		sort.Ints(expectedBits)

		actualBits := concatenatedBs.GetSliceString()
		sort.Ints(actualBits)

		if !reflect.DeepEqual(actualBits, expectedBits) {
			t.Errorf("Concatenate failed.\nExpected: %v\nGot:      %v", expectedBits, actualBits)
		}

		// Test size of concatenated BitSet
		if concatenatedBs.GetSize() != (bs1.GetSize() + bs2.GetSize()) {
			t.Errorf("Concatenated BitSet size incorrect. Expected %d, got %d", (bs1.GetSize() + bs2.GetSize()), concatenatedBs.GetSize())
		}
		// Test maxSize of concatenated BitSet
		if concatenatedBs.maxSize != (bs1.maxSize + bs2.maxSize) {
			t.Errorf("Concatenated BitSet maxSize incorrect. Expected %d, got %d", (bs1.maxSize + bs2.maxSize), concatenatedBs.maxSize)
		}
	})
}
