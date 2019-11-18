package internal

import (
	"reflect"

	"github.com/lyraproj/dgo/dgo"
)

type nilValue int

func (n nilValue) Assignable(other interface{}) bool {
	_, ok := other.(nilValue)
	if !ok {
		ok = CheckAssignableTo(nil, other, n)
	}
	return ok
}

func (nilValue) AppendTo(w dgo.Indenter) {
	w.Append(`nil`)
}

func (nilValue) CompareTo(other interface{}) (int, bool) {
	if Nil == other || nil == other {
		return 0, true
	}
	return -1, true
}

func (nilValue) HashCode() int {
	return 131
}

func (nilValue) Equals(other interface{}) bool {
	return Nil == other || nil == other
}

func (nilValue) GoNil() interface{} {
	return nil
}

func (nilValue) ReflectTo(value reflect.Value) {
	value.Set(reflect.Zero(value.Type()))
}

func (nilValue) ReflectType() reflect.Type {
	return reflect.TypeOf((*dgo.Value)(nil)).Elem()
}

func (nilValue) String() string {
	return `null`
}

func (nilValue) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiNil
}

// Nil is the singleton dgo.Value for Nil
const Nil = nilValue(0)
