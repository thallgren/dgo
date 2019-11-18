package tf

import (
	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/internal"
)

// AllOf returns a type that represents all values that matches all of the included types
func AllOf(types ...interface{}) dgo.Value {
	return internal.AllOfType(types)
}

// AnyOf returns a type that represents all values that matches at least one of the included types
func AnyOf(types ...interface{}) dgo.Value {
	return internal.AnyOfType(types)
}

// OneOf returns a type that represents all values that matches exactly one of the included types
func OneOf(types ...interface{}) dgo.Value {
	return internal.OneOfType(types)
}

// Not returns a type that represents all values that are not represented by the given type
func Not(t dgo.Value) dgo.Value {
	return internal.NotType(t)
}
