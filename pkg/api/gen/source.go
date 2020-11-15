package gen

// Source is the source of the API.
// It's recommended to initialize it
// at the very list to the URL or the
// name of the API for brevity.
type Source struct {
	src string
}

// NewSource returns API source.
func NewSource(s string) *Source {
	return &Source{
		src: s,
	}
}

// String implements Stringer interface.
func (s *Source) String() string {
	return s.src
}
