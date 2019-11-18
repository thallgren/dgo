package dgo

import (
	"reflect"
	"regexp"
)

type (
	// IntegerType describes integers that are within an inclusive or exclusive range
	IntegerType interface {
		Value

		// Inclusive returns true if this range has an inclusive end
		Inclusive() bool

		// IsInstance returns true if the given int64 is an instance of this type
		IsInstance(int64) bool

		// Max returns the maximum constraint
		Max() int64

		// Min returns the minimum constraint
		Min() int64
	}

	// FloatRangeType describes floating point numbers that are within an inclusive or exclusive range
	FloatRangeType interface {
		Value

		// Inclusive returns true if this range has an inclusive end
		Inclusive() bool

		// IsInstance returns true if the given float64 is an instance of this type
		IsInstance(float64) bool

		// Max returns the maximum constraint
		Max() float64

		// Min returns the minimum constraint
		Min() float64
	}

	// BooleanType matches the true and false literals
	BooleanType interface {
		Value

		// IsInstance returns true if the Go native value is represented by this type
		IsInstance(value bool) bool
	}

	// SizedType is implemented by types that may have a size constraint
	// such as String, Array, or Map
	SizedType interface {
		// Max returns the maximum size for instances of this type
		Max() int

		// Min returns the minimum size for instances of this type
		Min() int

		// Unbounded returns true when the type has no size constraint
		Unbounded() bool
	}

	// StringType is a SizedType.
	StringType interface {
		Value
		SizedType

		// IsInstance returns true if the given string is an instance of this type
		IsInstance(s string) bool
	}

	PatternType interface {
		StringType

		GoRegexp() *regexp.Regexp
	}

	// NativeType is the type for all Native values
	NativeType interface {
		Value

		// GoType returns the reflect.Type
		GoType() reflect.Type
	}

	// AliasProvider replaces aliases with their concrete value.
	//
	// The parser uses this interface to perform in-place replacement of aliases
	AliasProvider interface {
		Replace(Value) Value
	}

	// AliasContainer is implemented by types that can contain other types.
	//
	// The parser uses this interface to perform in-place replacement of aliases
	AliasContainer interface {
		Resolve(AliasProvider)
	}

	// An AliasMap maps names to types and vice versa.
	AliasMap interface {
		// GetName returns the name for the given type or nil if the type isn't found
		GetName(t Value) String

		// GetType returns the type with the given name or nil if the type isn't found
		GetType(n String) Value

		// Add adds the type t with the given name to this map
		Add(t Value, name String)
	}

	// GenericType is implemented by types that represent themselves stripped from
	// range and size constraints.
	GenericType interface {
		// Generic returns the generic type that this type represents stripped
		// from range and size constraints
		Generic() Value
	}

	// Factory provides the New method that types use to create new instances
	Factory interface {
		// New creates instances of this type.
		New(Value) Value
	}

	// Named is implemented by named types such as the StructMap
	Named interface {
		Name() string
	}

	// DeepAssignable is implemented by values that need deep Assignable comparisons.
	DeepAssignable interface {
		DeepAssignable(guard RecursionGuard, other interface{}) bool
	}

	// ReverseAssignable indicates that the check for assignable must continue by delegating to the
	// type passed as an argument to the Assignable method. The reason is that types like AllOf, AnyOf
	// OneOf or types representing exact slices or maps, might need to check if individual types are
	// assignable.
	//
	// All implementations of Assignable must take into account the argument may implement this interface
	// do a reverse by calling the CheckAssignableTo function
	ReverseAssignable interface {
		// AssignableTo returns true if a variable or parameter of the other type can be hold a value of this type.
		// All implementations of Assignable must take into account that the given type might implement this method
		// do a reverse check before returning false.
		//
		// The guard is part of the internal endless recursion mechanism and should be passed as nil unless provided
		// by a DeepAssignable caller.
		AssignableTo(guard RecursionGuard, other Value) bool
	}
)
