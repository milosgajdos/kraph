package graph

import (
	"strings"

	"github.com/milosgajdos/kraph/pkg/api"
)

// DOTID returns GraphViz DOT ID for the given api.Object.
// NOTE: the returned DOTID follows this naming convention:
// resourceGroup/resourceVersion/resourceKind/objectNamespace/objectName
func DOTID(obj api.Object) (string, error) {
	if obj.Resource() == nil {
		return "", ErrMissingResource
	}

	return strings.Join([]string{
		obj.Resource().Group(),
		obj.Resource().Version(),
		obj.Resource().Kind(),
		obj.Namespace(),
		obj.Name()}, "/"), nil
}
