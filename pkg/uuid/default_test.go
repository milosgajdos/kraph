package uuid

import "testing"

func TestNewFromString(t *testing.T) {
	u := "randomUID"

	uid := NewFromString(u)

	if u != uid.String() {
		t.Errorf("expected: %s, got: %s", u, uid)
	}
}

func TestNew(t *testing.T) {
	u1 := New()
	u2 := New()

	if len(u1.String()) == 0 || len(u2.String()) == 0 {
		t.Errorf("empty uid returned")
		return
	}

	if u1.String() == u2.String() {
		t.Errorf("non-unique uids genericerated")
	}
}
