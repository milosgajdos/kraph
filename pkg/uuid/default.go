package uuid

import (
	"github.com/google/uuid"
)

// uid implements UID
type uid struct {
	id string
}

// NewFromString returns new UID from uid string
func NewFromString(u string) *uid {
	return &uid{
		id: u,
	}
}

// New returns new UID
func New() *uid {
	return &uid{
		id: uuid.New().String(),
	}
}

// String returns API Object UID as string
func (u *uid) String() string {
	return u.id
}
