package store

// attrs are graph attributes
type attrs map[string]string

// NewAttributes creates new attributes and returns it
func NewAttributes() *attrs {
	attrs := make(attrs)

	return &attrs
}

// Keys returns all the metadata keys
func (a attrs) Keys() []string {
	keys := make([]string, len(a))

	i := 0
	for key := range a {
		keys[i] = key
		i++
	}

	return keys
}

// Get reads an attribute value for the given key and returns it.
// It returns an empty string if the attribute was not found.
func (a attrs) Get(key string) string {
	return a[key]
}

// Set sets an attribute to the given value
func (a *attrs) Set(key, val string) {
	(*a)[key] = val
}
