package internal

import (
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/lyraproj/dgo/util"

	"github.com/lyraproj/dgo/dgo"
)

type (
	// floatVal is a float64 that implements the dgo.Value interface
	floatVal float64

	floatType int

	floatRangeType struct {
		min       float64
		max       float64
		inclusive bool
	}
)

// DefaultFloatType is the unconstrained floatVal type
const DefaultFloatType = floatType(0)

var reflectFloatType = reflect.TypeOf(float64(0))

// FloatRangeType returns a dgo.FloatRangeType that is limited to the inclusive range given by min and max
// If inclusive is true, then the range has an inclusive end.
func FloatRangeType(min, max float64, inclusive bool) dgo.FloatRangeType {
	if min == max {
		if !inclusive {
			panic(fmt.Errorf(`non inclusive range cannot have equal min and max`))
		}
		return floatVal(min)
	}
	if max < min {
		t := max
		max = min
		min = t
	}
	if min == -math.MaxFloat64 && max == math.MaxFloat64 {
		return DefaultFloatType
	}
	return &floatRangeType{min: min, max: max, inclusive: inclusive}
}

func (t *floatRangeType) Assignable(other interface{}) bool {
	switch ot := other.(type) {
	case floatVal:
		return t.IsInstance(float64(ot))
	case *floatRangeType:
		if t.min > ot.min {
			return false
		}
		if t.inclusive || t.inclusive == ot.inclusive {
			return t.max >= ot.max
		}
		return t.max > ot.max
	default:
		if fv, ok := ToFloat(other); ok {
			return t.IsInstance(fv)
		}
		return CheckAssignableTo(nil, other, t)
	}
}

func (t *floatRangeType) Equals(other interface{}) bool {
	if ot, ok := other.(*floatRangeType); ok {
		return *t == *ot
	}
	return false
}

func (t *floatRangeType) HashCode() int {
	h := int(dgo.TiFloatRange)
	if t.min > 0 {
		h = h*31 + int(t.min)
	}
	if t.max < math.MaxInt64 {
		h = h*31 + int(t.max)
	}
	if t.inclusive {
		h *= 3
	}
	return h
}

func (t *floatRangeType) IsInstance(value float64) bool {
	if t.min <= value {
		if t.inclusive {
			return value <= t.max
		}
		return value < t.max
	}
	return false
}

func (t *floatRangeType) Max() float64 {
	return t.max
}

func (t *floatRangeType) Inclusive() bool {
	return t.inclusive
}

func (t *floatRangeType) Min() float64 {
	return t.min
}

func (t *floatRangeType) New(arg dgo.Value) dgo.Value {
	return newFloat(t, arg)
}

func (t *floatRangeType) String() string {
	return TypeString(t)
}

func (t *floatRangeType) ReflectType() reflect.Type {
	return reflectFloatType
}

func (t *floatRangeType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiFloatRange
}

func (t floatType) Assignable(other interface{}) bool {
	switch other := other.(type) {
	case dgo.FloatRangeType:
		return true
	default:
		_, ok := ToFloat(other)
		return ok || CheckAssignableTo(nil, other, t)
	}
}

func (t floatType) Equals(other interface{}) bool {
	_, ok := other.(floatType)
	return ok
}

func (t floatType) HashCode() int {
	return int(dgo.TiFloat)
}

func (t floatType) Inclusive() bool {
	return true
}

func (t floatType) Instance(value interface{}) bool {
	_, ok := ToFloat(value)
	return ok
}

func (t floatType) IsInstance(value float64) bool {
	return true
}

func (t floatType) Max() float64 {
	return math.MaxFloat64
}

func (t floatType) Min() float64 {
	return -math.MaxFloat64
}

func (t floatType) New(arg dgo.Value) dgo.Value {
	return newFloat(t, arg)
}

func (t floatType) ReflectType() reflect.Type {
	return reflectFloatType
}

func (t floatType) String() string {
	return TypeString(t)
}

func (t floatType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiFloat
}

// Float returns the dgo.Float for the given float64
func Float(f float64) dgo.Float {
	return floatVal(f)
}

func (v floatVal) Assignable(other interface{}) bool {
	return v.Equals(other) || CheckAssignableTo(nil, other, v)
}

func (v floatVal) CompareTo(other interface{}) (r int, ok bool) {
	ok = true
	if ov, isFloat := ToFloat(other); isFloat {
		fv := float64(v)
		switch {
		case fv > ov:
			r = 1
		case fv < ov:
			r = -1
		default:
			r = 0
		}
		return
	}
	if oi, isInt := ToInt(other); isInt {
		fv := float64(v)
		ov := float64(oi)
		switch {
		case fv > ov:
			r = 1
		case fv < ov:
			r = -1
		default:
			r = 0
		}
		return
	}
	if other == Nil || other == nil {
		r = 1
	} else {
		ok = false
	}
	return
}

func (v floatVal) Equals(other interface{}) bool {
	f, ok := ToFloat(other)
	return ok && float64(v) == f
}

func (v floatVal) Generic() dgo.Value {
	return DefaultFloatType
}

func (v floatVal) GoFloat() float64 {
	return float64(v)
}

func (v floatVal) HashCode() int {
	return int(v)
}

func (v floatVal) Inclusive() bool {
	return true
}

func (v floatVal) IsInstance(value float64) bool {
	return float64(v) == value
}

func (v floatVal) Max() float64 {
	return float64(v)
}

func (v floatVal) Min() float64 {
	return float64(v)
}

func (v floatVal) New(arg dgo.Value) dgo.Value {
	return newFloat(v, arg)
}

func (v floatVal) ReflectTo(value reflect.Value) {
	switch value.Kind() {
	case reflect.Interface:
		value.Set(reflect.ValueOf(float64(v)))
	case reflect.Ptr:
		if value.Type().Elem().Kind() == reflect.Float32 {
			gv := float32(v)
			value.Set(reflect.ValueOf(&gv))
		} else {
			gv := float64(v)
			value.Set(reflect.ValueOf(&gv))
		}
	default:
		value.SetFloat(float64(v))
	}
}

func (v floatVal) ReflectType() reflect.Type {
	return reflectFloatType
}

func (v floatVal) String() string {
	return util.Ftoa(float64(v))
}

func (v floatVal) ToFloat() float64 {
	return float64(v)
}

func (v floatVal) ToInt() int64 {
	return int64(v)
}

func (v floatVal) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiFloatExact
}

// ToFloat returns the given value as a float64 if, and only if, the value is a float32 or float64. An
// additional boolean is returned to indicate if that was the case or not.
func ToFloat(value interface{}) (v float64, ok bool) {
	ok = true
	switch value := value.(type) {
	case floatVal:
		v = float64(value)
	case float64:
		v = value
	case float32:
		v = float64(value)
	default:
		ok = false
	}
	return
}

func newFloat(t dgo.Value, arg dgo.Value) (f dgo.Float) {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`float`, 1, 1)
		arg = args.Get(0)
	}
	f = Float(floatFromConvertible(arg))
	if !t.Assignable(f) {
		panic(IllegalAssignment(t, f))
	}
	return f
}

func floatFromConvertible(from dgo.Value) float64 {
	switch from := from.(type) {
	case dgo.Float:
		return from.GoFloat()
	case dgo.Integer:
		return float64(from.GoInt())
	case *timeVal:
		return from.SecondsWithFraction()
	case dgo.Boolean:
		if from.GoBool() {
			return 1
		}
		return 0
	case dgo.String:
		if i, err := strconv.ParseFloat(from.GoString(), 64); err == nil {
			return i
		}
	}
	panic(fmt.Errorf(`the value '%s' cannot be converted to a float`, from))
}
