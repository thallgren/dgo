package internal

import (
	"reflect"

	"github.com/lyraproj/dgo/dgo"
)

type (
	sensitive struct {
		wrapped dgo.Value
	}
)

// DefaultSensitiveType is the unconstrained Sensitive type
var DefaultSensitiveType = &sensitive{wrapped: DefaultAnyType}

// SensitiveType returns a Sensitive that wraps the given dgo.Value or, if args is of zero
// length, the default Sensitive type.
func SensitiveType(args []interface{}) dgo.Value {
	switch len(args) {
	case 0:
		return DefaultSensitiveType
	case 1:
		return &sensitive{wrapped: Value(args[0])}
	}
	panic(illegalArgumentCount(`SensitiveType`, 0, 1, len(args)))
}

// Sensitive creates a new Sensitive that wraps the given value
func Sensitive(v interface{}) dgo.Sensitive {
	return &sensitive{Value(v)}
}

func (v *sensitive) Assignable(other interface{}) bool {
	return Assignable(nil, v, other)
}

func (v *sensitive) DeepAssignable(guard dgo.RecursionGuard, other interface{}) bool {
	if ot, ok := other.(*sensitive); ok {
		return Assignable(guard, v.wrapped, ot.wrapped)
	}
	return CheckAssignableTo(guard, other, v)
}

func (v *sensitive) Equals(other interface{}) bool {
	return equals(nil, v, other)
}

func (v *sensitive) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ot, ok := other.(*sensitive); ok {
		return equals(seen, v.wrapped, ot.wrapped)
	}
	return false
}

func (v *sensitive) Freeze() {
	if f, ok := v.wrapped.(dgo.Freezable); ok {
		f.Freeze()
	}
}

func (v *sensitive) Frozen() bool {
	if f, ok := v.wrapped.(dgo.Freezable); ok {
		return f.Frozen()
	}
	return true
}

func (v *sensitive) FrozenCopy() dgo.Value {
	if f, ok := v.wrapped.(dgo.Freezable); ok && !f.Frozen() {
		return &sensitive{f.FrozenCopy()}
	}
	return v
}

func (v *sensitive) HashCode() int {
	return deepHashCode(nil, v)
}

func (v *sensitive) deepHashCode(seen []dgo.Value) int {
	return int(dgo.TiSensitive)*31 + deepHashCode(seen, v.wrapped)
}

var reflectSensitiveType = reflect.TypeOf((*dgo.Sensitive)(nil)).Elem()

func (v *sensitive) ReflectType() reflect.Type {
	return reflectSensitiveType
}

func (v *sensitive) Operand() dgo.Value {
	return v.wrapped
}

func (v *sensitive) Operator() dgo.TypeOp {
	return dgo.OpSensitive
}

func (v *sensitive) New(arg dgo.Value) dgo.Value {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`sensitive`, 1, 1)
		arg = args.Get(0)
	}
	if s, ok := arg.(dgo.Sensitive); ok {
		return s
	}
	return Sensitive(arg)
}

func (v *sensitive) String() string {
	return TypeString(v)
}

func (v *sensitive) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiSensitive
}

func (v *sensitive) Unwrap() dgo.Value {
	return v.wrapped
}
