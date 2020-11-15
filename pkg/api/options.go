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

// AddOptions are store options
type AddOptions struct {
	// MergeLinks requests to merge links with an existing object links
	MergeLinks bool
	// MultiLink allows multiple links b/w the same objects
	// NOTE: this option is ignored if MergeLinks is false
	MultiLink bool
}

// AddOption sets AddOptions
type AddOption func(*AddOptions)

// NewOptions returns default add options
func NewAddOptions() AddOptions {
	return AddOptions{
		MergeLinks: false,
		MultiLink:  false,
	}
}

// LinkOptions are link options
type LinkOptions struct {
	// UID is optional link UID
	UID uuid.UID
	// Multi allows multiple links b/w the same objects
	Multi bool
	// Metadata
	Metadata metadata.Metadata
}

// LinkOption sets LinkOptions
type LinkOption func(*LinkOptions)

// NewLinkOptions returns default link options
func NewLinkOptions() LinkOptions {
	return LinkOptions{
		Multi:    false,
		Metadata: metadata.New(),
	}
}
