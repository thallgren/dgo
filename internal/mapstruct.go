package internal

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/lyraproj/dgo/dgo"
)

type (
	// structType describes each mapEntry of a map
	structType struct {
		additional bool
		keys       array
		values     array
		required   []byte
	}

	structEntry struct {
		mapEntry
		required bool
	}
)

// StructMapType returns a new StructMapType type built from the given StructMapEntries.
func StructMapType(additional bool, entries []dgo.StructMapEntry) dgo.StructMapType {
	l := len(entries)
	keys := make([]dgo.Value, l)
	values := make([]dgo.Value, l)
	required := make([]byte, l)
	for i := range entries {
		e := entries[i]
		switch k := e.Key().(type) {
		case dgo.String, dgo.Integer, dgo.Float, dgo.Array, dgo.Map, dgo.Boolean, dgo.Binary:
			keys[i] = k
			values[i] = e.Value()
			if e.Required() {
				required[i] = 1
			}
		default:
			panic(`non exact key types is not yet supported`)
		}
	}
	return &structType{
		additional: additional,
		keys:       array{slice: keys, frozen: true},
		values:     array{slice: values, frozen: true},
		required:   required}
}

var sfmType dgo.MapType

// StructFromMapType returns the map type used when validating the map sent to
// StructMapTypeFromMap
func StructFromMapType() dgo.MapType {
	if sfmType == nil {
		sfmType = Parse(`map[string](dgo|type|{type:dgo|type,required?:bool,...})`).(dgo.MapType)
	}
	return sfmType
}

// StructMapTypeFromMap returns a new type built from a map[string](dgo|type|{type:dgo|type,required?:bool,...})
func StructMapTypeFromMap(additional bool, entries dgo.Map) dgo.StructMapType {
	if !StructFromMapType().Assignable(entries) {
		panic(IllegalAssignment(sfmType, entries))
	}
	l := entries.Len()
	keys := make([]dgo.Value, l)
	values := make([]dgo.Value, l)
	required := make([]byte, l)
	i := 0

	// turn dgo|type into type
	asType := func(v dgo.Value) dgo.Value {
		if s, ok := v.(dgo.String); ok {
			v = Parse(s.GoString())
		}
		return v
	}

	entries.EachEntry(func(e dgo.MapEntry) {
		keys[i] = e.Key()
		if vm, ok := e.Value().(dgo.Map); ok {
			values[i] = asType(vm.Get(`type`))
			if rq := vm.Get(`required`); rq != nil && rq.(dgo.Boolean).GoBool() {
				required[i] = 1
			}
		} else {
			values[i] = asType(e.Value())
		}
		i++
	})
	return &structType{
		additional: additional,
		keys:       array{slice: keys, frozen: true},
		values:     array{slice: values, frozen: true},
		required:   required}
}

func (t *structType) Additional() bool {
	return t.additional
}

func (t *structType) CheckEntry(k, v dgo.Value) dgo.Value {
	ks := t.keys.slice
	for i := range ks {
		kt := ks[i]
		if kt.Assignable(k) {
			vt := t.values.slice[i]
			if vt.Assignable(v) {
				return nil
			}
			return IllegalAssignment(vt, v)
		}
	}
	if t.additional {
		return nil
	}
	return IllegalMapKey(t, k)
}

func (t *structType) Assignable(other interface{}) bool {
	return Assignable(nil, t, other)
}

func (t *structType) DeepAssignable(guard dgo.RecursionGuard, other interface{}) bool {
	switch ot := other.(type) {
	case *structType:
		mrs := t.required
		mks := t.keys.slice
		mvs := t.values.slice
		ors := ot.required
		oks := ot.keys.slice
		ovs := ot.values.slice
		oc := 0

	nextKey:
		for mi := range mks {
			rq := mrs[mi] != 0
			mk := mks[mi]
			for oi := range oks {
				ok := oks[oi]
				if mk.Equals(ok) {
					if rq && ors[oi] == 0 {
						return false
					}
					if !Assignable(guard, mvs[mi], ovs[oi]) {
						return false
					}
					oc++
					continue nextKey
				}
			}
			if rq || ot.additional { // additional included since key is allowed with unconstrained value
				return false
			}
		}
		return t.additional || oc == len(oks)
	case dgo.Map:
		return t.deepInstance(guard, ot)
	case dgo.Value:
		return CheckAssignableTo(guard, other, t)
	default:
		if reflect.TypeOf(other).Kind() == reflect.Map {
			return t.deepInstance(guard, Value(other).(dgo.Map))
		}
		return false
	}
}

func (t *structType) Each(actor func(dgo.StructMapEntry)) {
	ks := t.keys.slice
	vs := t.values.slice
	rs := t.required
	for i := range ks {
		actor(&structEntry{mapEntry: mapEntry{key: ks[i], value: vs[i]}, required: rs[i] != 0})
	}
}

func (t *structType) Equals(other interface{}) bool {
	return equals(nil, t, other)
}

func (t *structType) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ot, ok := other.(*structType); ok {
		return t.additional == ot.additional &&
			bytes.Equal(t.required, ot.required) &&
			equals(seen, &t.keys, &ot.keys) &&
			equals(seen, &t.values, &ot.values)
	}
	return false
}

func (t *structType) Generic() dgo.Value {
	return newMapType(Generic(t.KeyType()), Generic(t.ValueType()), 0, math.MaxInt64)
}

func (t *structType) HashCode() int {
	return deepHashCode(nil, t)
}

func (t *structType) deepHashCode(seen []dgo.Value) int {
	h := bytesHash(t.required)*31 + t.keys.HashCode()*31 + t.values.HashCode()
	if t.additional {
		h *= 3
	}
	return h
}

func (t *structType) deepInstance(guard dgo.RecursionGuard, om dgo.Map) bool {
	ks := t.keys.slice
	vs := t.values.slice
	rs := t.required
	oc := 0
	for i := range ks {
		k := ks[i]
		if ov := om.Get(k); ov != nil {
			oc++
			if !Assignable(guard, vs[i], ov) {
				return false
			}
		} else if rs[i] != 0 {
			return false
		}
	}
	return t.additional || oc == om.Len()
}

func (t *structType) Get(key interface{}) dgo.MapEntry {
	kv := Value(key)
	i := t.keys.IndexOf(kv)
	if i >= 0 {
		return StructMapEntry(kv, t.values.slice[i], t.required[i] != 0)
	}
	return nil
}

func (t *structType) KeyType() dgo.Value {
	switch t.keys.Len() {
	case 0:
		return DefaultAnyType
	case 1:
		return t.keys.Get(0)
	default:
		return (*allOfType)(&t.keys)
	}
}

func (t *structType) Len() int {
	return len(t.required)
}

func (t *structType) Max() int {
	m := len(t.required)
	if m == 0 || t.additional {
		return math.MaxInt64
	}
	return m
}

func (t *structType) Min() int {
	min := 0
	rs := t.required
	for i := range rs {
		if rs[i] != 0 {
			min++
		}
	}
	return min
}

func (t *structType) New(arg dgo.Value) dgo.Value {
	return newMap(t, arg)
}

func (t *structType) ReflectType() reflect.Type {
	return reflect.MapOf(t.KeyType().ReflectType(), t.ValueType().ReflectType())
}

func (t *structType) Resolve(ap dgo.AliasProvider) {
	ks := t.keys.slice
	vs := t.values.slice
	t.keys.slice = []dgo.Value{}
	t.values.slice = []dgo.Value{}
	for i := range ks {
		ks[i] = ap.Replace(ks[i])
		vs[i] = ap.Replace(vs[i])
	}
	t.keys.slice = ks
	t.values.slice = vs
}

func (t *structType) String() string {
	return TypeString(t)
}

func (t *structType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiStruct
}

func (t *structType) Unbounded() bool {
	return len(t.required) == 0 || t.additional && t.Min() == 0
}

func parameterLabel(key dgo.Value) string {
	return fmt.Sprintf(`parameter '%s'`, key)
}

func (t *structType) Validate(keyLabel func(key dgo.Value) string, value interface{}) []error {
	var errs []error
	pm, ok := Value(value).(dgo.Map)
	if !ok {
		return []error{errors.New(`value is not a Map`)}
	}

	if keyLabel == nil {
		keyLabel = parameterLabel
	}
	t.Each(func(e dgo.StructMapEntry) {
		ek := e.Key()
		if v := pm.Get(ek); v != nil {
			ev := e.Value()
			if !ev.Assignable(v) {
				errs = append(errs, fmt.Errorf(`%s is not an instance of type %s`, keyLabel(ek), ev))
			}
		} else if e.Required() {
			errs = append(errs, fmt.Errorf(`missing required %s`, keyLabel(ek)))
		}
	})
	pm.EachKey(func(k dgo.Value) {
		if t.Get(k) == nil {
			errs = append(errs, fmt.Errorf(`unknown %s`, keyLabel(k)))
		}
	})
	return errs
}

func (t *structType) ValidateVerbose(value interface{}, out dgo.Indenter) bool {
	pm, ok := Value(value).(dgo.Map)
	if !ok {
		out.Append(`value is not a Map`)
		return false
	}

	inner := out.Indent()
	t.Each(func(e dgo.StructMapEntry) {
		ek := e.Key()
		ev := e.Value()
		out.Printf(`Validating '%s' against definition %s`, ek, ev)
		inner.NewLine()
		inner.Printf(`'%s' `, ek)
		if v := pm.Get(ek); v != nil {
			if ev.Assignable(v) {
				inner.Append(`OK!`)
			} else {
				ok = false
				inner.Append(`FAILED!`)
				inner.NewLine()
				inner.Printf(`Reason: expected a value of type %s, got %s`, ev, v)
			}
		} else if e.Required() {
			ok = false
			inner.Append(`FAILED!`)
			inner.NewLine()
			inner.Append(`Reason: required key not found in input`)
		}
		out.NewLine()
	})
	pm.EachKey(func(k dgo.Value) {
		if t.Get(k) == nil {
			ok = false
			out.Printf(`Validating '%s'`, k)
			inner.NewLine()
			inner.Printf(`'%s' FAILED!`, k)
			inner.NewLine()
			inner.Append(`Reason: key is not found in definition`)
			out.NewLine()
		}
	})
	return ok
}

func (t *structType) Value() dgo.Value {
	ks := t.keys.slice
	vs := t.values.slice
	m := MapWithCapacity(len(ks), nil)
	for i := range ks {
		m.Put(ks[i], vs[i])
	}
	return m
}

func (t *structType) ValueType() dgo.Value {
	switch t.values.Len() {
	case 0:
		return DefaultAnyType
	case 1:
		return t.values.Get(0)
	default:
		return (*allOfType)(&t.values)
	}
}

// StructMapEntry returns a new StructMapEntry initiated with the given parameters
func StructMapEntry(key interface{}, value interface{}, required bool) dgo.StructMapEntry {
	return &structEntry{mapEntry: mapEntry{key: Value(key), value: Value(value)}, required: required}
}

func (t *structEntry) Equals(other interface{}) bool {
	return equals(nil, t, other)
}

func (t *structEntry) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ot, ok := other.(*structEntry); ok {
		return t.required == ot.required && equals(seen, &t.mapEntry, &ot.mapEntry)
	}
	return false
}

func (t *structEntry) Required() bool {
	return t.required
}
