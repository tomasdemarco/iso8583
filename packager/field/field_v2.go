package field

type BasePackager interface {
	GetLength()
}

type Packager[T any] interface {
	BasePackager
	Pack(T) ([]byte, string, error)
	Unpack([]byte, int) (T, int, error)
}
