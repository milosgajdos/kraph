package gen

type Source struct {
	src string
}

func NewSource(s string) *Source {
	return &Source{
		src: s,
	}
}

func (s *Source) String() string {
	return s.src
}
