package prefix

type Option func(Prefixer)

func IsHex() Option {
	return func(s Prefixer) {
		s.IsHex()
	}
}

func IsInclusive() Option {
	return func(s Prefixer) {
		s.IsInclusive()
	}
}
