package padding

type Option func(Padder)

func WithChar(char string) Option {
	return func(s Padder) {
		s.SetChar(char)
	}
}
