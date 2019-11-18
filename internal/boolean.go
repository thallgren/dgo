package internal

import (
	"fmt"
	"reflect"

	"github.com/lyraproj/dgo/dgo"
)

type (
	// booleanType represents an boolean without constraints (-1), constrained to false (0) or constrained to true(1)
	booleanType int

	boolean bool
)

// DefaultBooleanType is the unconstrained Boolean type
const DefaultBooleanType = booleanType(0)

// True is the dgo.Value for true
const True = boolean(true)

// False is the dgo.Value for false
const False = boolean(false)

func (t booleanType) Assignable(other interface{}) bool {
	switch other := other.(type) {
	case dgo.BooleanType, []byte:
		return true
	default:
		return CheckAssignableTo(nil, other, t)
	}
}

func (t booleanType) Equals(v interface{}) bool {
	return t == v
}

func (t booleanType) HashCode() int {
	return int(t.TypeIdentifier())
}

func (t booleanType) IsInstance(v bool) bool {
	return true
}

var boolStringType = CiEnumType([]string{`true`, `false`, `yes`, `no`, `y`, `n`})
var trueStringType = CiEnumType([]string{`true`, `yes`, `y`})

func (t booleanType) New(arg dgo.Value) dgo.Value {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`bool`, 1, 1)
		arg = args.Get(0)
	}
	switch arg := arg.(type) {
	case boolean:
		return arg
	case intVal:
		return boolean(arg != 0)
	case floatVal:
		return boolean(arg != 0)
	default:
		if boolStringType.Assignable(arg) {
			return boolean(trueStringType.Assignable(arg))
		}
		panic(fmt.Errorf(`unable to create a bool from %s`, arg))
	}
}

var reflectBooleanType = reflect.TypeOf(true)

func (t booleanType) ReflectType() reflect.Type {
	return reflectBooleanType
}

func (t booleanType) String() string {
	return TypeString(t)
}

func (t booleanType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiBoolean
}

func (v boolean) Assignable(other interface{}) bool {
	switch other := other.(type) {
	case dgo.Boolean:
		return bool(v) == other.GoBool()
	case bool:
		return bool(v) == other
	default:
		return CheckAssignableTo(nil, other, v)
	}
}

func (v boolean) GoBool() bool {
	return bool(v)
}

func (v boolean) CompareTo(other interface{}) (r int, ok bool) {
	ok = true
	switch ov := other.(type) {
	case boolean:
		r = 0
		if v {
			if !ov {
				r = 1
			}
		} else if ov {
			r = -1
		}
	case nilValue:
		r = 1
	default:
		ok = false
	}
	return
}

func (v boolean) Equals(other interface{}) bool {
	if ov, ok := other.(boolean); ok {
		return v == ov
	}
	if ov, ok := other.(bool); ok {
		return bool(v) == ov
	}
	return false
}

func (v boolean) HashCode() int {
	if v {
		return 1231
	}
	return 1237
}

func (v boolean) IsInstance(b bool) bool {
	return bool(v) == b
}

func (v boolean) ReflectTo(value reflect.Value) {
	b := bool(v)
	switch value.Kind() {
	case reflect.Interface:
		value.Set(reflect.ValueOf(b))
	case reflect.Ptr:
		value.Set(reflect.ValueOf(&b))
	default:
		value.SetBool(b)
	}
}

func (v boolean) ReflectType() reflect.Type {
	return reflectBooleanType
}

func (v boolean) String() string {
	if v {
		return `true`
	}
	return `false`
}

func (v boolean) TypeIdentifier() dgo.TypeIdentifier {
	if v {
		return dgo.TiTrue
	}
	return dgo.TiFalse
}
