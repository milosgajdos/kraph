package metadata

// metadata is a simple key-value store for arbitrary data
type metadata map[string]interface{}

// NewMetadata creates new metadata and returns it
func New() *metadata {
	md := make(metadata)

	return &md
}

// Get reads the value for the given key and returns it
func (m metadata) Get(key string) interface{} {
	return m[key]
}

// Set sets the value for the given key
func (m *metadata) Set(key string, val interface{}) {
	(*m)[key] = val
}

// Keys returns all the metadata keys
func (m metadata) Keys() []string {
	keys := make([]string, len(m))

	i := 0
	for key := range m {
		keys[i] = key
	}

	return keys
}
