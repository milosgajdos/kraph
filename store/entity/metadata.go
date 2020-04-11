package entity

// Metadata is a simple key-value store for arbitrary data
type Metadata map[string]interface{}

// Get reads the value for the given key and returns it
func (m Metadata) Get(key string) interface{} {
	return m[key]
}

// Set sets the value for the given key
func (m *Metadata) Set(key string, val interface{}) {
	(*m)[key] = val
}
