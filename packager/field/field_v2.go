// Package field defines the structure and behavior of ISO 8583 fields.
package field

// BasePackager defines the common interface for packagers.
type BasePackager interface {
	// GetLength returns the length of the packager.
	GetLength()
}

// Packager defines the interface for packing and unpacking data of a specific type.
type Packager[T any] interface {
	BasePackager
	// Pack packs a value of type T into a byte slice.
	// It returns the packed data, a plain string representation, and an error if packing fails.
	Pack(T) ([]byte, string, error)
	// Unpack unpacks a byte slice into a value of type T.
	// It returns the unpacked value, the number of bytes consumed, and an error if unpacking fails.
	Unpack([]byte, int) (T, int, error)
}
