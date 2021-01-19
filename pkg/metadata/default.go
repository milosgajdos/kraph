package metadata

// metadata is a simple key-value store for arbitrary data.
type metadata map[string]interface{}

// NewMetadata creates new metadata and returns it.
func New() *metadata {
	md := make(metadata)

	return &md
}

// NewCopyFrom copies metadata from m and returns it.
func NewCopyFrom(m Metadata) *metadata {
	md := make(metadata)

	if m != nil {
		for _, k := range m.Keys() {
			md[k] = m.Get(k)
		}
	}

	return &md
}

// NewFromMap creates new metadata from m and returns it.
func NewFromMap(m map[string]interface{}) *metadata {
	md := make(metadata)

	if m != nil {
		for k, v := range m {
			md[k] = v
		}
	}

	return &md
}

// Get reads the value for the given key and returns it.
func (m metadata) Get(key string) interface{} {
	return m[key]
}

// Set sets the value for the given key.
func (m *metadata) Set(key string, val interface{}) {
	(*m)[key] = val
}

// Keys returns all the metadata keys.
func (m metadata) Keys() []string {
	keys := make([]string, len(m))

	i := 0
	for key := range m {
		keys[i] = key
	}

	return keys
}
