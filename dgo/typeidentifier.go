package dgo

import "fmt"

// TypeIdentifier is a unique identifier for each type known to the system. The order of the TypeIdentifier
// determines the sort order for elements that are not comparable
type TypeIdentifier int

const (
	// TiNil is the type identifier for the Nil type
	TiNil = TypeIdentifier(iota)

	// TiAny is the type identifier for the Any type
	TiAny

	// TiMeta is the type identifier for the Meta type
	TiMeta

	// TiBoolean is the type identifier for the Boolean type
	TiBoolean

	// TiFalse is the type identifier for the False type
	TiFalse

	// TiTrue is the type identifier for the True type
	TiTrue

	// TiInteger is the type identifier for the Integer type
	TiInteger

	// TiIntegerExact is the type identifier for the exact Integer type
	TiIntegerExact

	// TiIntegerRange is the type identifier for the Integer range type
	TiIntegerRange

	// TiFloat is the type identifier for the Float type
	TiFloat

	// TiFloatExact is the type identifier for the exact Float type
	TiFloatExact

	// TiFloatRange is the type identifier for the Float range type
	TiFloatRange

	// TiBinary is the type identifier for the Binary type
	TiBinary

	// TiBinaryExact is the type identifier for the exact Binary type
	TiBinaryExact

	// TiString is the type identifier for the String type
	TiString

	// TiStringExact is the type identifier for the exact String type
	TiStringExact

	// TiStringSized is the type identifier for the size constrained String type
	TiStringSized

	// TiStringPattern is the type identifier for the String pattern type
	TiStringPattern

	// TiCiString is the type identifier for the case insensitive String type
	TiCiString

	// TiRegexp is the type identifier for the Regexp type
	TiRegexp

	// TiRegexpExact is the type identifier for the exact Regexp type
	TiRegexpExact

	// TiNative is the type identifier for the Native type
	TiNative

	// TiNativeExact is the type identifier for the exact Native type
	TiNativeExact

	// TiArray is the type identifier for the Array type
	TiArray

	// TiArrayExact is the type identifier for the exact Array type
	TiArrayExact

	// TiTuple is the type identifier for the Tuple type
	TiTuple

	// TiMap is the type identifier for the Map type
	TiMap

	// TiMapExact is the type identifier for exact Map type
	TiMapExact

	// TiMapEntryExact is the type identifier the map entry type of the exact Map type
	TiMapEntryExact

	// TiStruct is the type identifier for the Struct type
	TiStruct

	// TiNot is the type identifier for the Not type
	TiNot

	// TiAllOf is the type identifier for the AllOf type
	TiAllOf

	// TiAllOfValue is the type identifier for the AllOf type that uses the type of its contained values
	TiAllOfValue

	// TiAnyOf is the type identifier for the AnyOf type
	TiAnyOf

	// TiOneOf is the type identifier for the OneOf type
	TiOneOf

	// TiError is the type identifier for for the Error type
	TiError

	// TiErrorExact is the type identifier for for the exact Error type
	TiErrorExact

	// TiDgoString is the type identifier for for the DgoString type
	TiDgoString

	// TiSensitive is the type identifier for for the Sensitive type
	TiSensitive

	// TiTime is the type identifier for for the Time type
	TiTime
	// TiTimeExact is the type identifier for the exact Time type
	TiTimeExact

	// TiFunction is the type identifier for for the Function type
	TiFunction

	// TiNamed is the type identifier for for named types
	TiNamed

	// TiNamedExact is the type identifier for for exact Named types
	TiNamedExact
)

var tiLabels = map[TypeIdentifier]string{
	TiNil:           `nil`,
	TiAny:           `any`,
	TiMeta:          `type`,
	TiBoolean:       `bool`,
	TiFalse:         `false`,
	TiTrue:          `true`,
	TiInteger:       `int`,
	TiIntegerExact:  `int`,
	TiIntegerRange:  `int range`,
	TiFloat:         `float`,
	TiFloatExact:    `float`,
	TiFloatRange:    `float range`,
	TiBinary:        `binary`,
	TiBinaryExact:   `binary`,
	TiString:        `string`,
	TiStringExact:   `string`,
	TiStringSized:   `string`,
	TiStringPattern: `pattern`,
	TiCiString:      `string`,
	TiRegexp:        `regexp`,
	TiRegexpExact:   `regexp`,
	TiTime:          `time`,
	TiTimeExact:     `time`,
	TiNative:        `native`,
	TiNativeExact:   `native`,
	TiArray:         `slice`,
	TiArrayExact:    `slice`,
	TiTuple:         `tuple`,
	TiMap:           `map`,
	TiMapExact:      `map`,
	TiMapEntryExact: `map entry`,
	TiStruct:        `struct`,
	TiNot:           `not`,
	TiAllOf:         `all of`,
	TiAllOfValue:    `all of`,
	TiAnyOf:         `any of`,
	TiOneOf:         `one of`,
	TiError:         `error`,
	TiErrorExact:    `error`,
	TiDgoString:     `dgo`,
	TiSensitive:     `sensitive`,
	TiFunction:      `function`,
	TiNamed:         `named`,
}

func (ti TypeIdentifier) String() string {
	if s, ok := tiLabels[ti]; ok {
		return s
	}
	panic(fmt.Errorf("unhandled TypeIdentifier %d", ti))
}
