// Package prefix provides functionalities for handling length prefixes in ISO 8583 messages.
package prefix

// Prefixer defines the interface for encoding and decoding length prefixes.
type Prefixer interface {
	// EncodeLength encodes an integer length into a byte slice according to the prefixer's rules.
	// It returns the encoded length as a byte slice and an error if encoding fails.
	EncodeLength(length int) ([]byte, error)
	// DecodeLength decodes a length from a byte slice starting at a given offset.
	// It returns the decoded integer length, the number of bytes consumed, and an error if decoding fails.
	DecodeLength(b []byte, offset int) (int, error)
	// GetPackedLength returns the fixed number of bytes this prefixer uses to encode the length.
	GetPackedLength() int
}

// Prefixers is a collection of common Prefixer implementations for various length types.
type Prefixers struct {
	FIXED  Prefixer
	L      Prefixer
	LL     Prefixer
	LLL    Prefixer
	LLLL   Prefixer
	LLLLL  Prefixer
	LLLLLL Prefixer
}

// BinaryPrefixers is a collection of common Prefixer implementations for binary length types.
type BinaryPrefixers struct {
	FIXED  Prefixer
	B      Prefixer
	BB     Prefixer
	BBB    Prefixer
	BBBB   Prefixer
	BBBBB  Prefixer
	BBBBBB Prefixer
}
