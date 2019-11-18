package internal

import (
	"reflect"

	"github.com/lyraproj/dgo/dgo"
)

type notType struct {
	negated dgo.Value
}

// DefaultNotType is the unconstrained Not type
var DefaultNotType = &notType{DefaultAnyType}

// NotType returns a type that represents all values that are not represented by the given type
func NotType(t dgo.Value) dgo.Value {
	// Avoid double negation
	if nt, ok := t.(*notType); ok {
		return nt.negated
	}
	return &notType{negated: t}
}

func (t *notType) Assignable(other interface{}) bool {
	switch ot := other.(type) {
	case *notType:
		// Reverse order of Negated test
		return ot.negated.Assignable(t.negated)
	case *anyOfType:
		ts := ot.slice
		for i := range ts {
			if t.Assignable(ts[i]) {
				return true
			}
		}
		return false
	case *allOfType:
		ts := ot.slice
		for i := range ts {
			if !t.Assignable(ts[i]) {
				return false
			}
		}
		return true
	case *oneOfType:
		f := false
		ts := ot.slice
		for i := range ts {
			if t.Assignable(ts[i]) {
				if f {
					return false
				}
				f = true
			}
		}
		return f
	default:
		return !t.negated.Assignable(other)
	}
}

func (t *notType) Equals(other interface{}) bool {
	if ot, ok := other.(*notType); ok {
		return t.negated.Equals(ot.negated)
	}
	return false
}

func (t *notType) HashCode() int {
	return 1579 + t.negated.HashCode()
}

func (t *notType) Operand() dgo.Value {
	return t.negated
}

func (t *notType) Operator() dgo.TypeOp {
	return dgo.OpNot
}

func (t *notType) ReflectType() reflect.Type {
	return reflectAnyType
}

func (t *notType) Resolve(ap dgo.AliasProvider) {
	tn := t.negated
	t.negated = DefaultAnyType
	t.negated = ap.Replace(tn)
}

func (t *notType) String() string {
	return TypeString(t)
}

func (t *notType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiNot
}
