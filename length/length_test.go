package length_test

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/tomasdemarco/iso8583/length"
	"github.com/tomasdemarco/iso8583/prefix"
)

// MockPrefixer implements the prefix.Prefixer interface for testing purposes.
type MockPrefixer struct {
	hex,
	isInclusive bool
	encodeFunc    func(lenMessage int) ([]byte, error)
	decodeFunc    func(b []byte, offset int) (int, error)
	getPackedFunc func() int
}

func (m *MockPrefixer) EncodeLength(lenMessage int) ([]byte, error) {
	if m.encodeFunc != nil {
		return m.encodeFunc(lenMessage)
	}
	return []byte{byte(lenMessage)}, nil // Default simple encoding
}

func (m *MockPrefixer) DecodeLength(b []byte, offset int) (int, error) {
	if m.decodeFunc != nil {
		return m.decodeFunc(b, offset)
	}
	if len(b) > offset {
		return int(b[offset]), nil // Default simple decoding
	}
	return 0, errors.New("not enough bytes to decode")
}

func (m *MockPrefixer) GetPackedLength() int {
	if m.getPackedFunc != nil {
		return m.getPackedFunc()
	}
	return 1 // Default packed length
}

func (p *MockPrefixer) IsHex() {
	p.hex = true
}

func (p *MockPrefixer) IsInclusive() {
	p.isInclusive = true
}

func TestPack(t *testing.T) {
	tests := []struct {
		name        string
		prefixer    prefix.Prefixer
		lenMessage  int
		expected    []byte
		expectedErr error
	}{
		{
			name: "Successful Parse",
			prefixer: &MockPrefixer{
				encodeFunc: func(lenMessage int) ([]byte, error) {
					return []byte{byte(lenMessage)}, nil
				},
			},
			lenMessage:  10,
			expected:    []byte{10},
			expectedErr: nil,
		},
		{
			name: "Parse Error from Prefixer",
			prefixer: &MockPrefixer{
				encodeFunc: func(lenMessage int) ([]byte, error) {
					return nil, errors.New("encoding error")
				},
			},
			lenMessage:  20,
			expected:    nil,
			expectedErr: errors.New("encoding error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := length.Pack(tt.prefixer, tt.lenMessage)
			if !bytes.Equal(got, tt.expected) {
				t.Errorf("Parse() got = %v, want %v", got, tt.expected)
			}
			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) || (err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestUnpack(t *testing.T) {
	tests := []struct {
		name        string
		reader      *bufio.Reader
		prefixer    prefix.Prefixer
		expectedLen int
		expectedErr error
	}{
		{
			name:        "Successful Unparse",
			reader:      bufio.NewReader(bytes.NewBuffer([]byte{10})),
			prefixer:    &MockPrefixer{getPackedFunc: func() int { return 1 }},
			expectedLen: 10,
			expectedErr: nil,
		},
		{
			name:        "Unparse EOF Error",
			reader:      bufio.NewReader(bytes.NewBuffer([]byte{})),
			prefixer:    &MockPrefixer{getPackedFunc: func() int { return 1 }},
			expectedLen: 0,
			expectedErr: io.EOF,
		},
		{
			name:   "Unparse Decode Error from Prefixer",
			reader: bufio.NewReader(bytes.NewBuffer([]byte{99})),
			prefixer: &MockPrefixer{
				getPackedFunc: func() int { return 1 },
				decodeFunc: func(b []byte, offset int) (int, error) {
					return 0, errors.New("decode error")
				},
			},
			expectedLen: 0,
			expectedErr: errors.New("decode error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLen, err := length.Unpack(tt.reader, tt.prefixer)
			if gotLen != tt.expectedLen {
				t.Errorf("Unparse() gotLen = %v, want %v", gotLen, tt.expectedLen)
			}
			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) || (err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("Unparse() error = %v, wantErr %v", err, tt.expectedErr)
			}
		})
	}
}
