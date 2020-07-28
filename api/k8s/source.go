package k8s

type source struct {
	src string
}

// String returns source as a string
func (s *source) String() string {
	return s.src
}
