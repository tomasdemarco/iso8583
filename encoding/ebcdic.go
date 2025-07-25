// Package encoding provides various data encoding and decoding functionalities for ISO 8583 fields.
package encoding

import "fmt"

var (
	// asciiToEbcdic is a lookup table for converting ASCII characters to EBCDIC.
	asciiToEbcdic = []byte{
		'\x00', '\x01', '\x02', '\x03', '\x37', '\x2D', '\x2E', '\x2F',
		'\x16', '\x05', '\x25', '\x0B', '\x0C', '\x0D', '\x0E', '\x0F',
		'\x10', '\x11', '\x12', '\x13', '\x3C', '\x3D', '\x32', '\x26',
		'\x18', '\x19', '\x3F', '\x27', '\x1C', '\x1D', '\x1E', '\x1F',
		'\x40', '\x4F', '\x7F', '\x7B', '\x5B', '\x6C', '\x50', '\x7D',
		'\x4D', '\x5D', '\x5C', '\x4E', '\x6B', '\x60', '\x4B', '\x61',
		'\xF0', '\xF1', '\xF2', '\xF3', '\xF4', '\xF5', '\xF6', '\xF7',
		'\xF8', '\xF9', '\x7A', '\x5E', '\x4C', '\x7E', '\x6E', '\x6F',
		'\x7C', '\xC1', '\xC2', '\xC3', '\xC4', '\xC5', '\xC6', '\xC7',
		'\xC8', '\xC9', '\xD1', '\xD2', '\xD3', '\xD4', '\xD5', '\xD6',
		'\xD7', '\xD8', '\xD9', '\xE2', '\xE3', '\xE4', '\xE5', '\xE6',
		'\xE7', '\xE8', '\xE9', '\x4A', '\xE0', '\x5A', '\x5F', '\x6D',
		'\x79', '\x81', '\x82', '\x83', '\x84', '\x85', '\x86', '\x87',
		'\x88', '\x89', '\x91', '\x92', '\x93', '\x94', '\x95', '\x96',
		'\x97', '\x98', '\x99', '\xA2', '\xA3', '\xA4', '\xA5', '\xA6',
		'\xA7', '\xA8', '\xA9', '\xC0', '\x6A', '\xD0', '\xA1', '\x07',
		'\x20', '\x21', '\x22', '\x23', '\x24', '\x15', '\x06', '\x17',
		'\x28', '\x29', '\x2A', '\x2B', '\x2C', '\x09', '\x0A', '\x1B',
		'\x30', '\x31', '\x1A', '\x33', '\x34', '\x35', '\x36', '\x08',
		'\x38', '\x39', '\x3A', '\x3B', '\x04', '\x14', '\x3E', '\xE1',
		'\x41', '\x42', '\x43', '\x44', '\x45', '\x46', '\x47', '\x48',
		'\x49', '\x51', '\x52', '\x53', '\x54', '\x55', '\x56', '\x57',
		'\x58', '\x59', '\x62', '\x63', '\x64', '\x65', '\x66', '\x67',
		'\x68', '\x69', '\x70', '\x71', '\x72', '\x73', '\x74', '\x75',
		'\x76', '\x77', '\x78', '\x80', '\x8A', '\x8B', '\x8C', '\x8D',
		'\x8E', '\x8F', '\x90', '\x9A', '\x9B', '\x9C', '\x9D', '\x9E',
		'\x9F', '\xA0', '\xAA', '\xAB', '\xAC', '\xAD', '\xAE', '\xAF',
		'\xB0', '\xB1', '\xB2', '\xB3', '\xB4', '\xB5', '\xB6', '\xB7',
		'\xB8', '\xB9', '\xBA', '\xBB', '\xBC', '\xBD', '\xBE', '\xBF',
		'\xCA', '\xCB', '\xCC', '\xCD', '\xCE', '\xCF', '\xDA', '\xDB',
		'\xDC', '\xDD', '\xDE', '\xDF', '\xEA', '\xEB', '\xEC', '\xED',
		'\xEE', '\xEF', '\xFA', '\xFB', '\xFC', '\xFD', '\xFE', '\xFF'}

	// ebcdicToAscii is a lookup table for converting EBCDIC characters to ASCII.
	ebcdicToAscii = []byte{
		'\x00', '\x01', '\x02', '\x03', '\x9C', '\x09', '\x86', '\x7F',
		'\x97', '\x8D', '\x8E', '\x0B', '\x0C', '\x0D', '\x0E', '\x0F',
		'\x10', '\x11', '\x12', '\x13', '\x9D', '\x85', '\x08', '\x87',
		'\x18', '\x19', '\x92', '\x8F', '\x1C', '\x1D', '\x1E', '\x1F',
		'\x80', '\x81', '\x82', '\x83', '\x84', '\x0A', '\x17', '\x1B',
		'\x88', '\x89', '\x8A', '\x8B', '\x8C', '\x05', '\x06', '\x07',
		'\x90', '\x91', '\x16', '\x93', '\x94', '\x95', '\x96', '\x04',
		'\x98', '\x99', '\x9A', '\x9B', '\x14', '\x15', '\x9E', '\x1A',
		'\x20', '\xA0', '\xA1', '\xA2', '\xA3', '\xA4', '\xA5', '\xA6',
		'\xA7', '\xA8', '\x5B', '\x2E', '\x3C', '\x28', '\x2B', '\x21',
		'\x26', '\xA9', '\xAA', '\xAB', '\xAC', '\xAD', '\xAE', '\xAF',
		'\xB0', '\xB1', '\x5D', '\x24', '\x2A', '\x29', '\x3B', '\x5E',
		'\x2D', '\x2F', '\xB2', '\xB3', '\xB4', '\xB5', '\xB6', '\xB7',
		'\xB8', '\xB9', '\x7C', '\x2C', '\x25', '\x5F', '\x3E', '\x3F',
		'\xBA', '\xBB', '\xBC', '\xBD', '\xBE', '\xBF', '\xC0', '\xC1',
		'\xC2', '\x60', '\x3A', '\x23', '\x40', '\x27', '\x3D', '\x22',
		'\xC3', '\x61', '\x62', '\x63', '\x64', '\x65', '\x66', '\x67',
		'\x68', '\x69', '\xC4', '\xC5', '\xC6', '\xC7', '\xC8', '\xC9',
		'\xCA', '\x6A', '\x6B', '\x6C', '\x6D', '\x6E', '\x6F', '\x70',
		'\x71', '\x72', '\xCB', '\xCC', '\xCD', '\xCE', '\xCF', '\xD0',
		'\xD1', '\x7E', '\x73', '\x74', '\x75', '\x76', '\x77', '\x78',
		'\x79', '\x7A', '\xD2', '\xD3', '\xD4', '\xD5', '\xD6', '\xD7',
		'\xD8', '\xD9', '\xDA', '\xDB', '\xDC', '\xDD', '\xDE', '\xDF',
		'\xE0', '\xE1', '\xE2', '\xE3', '\xE4', '\xE5', '\xE6', '\xE7',
		'\x7B', '\x41', '\x42', '\x43', '\x44', '\x45', '\x46', '\x47',
		'\x48', '\x49', '\xE8', '\xE9', '\xEA', '\xEB', '\xEC', '\xED',
		'\x7D', '\x4A', '\x4B', '\x4C', '\x4D', '\x4E', '\x4F', '\x50',
		'\x51', '\x52', '\xEE', '\xEF', '\xF0', '\xF1', '\xF2', '\xF3',
		'\x5C', '\x9F', '\x53', '\x54', '\x55', '\x56', '\x57', '\x58',
		'\x59', '\x5A', '\xF4', '\xF5', '\xF6', '\xF7', '\xF8', '\xF9',
		'\x30', '\x31', '\x32', '\x33', '\x34', '\x35', '\x36', '\x37',
		'\x38', '\x39', '\xFA', '\xFB', '\xFC', '\xFD', '\xFE', '\xFF'}
)

// EBCDIC implements the Encoder interface for EBCDIC encoding.
// It converts ASCII strings to EBCDIC byte slices and vice-versa using lookup tables.
type EBCDIC struct {
	length int
}

// NewEbcdicEncoder creates a new EBCDIC encoder.
func NewEbcdicEncoder() Encoder {
	return &EBCDIC{}
}

// Encode converts an ASCII string to an EBCDIC byte slice.
// It uses the asciiToEbcdic lookup table.
func (e *EBCDIC) Encode(src string) ([]byte, error) {
	var dst []byte
	for _, v := range []byte(src) {
		dst = append(dst, asciiToEbcdic[v])
	}
	return dst, nil
}

// Decode converts an EBCDIC byte slice to an ASCII string.
// It reads up to the configured length and uses the ebcdicToAscii lookup table.
func (e *EBCDIC) Decode(src []byte) (string, error) {
	if len(src) < e.length {
		return "", fmt.Errorf("%w: expected %d, got %d", ErrNotEnoughDataToDecode, e.length, len(src))
	}

	var dst []byte
	for _, v := range src[:e.length] {
		dst = append(dst, ebcdicToAscii[v])
	}

	return string(dst), nil
}

// SetLength sets the length for the EBCDIC encoder.
func (e *EBCDIC) SetLength(length int) {
	e.length = length
}
