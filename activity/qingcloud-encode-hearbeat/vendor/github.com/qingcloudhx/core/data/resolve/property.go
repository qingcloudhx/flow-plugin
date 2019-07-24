package resolve

import (
	"fmt"
	"github.com/qingcloudhx/core/data/property"

	"github.com/qingcloudhx/core/data"
)

var propertyResolverInfo = NewResolverInfo(true, true)

type PropertyResolver struct {
}

func (*PropertyResolver) GetResolverInfo() *ResolverInfo {
	return propertyResolverInfo
}

//PropertyResolver Property Resolver $property[item]
func (*PropertyResolver) Resolve(scope data.Scope, item string, field string) (interface{}, error) {

	manager := property.DefaultManager()
	value, exists := manager.GetProperty(item) //should we add the path and reset it to ""
	if !exists {
		return nil, fmt.Errorf("failed to resolve Property: '%s', ensure that property is configured in the application", item)
	}

	return value, nil
}
