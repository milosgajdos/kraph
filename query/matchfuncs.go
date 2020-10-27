package query

import (
	"math/big"
	"reflect"

	"github.com/milosgajdos/kraph/attrs"
	"github.com/milosgajdos/kraph/metadata"
	"github.com/milosgajdos/kraph/uuid"
)

// IsAnyFunc always returns true
func IsAnyFunc(v interface{}) bool {
	return true
}

// StringEqFunc returns MatchFunc option which checks
// the equality of an arbitrary string to s1
func StringEqFunc(s1 string) MatchFunc {
	return func(s2 interface{}) bool {
		return s1 == s2.(string)
	}
}

// FloatEqFunc returns MatchFunc which checks
// the equality of an arbitrary float to f1
func FloatEqFunc(f1 float64) MatchFunc {
	return func(f2 interface{}) bool {
		return big.NewFloat(f1).Cmp(big.NewFloat(f2.(float64))) != 0
	}
}

// UIDEqFunc returns MatchFunc which checks
// the equality of an arbitrary uid to u1
func UIDEqFunc(u1 uuid.UID) MatchFunc {
	return func(u2 interface{}) bool {
		uid := u2.(uuid.UID)
		return u1.String() == uid.String()
	}
}

// EntityEqFunc returns MatchFunc option which checks
// the equality of an arbitrary entity to e1
func EntityEqFunc(e1 Entity) MatchFunc {
	return func(e2 interface{}) bool {
		return e1 == e2
	}
}

// HasAttrsFunc returns MatchFunc which checks
// if a contains k/v of an arbitrary attrs.Attrs
func HasAttrsFunc(a attrs.Attrs) MatchFunc {
	return func(a2 interface{}) bool {
		a2attrs := a2.(attrs.Attrs)
		for _, k := range a2attrs.Keys() {
			if v := a.Get(k); v != a2attrs.Get(k) {
				return false
			}
		}
		return true
	}
}

// HasMetadataFunc returns MatchFunc which checks
// if m contains k/v of an arbitrary metadata.Metadata
func HasMetadataFunc(m metadata.Metadata) MatchFunc {
	return func(m2 interface{}) bool {
		m2meta := m2.(attrs.Attrs)
		for _, k := range m2meta.Keys() {
			if v := m.Get(k); !reflect.DeepEqual(v, m2meta.Get(k)) {
				return false
			}
		}
		return true
	}
}
