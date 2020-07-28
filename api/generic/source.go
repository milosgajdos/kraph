package generic

type source struct {
	src string
}

func (s *source) String() string {
	return s.src
}
