// Package encoding provides various data encoding and decoding functionalities for ISO 8583 fields.
package encoding

// Encoder defines the interface for encoding and decoding data.
type Encoder interface {
	// Encode converts a string into a byte slice according to the encoder's rules.
	// It returns the encoded data as a byte slice and an error if encoding fails.
	Encode(src string) ([]byte, error)
	// Decode converts a byte slice into a string according to the encoder's rules.
	// It returns the decoded string and an error if decoding fails.
	Decode(src []byte) (string, error)
	// SetLength sets the length for the encoder, which might be used during encoding or decoding.
	SetLength(length int)
}
