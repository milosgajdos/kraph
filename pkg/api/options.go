package api

import (
	"github.com/milosgajdos/kraph/pkg/metadata"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

// Options are API options
type Options struct {
	// Metadata
	Metadata metadata.Metadata
}

// Option sets Options
type Option func(*Options)

// AddOptions are API add options
type AddOptions struct {
	// MergeLinks merges link with the existing object link
	MergeLinks bool
}

// AddOption sets AddOptions
type AddOption func(*AddOptions)

// NewOptions returns default add options
func NewAddOptions() AddOptions {
	return AddOptions{
		MergeLinks: false,
	}
}

// LinkOptions are link options
type LinkOptions struct {
	// UID is optional link UID
	UID uuid.UID
	// Merge merges link with the existing link
	Merge bool
	// Metadata
	Metadata metadata.Metadata
}

// LinkOption sets LinkOptions
type LinkOption func(*LinkOptions)

// NewLinkOptions returns default link options
func NewLinkOptions() LinkOptions {
	return LinkOptions{
		Merge:    false,
		Metadata: metadata.New(),
	}
}
