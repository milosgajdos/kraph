package graph

import (
	"strings"
	"testing"

	"github.com/milosgajdos/kraph/pkg/api"
	"github.com/milosgajdos/kraph/pkg/api/generic"
	"github.com/milosgajdos/kraph/pkg/uuid"
)

const (
	resName    = "resName"
	resGroup   = "resGroup"
	resVersion = "resVersion"
	resKind    = "resKind"
	objUID     = "testID"
	objName    = "testName"
	objNs      = "testNs"
)

func TestDOTID(t *testing.T) {
	o := generic.NewObject(uuid.NewFromString(objUID), objName, objNs, nil, api.Options{})

	if _, err := DOTID(o); err == nil {
		t.Errorf("expected error, got: %v", err)
	}

	r := generic.NewResource(resName, resGroup, resVersion, resKind, true, api.Options{})
	o = generic.NewObject(uuid.NewFromString(objUID), objName, objNs, r, api.Options{})

	exp := strings.Join([]string{
		resGroup,
		resVersion,
		resKind,
		objNs,
		objName}, "/")

	dotid, err := DOTID(o)
	if err != nil {
		t.Fatalf("failed to generate DOTID: %v", err)
	}

	if dotid != exp {
		t.Errorf("expected DOTID: %s, got: %s", exp, dotid)
	}
}
