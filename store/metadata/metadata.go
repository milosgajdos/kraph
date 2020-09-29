package metadata

// Metadata is a simple key-value store for arbitrary data
type Metadata map[string]interface{}

// NewMetadata creates new metadata and returns it
func New() *Metadata {
	md := make(Metadata)

	return &md
}

// Get reads the value for the given key and returns it
func (m Metadata) Get(key string) interface{} {
	return m[key]
}

// Set sets the value for the given key
func (m *Metadata) Set(key string, val interface{}) {
	(*m)[key] = val
}

// Keys returns all the metadata keys
func (m Metadata) Keys() []string {
	keys := make([]string, len(m))

	i := 0
	for key := range m {
		keys[i] = key
	}

	return keys
}
