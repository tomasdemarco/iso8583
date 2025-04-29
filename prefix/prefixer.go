package prefix

type Prefixer interface {
	EncodeLength(length int) ([]byte, error)
	DecodeLength(b []byte, offset int) (int, error)
	GetPackedLength() int
	SetHex()
}

type Prefixers struct {
	Fixed  Prefixer
	L      Prefixer
	LL     Prefixer
	LLL    Prefixer
	LLLL   Prefixer
	LLLLL  Prefixer
	LLLLLL Prefixer
}

type BinaryPrefixers struct {
	Fixed  Prefixer
	B      Prefixer
	BB     Prefixer
	BBB    Prefixer
	BBBB   Prefixer
	BBBBB  Prefixer
	BBBBBB Prefixer
}
