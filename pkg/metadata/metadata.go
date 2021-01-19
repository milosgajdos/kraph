package metadata

// Metadata provides a simple key-valule store
// for arbitrary entity data of arbitrary type.
type Metadata interface {
	// Keys returns all metadata keys.
	Keys() []string
	// Get returns the attribute value for the given key.
	Get(string) interface{}
	// Set sets the value of the attribute for the given key.
	Set(string, interface{})
}
