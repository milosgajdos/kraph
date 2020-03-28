package k8s

import (
	"testing"

	"github.com/milosgajdos/kraph/query"
)

func TestResources(t *testing.T) {
	api := MockAPI()

	resources := api.Resources()
	if len(resources) != MockAPIResCount {
		t.Errorf("expected resource count: %d, got: %d", len(resources), MockAPIResCount)
	}

	group := "odd"
	expCount := 4
	oddResources := api.Resources(query.Group(group))
	if len(oddResources) != expCount {
		t.Errorf("expected %d odd resources, got: %d", expCount, len(oddResources))
	}

	group = "even"
	expCount = MockAPIResCount - 4
	evenResources := api.Resources(query.Group(group))
	if len(evenResources) != expCount {
		t.Errorf("expected %d even resources, got: %d", expCount, len(evenResources))
	}

	expCount = 1
	version := "v2"
	v2Res := api.Resources(query.Version(version))
	if len(v2Res) != expCount {
		t.Errorf("expected %d version %s resources, got: %d", expCount, version, len(v2Res))
	}

	expCount = 0
	group = "odd"
	v2OddRes := api.Resources(query.Group(group), query.Version(version))
	if len(v2OddRes) != expCount {
		t.Errorf("expected %d resources version: %s, group: %s, got: %d", expCount, version, group, len(v2OddRes))
	}

	name := "evenRes"
	group = "even"
	v2EvenRes := api.Resources(query.Name(name), query.Group(group), query.Version(version))
	if len(v2EvenRes) != 1 {
		t.Errorf("expected to find resource: %s version: %s group: %s", name, version, group)
	}
}
