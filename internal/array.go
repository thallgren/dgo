package internal

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"

	"github.com/lyraproj/dgo/dgo"
	"github.com/lyraproj/dgo/util"
)

type (
	array struct {
		slice  []dgo.Value
		typ    dgo.ArrayType
		frozen bool
	}

	// defaultArrayType is the unconstrained array type
	defaultArrayType int

	// sizedArrayType represents array with element type constraint and a size constraint
	sizedArrayType struct {
		elementType dgo.Value
		min         int
		max         int
	}

	// tupleType represents an array with an exact number of ordered element types.
	tupleType struct {
		types    []dgo.Value
		variadic bool
	}
)

// DefaultArrayType is the unconstrained Array type
const DefaultArrayType = defaultArrayType(0)

func arrayTypeOne(args []interface{}) dgo.ArrayType {
	switch a0 := Value(args[0]).(type) {
	case dgo.Integer:
		return newArrayType(nil, int(a0.GoInt()), math.MaxInt64)
	case dgo.Value:
		return newArrayType(a0, 0, math.MaxInt64)
	default:
		panic(illegalArgument(`Array`, `Type or Integer`, args, 0))
	}
}

func arrayTypeTwo(args []interface{}) dgo.ArrayType {
	a1, ok := Value(args[1]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`Array`, `Integer`, args, 1))
	}
	switch a0 := Value(args[0]).(type) {
	case dgo.Integer:
		return newArrayType(nil, int(a0.GoInt()), int(a1.GoInt()))
	case dgo.Value:
		return newArrayType(a0, int(a1.GoInt()), math.MaxInt64)
	default:
		panic(illegalArgument(`Array`, `Type or Integer`, args, 0))
	}
}

func arrayTypeThree(args []interface{}) dgo.ArrayType {
	a0 := Value(args[0])
	a1, ok := Value(args[1]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`Array`, `Integer`, args, 1))
	}
	a2, ok := Value(args[2]).(dgo.Integer)
	if !ok {
		panic(illegalArgument(`ArrayType`, `Integer`, args, 2))
	}
	return newArrayType(a0, int(a1.GoInt()), int(a2.GoInt()))
}

// ArrayType returns a type that represents an Array value
func ArrayType(args []interface{}) dgo.ArrayType {
	switch len(args) {
	case 0:
		return DefaultArrayType
	case 1:
		return arrayTypeOne(args)
	case 2:
		return arrayTypeTwo(args)
	case 3:
		return arrayTypeThree(args)
	default:
		panic(illegalArgumentCount(`Array`, 0, 3, len(args)))
	}
}

func newArrayType(elementType dgo.Value, min, max int) dgo.ArrayType {
	if min < 0 {
		min = 0
	}
	if max < 0 {
		max = 0
	}
	if max < min {
		t := max
		max = min
		min = t
	}
	if elementType == nil {
		elementType = DefaultAnyType
	}
	if min == 0 && max == math.MaxInt64 && elementType == DefaultAnyType {
		// Unbounded
		return DefaultArrayType
	}
	return &sizedArrayType{elementType: elementType, min: min, max: max}
}

func (t defaultArrayType) Assignable(other interface{}) bool {
	switch other.(type) {
	case dgo.ArrayType:
		return true
	}
	return CheckAssignableTo(nil, other, t)
}

func (t defaultArrayType) ElementType() dgo.Value {
	return DefaultAnyType
}

func (t defaultArrayType) Equals(other interface{}) bool {
	return t == other
}

func (t defaultArrayType) HashCode() int {
	return int(dgo.TiArray)
}

func (t defaultArrayType) Max() int {
	return math.MaxInt64
}

func (t defaultArrayType) Min() int {
	return 0
}

func (t defaultArrayType) New(arg dgo.Value) dgo.Value {
	return newArray(t, arg)
}

func (t defaultArrayType) ReflectType() reflect.Type {
	return reflect.SliceOf(reflectAnyType)
}

func (t defaultArrayType) String() string {
	return TypeString(t)
}

func (t defaultArrayType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiArray
}

func (t defaultArrayType) Unbounded() bool {
	return true
}

func (t *sizedArrayType) Assignable(other interface{}) bool {
	return Assignable(nil, t, other)
}

func (t *sizedArrayType) DeepAssignable(guard dgo.RecursionGuard, other interface{}) bool {
	switch ot := other.(type) {
	case defaultArrayType:
		return false // lacks size
	case dgo.ArrayType:
		return t.min <= ot.Min() && ot.Max() <= t.max && t.elementType.Assignable(ot.ElementType())
	}
	return CheckAssignableTo(guard, other, t)
}

func (t *sizedArrayType) ElementType() dgo.Value {
	return t.elementType
}

func (t *sizedArrayType) Equals(other interface{}) bool {
	return equals(nil, t, other)
}

func (t *sizedArrayType) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ot, ok := other.(*sizedArrayType); ok {
		return t.min == ot.min && t.max == ot.max && equals(seen, t.elementType, ot.elementType)
	}
	return false
}

func (t *sizedArrayType) HashCode() int {
	return deepHashCode(nil, t)
}

func (t *sizedArrayType) deepHashCode(seen []dgo.Value) int {
	h := int(dgo.TiArray)
	if t.min > 0 {
		h = h*31 + t.min
	}
	if t.max < math.MaxInt64 {
		h = h*31 + t.max
	}
	if DefaultAnyType != t.elementType {
		h = h*31 + deepHashCode(seen, t.elementType)
	}
	return h
}

func (t *sizedArrayType) DeepInstance(guard dgo.RecursionGuard, value interface{}) bool {
	if ov, ok := value.(*array); ok {
		l := len(ov.slice)
		return t.min <= l && l <= t.max && allAssignable(guard, t.elementType, ov.slice)
	}
	return false
}

func (t *sizedArrayType) Max() int {
	return t.max
}

func (t *sizedArrayType) Min() int {
	return t.min
}

func (t *sizedArrayType) New(arg dgo.Value) dgo.Value {
	return newArray(t, arg)
}

func (t *sizedArrayType) Resolve(ap dgo.AliasProvider) {
	te := t.elementType
	t.elementType = DefaultAnyType
	t.elementType = ap.Replace(te)
}

func (t *sizedArrayType) ReflectType() reflect.Type {
	return reflect.SliceOf(t.elementType.ReflectType())
}

func (t *sizedArrayType) String() string {
	return TypeString(t)
}

func (t *sizedArrayType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiArray
}

func (t *sizedArrayType) Unbounded() bool {
	return t.min == 0 && t.max == math.MaxInt64
}

// DefaultTupleType is a tuple without size and type constraints
var DefaultTupleType = &tupleType{variadic: true, types: []dgo.Value{DefaultAnyType}}

// EmptyTupleType is a tuple that represents an empty array
var EmptyTupleType = &tupleType{variadic: false, types: []dgo.Value{}}

// TupleType creates a new TupleTupe based on the given types
func TupleType(types []interface{}) dgo.TupleType {
	return newTupleType(types, false)
}

// VariadicTupleType returns a type that represents an Array value with a variadic number of elements. Each
// given type determines the type of a corresponding element in an array except for the last one which
// determines the remaining elements.
func VariadicTupleType(types []interface{}) dgo.TupleType {
	n := len(types)
	if n == 0 {
		panic(errors.New(`a variadic tuple must have at least one element`))
	}
	return newTupleType(types, true)
}

func newTupleType(types []interface{}, variadic bool) dgo.TupleType {
	l := len(types)
	if l == 0 {
		return EmptyTupleType
	}
	if variadic && l == 1 && DefaultAnyType.Equals(types[0]) {
		return DefaultTupleType
	}
	es := make([]dgo.Value, l)
	for i := range types {
		es[i] = types[i].(dgo.Value)
	}
	return &tupleType{types: es, variadic: variadic}
}

func (t *tupleType) Assignable(other interface{}) bool {
	return Assignable(nil, t, other)
}

func (t *tupleType) DeepAssignable(guard dgo.RecursionGuard, other interface{}) bool {
	return tupleAssignable(guard, t, other)
}

func tupleAssignableTuple(guard dgo.RecursionGuard, t, ot dgo.TupleType) bool {
	if t.Min() > ot.Min() || ot.Max() > t.Max() {
		return false
	}

	var tv, ov dgo.Value
	tn := t.Len()
	if t.Variadic() {
		tn--
		tv = t.Element(tn)
	}
	on := ot.Len()
	if ot.Variadic() {
		on--
		ov = ot.Element(on)
	}

	// n := max(tn, on)
	n := tn
	if n < on {
		n = on
	}

	for i := 0; i < n; i++ {
		te := tv
		if i < tn {
			te = t.Element(i)
		}

		oe := ov
		if i < on {
			oe = ot.Element(i)
		}
		if te == nil || oe == nil || !Assignable(guard, te, oe) {
			return false
		}
	}
	return true
}

func tupleAssignableArray(guard dgo.RecursionGuard, t dgo.TupleType, ot *sizedArrayType) bool {
	if t.Min() <= ot.Min() && ot.Max() <= t.Max() {
		et := ot.ElementType()
		n := t.Len()
		if t.Variadic() {
			n--
		}
		for i := 0; i < n; i++ {
			if !Assignable(guard, t.Element(i), et) {
				return false
			}
		}
		return !t.Variadic() || Assignable(guard, t.Element(n), et)
	}
	return false
}

func tupleAssignable(guard dgo.RecursionGuard, t dgo.TupleType, other interface{}) bool {
	switch ot := other.(type) {
	case defaultArrayType:
		return false
	case dgo.Array:
		return tupleInstance(guard, t, ot)
	case dgo.TupleType:
		return tupleAssignableTuple(guard, t, ot)
	case *sizedArrayType:
		return tupleAssignableArray(guard, t, ot)
	}
	return CheckAssignableTo(guard, other, t)
}

func (t *tupleType) Element(index int) dgo.Value {
	return t.types[index]
}

func (t *tupleType) ElementType() dgo.Value {
	return tupleElementType(t)
}

func tupleElementType(t dgo.TupleType) (et dgo.Value) {
	switch t.Len() {
	case 0:
		et = DefaultAnyType
	case 1:
		et = t.Element(0)
	default:
		ea := t.ElementTypes().Unique()
		if ea.Len() == 1 {
			return ea.Get(0)
		}
		return (*allOfType)(ea.(*array))
	}
	return
}

func (t *tupleType) ElementTypes() dgo.Array {
	return &array{slice: t.types, frozen: true}
}

func (t *tupleType) Equals(other interface{}) bool {
	return equals(nil, t, other)
}

func (t *tupleType) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ot, ok := other.(*tupleType); ok {
		return t.variadic == ot.variadic && sliceEquals(seen, t.types, ot.types)
	}
	return tupleEquals(seen, t, other)
}

func tupleEquals(seen []dgo.Value, t dgo.TupleType, other interface{}) bool {
	if ot, ok := other.(dgo.TupleType); ok {
		n := t.Len()
		if t.Variadic() == ot.Variadic() && n == ot.Len() {
			for i := 0; i < n; i++ {
				if !equals(seen, t.Element(i), ot.Element(i)) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (t *tupleType) Generic() dgo.Value {
	return newArrayType(Generic(t.ElementType()), 0, math.MaxInt64)
}

func (t *tupleType) HashCode() int {
	return tupleHashCode(t, nil)
}

func (t *tupleType) deepHashCode(seen []dgo.Value) int {
	return tupleHashCode(t, seen)
}

func tupleHashCode(t dgo.TupleType, seen []dgo.Value) int {
	h := 1
	if t.Variadic() {
		h = 7
	}
	l := t.Len()
	for i := 0; i < l; i++ {
		h = h*31 + deepHashCode(seen, t.Element(i))
	}
	return h
}

func tupleInstance(guard dgo.RecursionGuard, t dgo.TupleType, value interface{}) bool {
	if ov, ok := value.(*array); ok {
		s := ov.slice
		n := len(s)
		if t.Variadic() {
			if t.Min() <= n {
				tn := t.Len() - 1
				for i := 0; i < tn; i++ {
					if !Assignable(guard, t.Element(i), s[i]) {
						return false
					}
				}
				vt := t.Element(tn)
				for ; tn < n; tn++ {
					if !Assignable(guard, vt, s[tn]) {
						return false
					}
				}
				return true
			}
		} else {
			if n == t.Len() {
				for i := range s {
					if !Assignable(guard, t.Element(i), s[i]) {
						return false
					}
				}
				return true
			}
		}
	}
	return false
}

func (t *tupleType) Len() int {
	return len(t.types)
}

func (t *tupleType) Max() int {
	return tupleMax(t)
}

func tupleMax(t dgo.TupleType) int {
	n := t.Len()
	if t.Variadic() {
		n = math.MaxInt64
	}
	return n
}

func (t *tupleType) Min() int {
	return tupleMin(t)
}

func (t *tupleType) New(arg dgo.Value) dgo.Value {
	return newArray(t, arg)
}

func tupleMin(t dgo.TupleType) int {
	n := t.Len()
	if t.Variadic() {
		n--
	}
	return n
}

func (t *tupleType) ReflectType() reflect.Type {
	return reflect.SliceOf(t.ElementType().ReflectType())
}

func (t *tupleType) Resolve(ap dgo.AliasProvider) {
	s := t.types
	t.types = nil
	resolveSlice(s, ap)
	t.types = s
}

func (t *tupleType) String() string {
	return TypeString(t)
}

func (t *tupleType) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiTuple
}

func (t *tupleType) Unbounded() bool {
	return t.variadic
}

func (t *tupleType) Value() dgo.Value {
	a := util.SliceCopy(t.types)
	return &array{slice: a, frozen: true}
}

func (t *tupleType) Variadic() bool {
	return t.variadic
}

// Array returns a frozen dgo.Array that represents a copy of the given value. The value can be
// a slice or an Iterable
func Array(value interface{}) dgo.Array {
	switch value := value.(type) {
	case dgo.Array:
		return value.FrozenCopy().(dgo.Array)
	case dgo.Iterable:
		return arrayFromIterator(value.Len(), value.Each)
	case []dgo.Value:
		arr := make([]dgo.Value, len(value))
		for i := range value {
			e := value[i]
			if f, ok := e.(dgo.Freezable); ok {
				e = f.FrozenCopy()
			} else if e == nil {
				e = Nil
			}
			arr[i] = e
		}
		return &array{slice: arr, frozen: true}
	case reflect.Value:
		return ValueFromReflected(value).(dgo.Array)
	default:
		return ValueFromReflected(reflect.ValueOf(value)).(dgo.Array)
	}
}

// arrayFromIterator creates an array from a size and an iterator goFunc. The
// iterator goFunc is expected to call its actor exactly size number of times.
func arrayFromIterator(size int, each func(dgo.Consumer)) *array {
	arr := make([]dgo.Value, size)
	i := 0
	each(func(e dgo.Value) {
		if f, ok := e.(dgo.Freezable); ok {
			e = f.FrozenCopy()
		}
		arr[i] = e
		i++
	})
	return &array{slice: arr, frozen: true}
}

func sliceFromIterable(ir dgo.Iterable) []dgo.Value {
	es := make([]dgo.Value, ir.Len())
	i := 0
	ir.Each(func(e dgo.Value) {
		es[i] = e
		i++
	})
	return es
}

// ArrayFromReflected creates a new array that contains a copy of the given reflected slice
func ArrayFromReflected(vr reflect.Value, frozen bool) dgo.Value {
	if vr.IsNil() {
		return Nil
	}

	var arr []dgo.Value
	if vr.CanInterface() {
		ix := vr.Interface()
		if bs, ok := ix.([]byte); ok {
			return Binary(bs, frozen)
		}

		if vs, ok := ix.([]dgo.Value); ok {
			arr = vs
			if frozen {
				arr = util.SliceCopy(arr)
			}
		}
	}

	if arr == nil {
		top := vr.Len()
		arr = make([]dgo.Value, top)
		for i := 0; i < top; i++ {
			arr[i] = ValueFromReflected(vr.Index(i))
		}
	}

	if frozen {
		for i := range arr {
			if f, ok := arr[i].(dgo.Freezable); ok {
				arr[i] = f.FrozenCopy()
			}
		}
	}
	return &array{slice: arr, frozen: frozen}
}

func asArrayType(typ interface{}) dgo.ArrayType {
	if typ == nil {
		return nil
	}

	parseArrayType := func(s string) dgo.ArrayType {
		if t, ok := Parse(s).(dgo.ArrayType); ok {
			return t
		}
		panic(fmt.Errorf("expression '%s' does not evaluate to an array type", s))
	}

	var mt dgo.ArrayType
	switch typ := typ.(type) {
	case dgo.ArrayType:
		mt = typ
	case string:
		mt = parseArrayType(typ)
	case dgo.String:
		mt = parseArrayType(typ.GoString())
	default:
		mt = TypeFromReflected(reflect.TypeOf(typ)).(dgo.ArrayType)
	}
	return mt
}

// ArrayWithCapacity creates a new mutable array of the given type and initial capacity. The type can be nil, the
// zero value of a go slice, a dgo.ArrayType, or a dgo string that parses to a dgo.ArrayType.
func ArrayWithCapacity(capacity int, typ interface{}) dgo.Array {
	mt := asArrayType(typ)
	return &array{slice: make([]dgo.Value, 0, capacity), typ: mt, frozen: false}
}

// WrapSlice wraps the given slice in an array. Unset entries in the slice will be replaced by Nil.
func WrapSlice(values []dgo.Value) dgo.Array {
	ReplaceNil(values)
	return &array{slice: values, frozen: false}
}

// MutableValues returns a frozen dgo.Array that represents the given values
func MutableValues(values []interface{}) dgo.Array {
	cp := make([]dgo.Value, len(values))
	for i := range values {
		cp[i] = Value(values[i])
	}
	return &array{slice: cp, frozen: false}
}

func newArray(t dgo.ArrayType, arg dgo.Value) dgo.Array {
	if args, ok := arg.(dgo.Arguments); ok {
		args.AssertSize(`array`, 1, 1)
		arg = args.Get(0)
	}
	a := Array(arg)
	if !t.Assignable(a) {
		panic(IllegalAssignment(t, a))
	}
	return a
}

func valueSlice(values []interface{}, frozen bool) []dgo.Value {
	cp := make([]dgo.Value, len(values))
	if frozen {
		for i := range values {
			v := Value(values[i])
			if f, ok := v.(dgo.Freezable); ok {
				v = f.FrozenCopy()
			}
			cp[i] = v
		}
	} else {
		for i := range values {
			cp[i] = Value(values[i])
		}
	}
	return cp
}

// Integers returns a dgo.Array that represents the given ints
func Integers(values []int) dgo.Array {
	cp := make([]dgo.Value, len(values))
	for i := range values {
		cp[i] = intVal(values[i])
	}
	return &array{slice: cp, frozen: true}
}

// Strings returns a dgo.Array that represents the given strings
func Strings(values []string) dgo.Array {
	cp := make([]dgo.Value, len(values))
	for i := range values {
		cp[i] = makeHString(values[i])
	}
	return &array{slice: cp, frozen: true}
}

// Values returns a frozen dgo.Array that represents the given values
func Values(values []interface{}) dgo.Array {
	return &array{slice: valueSlice(values, true), frozen: true}
}

func (v *array) assertType(e dgo.Value, pos int) {
	if t := v.typ; t != nil {
		sz := len(v.slice)
		if pos >= sz {
			sz++
			if sz > t.Max() {
				panic(IllegalSize(t, sz))
			}
		}
		var et dgo.Value
		if tp, ok := t.(dgo.TupleType); ok {
			if tp.Variadic() {
				lp := tp.Len() - 1
				if pos < lp {
					et = tp.Element(pos)
				} else {
					et = tp.Element(lp)
				}
			} else {
				et = tp.Element(pos)
			}
		} else {
			et = t.ElementType()
		}
		if !et.Assignable(e) {
			panic(IllegalAssignment(et, e))
		}
	}
}

func (v *array) assertTypes(values dgo.Iterable) {
	if t := v.typ; t != nil {
		addedSize := values.Len()
		if addedSize == 0 {
			return
		}
		sz := len(v.slice)
		if sz+addedSize > t.Max() {
			panic(IllegalSize(t, sz+addedSize))
		}
		et := t.ElementType()
		values.Each(func(e dgo.Value) {
			if !et.Assignable(e) {
				panic(IllegalAssignment(et, e))
			}
		})
	}
}

func (v *array) Add(vi interface{}) {
	if v.frozen {
		panic(frozenArray(`Add`))
	}
	val := Value(vi)
	v.assertType(val, len(v.slice))
	v.slice = append(v.slice, val)
}

func (v *array) AddAll(values dgo.Iterable) {
	if v.frozen {
		panic(frozenArray(`AddAll`))
	}
	v.assertTypes(values)
	a := v.slice
	if ar, ok := values.(*array); ok {
		a = ar.AppendToSlice(a)
	} else {
		values.Each(func(e dgo.Value) { a = append(a, e) })
	}
	v.slice = a
}

func (v *array) AddValues(values ...interface{}) {
	if v.frozen {
		panic(frozenArray(`AddValues`))
	}
	va := valueSlice(values, false)
	v.assertTypes(&array{slice: va})
	v.slice = append(v.slice, va...)
}

func (v *array) All(predicate dgo.Predicate) bool {
	a := v.slice
	for i := range a {
		if !predicate(a[i]) {
			return false
		}
	}
	return true
}

func (v *array) Any(predicate dgo.Predicate) bool {
	a := v.slice
	for i := range a {
		if predicate(a[i]) {
			return true
		}
	}
	return false
}

func (v *array) AppendTo(w dgo.Indenter) {
	w.AppendRune('{')
	ew := w.Indent()
	a := v.slice
	for i := range a {
		if i > 0 {
			w.AppendRune(',')
		}
		ew.NewLine()
		ew.AppendValue(a[i])
	}
	w.NewLine()
	w.AppendRune('}')
}

func (v *array) AppendToSlice(slice []dgo.Value) []dgo.Value {
	return append(slice, v.slice...)
}

func (v *array) Assignable(other interface{}) bool {
	return Assignable(nil, v, Value(other))
}

func (v *array) DeepAssignable(guard dgo.RecursionGuard, other dgo.Value) bool {
	es := v.slice
	switch ot := other.(type) {
	case defaultArrayType:
		return false // lacks size
	case *sizedArrayType:
		l := len(es)
		return ot.min == l && ot.max == l && assignableToAll(guard, ot.elementType, es)
	case *array:
		return sliceEquals(nil, es, ot.slice)
	case dgo.TupleType:
		return tupleAssignableTuple(guard, v, ot)
	}
	return CheckAssignableTo(guard, other, v)
}

func (v *array) CompareTo(other interface{}) (int, bool) {
	return compare(nil, v, Value(other))
}

func (v *array) deepCompare(seen []dgo.Value, other deepCompare) (int, bool) {
	ov, ok := other.(*array)
	if !ok {
		return 0, false
	}
	a := v.slice
	b := ov.slice
	top := len(a)
	max := len(b)
	r := 0
	if top < max {
		r = -1
		max = top
	} else if top > max {
		r = 1
	}

	for i := 0; i < max; i++ {
		if _, ok = a[i].(dgo.Comparable); !ok {
			r = 0
			break
		}
		var c int
		if c, ok = compare(seen, a[i], b[i]); !ok {
			r = 0
			break
		}
		if c != 0 {
			r = c
			break
		}
	}
	return r, ok
}

func (v *array) Copy(frozen bool) dgo.Array {
	if frozen && v.frozen {
		return v
	}
	cp := util.SliceCopy(v.slice)
	if frozen {
		for i := range cp {
			if f, ok := cp[i].(dgo.Freezable); ok {
				cp[i] = f.FrozenCopy()
			}
		}
	}
	return &array{slice: cp, typ: v.typ, frozen: frozen}
}

func (v *array) ContainsAll(other dgo.Iterable) bool {
	a := v.slice
	l := len(a)
	if l < other.Len() {
		return false
	}
	if l == 0 {
		return true
	}

	var vs []dgo.Value
	if oa, ok := other.(*array); ok {
		vs = util.SliceCopy(oa.slice)
	} else {
		vs = sliceFromIterable(other)
	}

	// Keep track of elements that have been found equal using v copy
	// where such elements are set to nil. This avoids excessive calls
	// to Equals
	for i := range vs {
		ea := a[i]
		f := false
		for j := range vs {
			if be := vs[j]; be != nil {
				if be.Equals(ea) {
					vs[j] = nil
					f = true
					break
				}
			}
		}
		if !f {
			return false
		}
	}
	return true
}

func (v *array) Each(actor dgo.Consumer) {
	a := v.slice
	for i := range a {
		actor(a[i])
	}
}

func (v *array) EachWithIndex(actor dgo.DoWithIndex) {
	a := v.slice
	for i := range a {
		actor(a[i], i)
	}
}

func (v *array) Element(index int) dgo.Value {
	return v.Get(index)
}

func (v *array) ElementType() dgo.Value {
	switch len(v.slice) {
	case 0:
		return DefaultAnyType
	case 1:
		return v.slice[0]
	}
	return (*allOfType)(v)
}

func (v *array) ElementTypes() dgo.Array {
	return v
}

func (v *array) Equals(other interface{}) bool {
	return equals(nil, v, other)
}

func (v *array) deepEqual(seen []dgo.Value, other deepEqual) bool {
	if ov, ok := other.(*array); ok {
		return sliceEquals(seen, v.slice, ov.slice)
	}
	return false
}

func (v *array) Find(finder dgo.Mapper) interface{} {
	a := v.slice
	for i := range a {
		if fv := finder(a[i]); fv != nil {
			return fv
		}
	}
	return nil
}

func (v *array) Flatten() dgo.Array {
	a := v.slice
	for i := range a {
		if _, ok := a[i].(*array); ok {
			fs := make([]dgo.Value, i, len(a)*2)
			copy(fs, a)
			return &array{slice: flattenElements(a[i:], fs), frozen: v.frozen}
		}
	}
	return v
}

func flattenElements(elements, receiver []dgo.Value) []dgo.Value {
	for i := range elements {
		e := elements[i]
		if a, ok := e.(*array); ok {
			receiver = flattenElements(a.slice, receiver)
		} else {
			receiver = append(receiver, e)
		}
	}
	return receiver
}

func (v *array) Freeze() {
	if v.frozen {
		return
	}
	v.frozen = true
	a := v.slice
	for i := range a {
		if f, ok := a[i].(dgo.Freezable); ok {
			f.Freeze()
		}
	}
}

func (v *array) Frozen() bool {
	return v.frozen
}

func (v *array) FrozenCopy() dgo.Value {
	return v.Copy(true)
}

func (v *array) Generic() dgo.Value {
	return &sizedArrayType{elementType: Generic(v.ElementType()), min: 0, max: math.MaxInt64}
}

func (v *array) GoSlice() []dgo.Value {
	if v.frozen {
		return util.SliceCopy(v.slice)
	}
	return v.slice
}

func (v *array) HashCode() int {
	return v.deepHashCode(nil)
}

func (v *array) deepHashCode(seen []dgo.Value) int {
	h := 1
	s := v.slice
	for i := range s {
		h = h*31 + deepHashCode(seen, s[i])
	}
	return h
}

func (v *array) Get(index int) dgo.Value {
	return v.slice[index]
}

func (v *array) IndexOf(vi interface{}) int {
	val := Value(vi)
	a := v.slice
	for i := range a {
		if val.Equals(a[i]) {
			return i
		}
	}
	return -1
}

func (v *array) Insert(pos int, vi interface{}) {
	if v.frozen {
		panic(frozenArray(`Insert`))
	}
	val := Value(vi)
	v.assertType(val, pos)
	v.slice = append(v.slice[:pos], append([]dgo.Value{val}, v.slice[pos:]...)...)
}

// InterfaceSlice returns the values held by the Array as a slice. The slice will
// contain dgo.Value instances.
func (v *array) InterfaceSlice() []interface{} {
	s := v.slice
	is := make([]interface{}, len(s))
	for i := range s {
		is[i] = s[i]
	}
	return is
}

func (v *array) Len() int {
	return len(v.slice)
}

func (v *array) MapTo(t dgo.ArrayType, mapper dgo.Mapper) dgo.Array {
	if t == nil {
		return v.Map(mapper)
	}
	a := v.slice
	l := len(a)
	if l < t.Min() {
		panic(IllegalSize(t, l))
	}
	if l > t.Max() {
		panic(IllegalSize(t, l))
	}
	et := t.ElementType()
	vs := make([]dgo.Value, len(a))

	for i := range a {
		mv := Value(mapper(a[i]))
		if !et.Assignable(mv) {
			panic(IllegalAssignment(et, mv))
		}
		vs[i] = mv
	}
	return &array{slice: vs, typ: t, frozen: v.frozen}
}

func (v *array) Map(mapper dgo.Mapper) dgo.Array {
	a := v.slice
	vs := make([]dgo.Value, len(a))
	for i := range a {
		vs[i] = Value(mapper(a[i]))
	}
	return &array{slice: vs, frozen: v.frozen}
}

func (v *array) One(predicate dgo.Predicate) bool {
	a := v.slice
	f := false
	for i := range a {
		if predicate(a[i]) {
			if f {
				return false
			}
			f = true
		}
	}
	return f
}

func (v *array) Reduce(mi interface{}, reductor func(memo dgo.Value, elem dgo.Value) interface{}) dgo.Value {
	memo := Value(mi)
	a := v.slice
	for i := range a {
		memo = Value(reductor(memo, a[i]))
	}
	return memo
}

func (v *array) ReflectTo(value reflect.Value) {
	vt := value.Type()
	ptr := vt.Kind() == reflect.Ptr
	if ptr {
		vt = vt.Elem()
	}
	if vt.Kind() == reflect.Interface && vt.Name() == `` {
		vt = v.ReflectType()
	}
	a := v.slice
	var s reflect.Value
	if !v.frozen && vt.Elem() == reflectValueType {
		s = reflect.ValueOf(a)
	} else {
		l := len(a)
		s = reflect.MakeSlice(vt, l, l)
		for i := range a {
			ReflectTo(a[i], s.Index(i))
		}
	}
	if ptr {
		// The created slice cannot be addressed. A pointer to it is necessary
		x := reflect.New(s.Type())
		x.Elem().Set(s)
		s = x
	}
	value.Set(s)
}

func (v *array) Max() int {
	return len(v.slice)
}

func (v *array) Min() int {
	return len(v.slice)
}

func (v *array) New(arg dgo.Value) dgo.Value {
	return newArray(v, arg)
}

func (v *array) ReflectType() reflect.Type {
	return reflect.SliceOf(v.ElementType().ReflectType())
}

func (v *array) TypeIdentifier() dgo.TypeIdentifier {
	return dgo.TiArrayExact
}

func (v *array) Unbounded() bool {
	return false
}

func (v *array) removePos(pos int) dgo.Value {
	a := v.slice
	if pos >= 0 && pos < len(a) {
		newLen := len(a) - 1
		if v.typ != nil {
			if v.typ.Min() > newLen {
				panic(IllegalSize(v.typ, newLen))
			}
		}
		val := a[pos]
		copy(a[pos:], a[pos+1:])
		a[newLen] = nil // release to GC
		v.slice = a[:newLen]
		return val
	}
	return nil
}

func (v *array) Remove(pos int) dgo.Value {
	if v.frozen {
		panic(frozenArray(`Remove`))
	}
	return v.removePos(pos)
}

func (v *array) RemoveValue(value interface{}) bool {
	if v.frozen {
		panic(frozenArray(`RemoveValue`))
	}
	return v.removePos(v.IndexOf(value)) != nil
}

func (v *array) Resolve(ap dgo.AliasProvider) {
	a := v.slice
	v.slice = nil
	for i := range a {
		a[i] = ap.Replace(a[i])
	}
	v.slice = a
}

func (v *array) Reject(predicate dgo.Predicate) dgo.Array {
	vs := make([]dgo.Value, 0)
	a := v.slice
	for i := range a {
		e := a[i]
		if !predicate(e) {
			vs = append(vs, e)
		}
	}
	return &array{slice: vs, typ: v.typ, frozen: v.frozen}
}

func (v *array) SameValues(other dgo.Iterable) bool {
	return len(v.slice) == other.Len() && v.ContainsAll(other)
}

func (v *array) Select(predicate dgo.Predicate) dgo.Array {
	vs := make([]dgo.Value, 0)
	a := v.slice
	for i := range a {
		e := a[i]
		if predicate(e) {
			vs = append(vs, e)
		}
	}
	return &array{slice: vs, typ: v.typ, frozen: v.frozen}
}

func (v *array) Set(pos int, vi interface{}) dgo.Value {
	if v.frozen {
		panic(frozenArray(`Set`))
	}
	val := Value(vi)
	v.assertType(val, pos)
	old := v.slice[pos]
	v.slice[pos] = val
	return old
}

func (v *array) SetType(ti interface{}) {
	if v.frozen {
		panic(frozenArray(`SetType`))
	}
	mt := asArrayType(ti)
	if mt == nil || mt.Assignable(v) {
		v.typ = mt
		return
	}
	panic(IllegalAssignment(mt, v))
}

func (v *array) Slice(i, j int) dgo.Array {
	if v.frozen && i == 0 && j == len(v.slice) {
		return v
	}
	ss := v.slice[i:j]
	if !v.frozen {
		// v copy is needed. Two non frozen arrays cannot share the same slice storage
		ss = util.SliceCopy(ss)
	}
	return &array{slice: ss, frozen: v.frozen}
}

func (v *array) Sort() dgo.Array {
	sa := v.slice
	if len(sa) < 2 {
		return v
	}
	sorted := util.SliceCopy(sa)
	sort.SliceStable(sorted, func(i, j int) bool {
		a := sorted[i]
		b := sorted[j]
		if ac, ok := a.(dgo.Comparable); ok {
			var c int
			if c, ok = ac.CompareTo(b); ok {
				return c < 0
			}
		}
		return a.TypeIdentifier() < b.TypeIdentifier()
	})
	return &array{slice: sorted, typ: v.typ, frozen: v.frozen}
}

func (v *array) String() string {
	return util.ToStringERP(v)
}

func (v *array) ToMap() dgo.Map {
	ms := v.slice
	top := len(ms)

	ts := top / 2
	if top%2 != 0 {
		ts++
	}
	tbl := make([]*hashNode, tableSizeFor(ts))
	hl := len(tbl) - 1
	m := &hashMap{table: tbl, len: ts, frozen: v.frozen}

	for i := 0; i < top; {
		mk := ms[i]
		i++
		var mv dgo.Value = Nil
		if i < top {
			mv = ms[i]
			i++
		}
		hk := hl & hash(mk.HashCode())
		nd := &hashNode{mapEntry: mapEntry{key: mk, value: mv}, hashNext: tbl[hk], prev: m.last}
		if m.first == nil {
			m.first = nd
		} else {
			m.last.next = nd
		}
		m.last = nd
		tbl[hk] = nd
	}
	return m
}

func (v *array) ToMapFromEntries() (dgo.Map, bool) {
	ms := v.slice
	top := len(ms)
	tbl := make([]*hashNode, tableSizeFor(top))
	hl := len(tbl) - 1
	m := &hashMap{table: tbl, len: top, frozen: v.frozen}

	for i := range ms {
		nd, ok := ms[i].(*hashNode)
		if !ok {
			var ea *array
			if ea, ok = ms[i].(*array); ok && len(ea.slice) == 2 {
				nd = &hashNode{mapEntry: mapEntry{key: ea.slice[0], value: ea.slice[1]}}
			} else {
				return nil, false
			}
		} else if nd.hashNext != nil {
			// Copy node, it belongs to another map
			c := *nd
			c.next = nil // this one might not get assigned below
			nd = &c
		}

		hk := hl & hash(nd.key.HashCode())
		nd.hashNext = tbl[hk]
		nd.prev = m.last
		if m.first == nil {
			m.first = nd
		} else {
			m.last.next = nd
		}
		m.last = nd
		tbl[hk] = nd
	}
	return m, true
}

func (v *array) Unique() dgo.Array {
	a := v.slice
	top := len(a)
	if top < 2 {
		return v
	}
	tbl := make([]*hashNode, tableSizeFor(int(float64(top)/loadFactor)))
	hl := len(tbl) - 1
	u := make([]dgo.Value, top)
	ui := 0

nextVal:
	for i := range a {
		k := a[i]
		hk := hl & hash(k.HashCode())
		for e := tbl[hk]; e != nil; e = e.hashNext {
			if k.Equals(e.key) {
				continue nextVal
			}
		}
		tbl[hk] = &hashNode{mapEntry: mapEntry{key: k}, hashNext: tbl[hk]}
		u[ui] = k
		ui++
	}
	if ui == top {
		return v
	}
	return &array{slice: u[:ui], typ: v.typ, frozen: v.frozen}
}

func (v *array) Pop() (dgo.Value, bool) {
	if v.frozen {
		panic(frozenArray(`Pop`))
	}
	p := len(v.slice) - 1
	if p >= 0 {
		return v.removePos(p), true
	}
	return nil, false
}

func (v *array) Variadic() bool {
	return false
}

func (v *array) With(vi interface{}) dgo.Array {
	val := Value(vi)
	v.assertType(val, len(v.slice))
	return &array{slice: append(v.slice, val), typ: v.typ, frozen: v.frozen}
}

func (v *array) WithAll(values dgo.Iterable) dgo.Array {
	if values.Len() == 0 {
		return v
	}
	c := v.Copy(false)
	if v.frozen {
		values = values.FrozenCopy().(dgo.Iterable)
	}
	c.AddAll(values)
	c.(*array).frozen = v.frozen
	return c
}

func (v *array) WithValues(values ...interface{}) dgo.Array {
	if len(values) == 0 {
		return v
	}
	va := valueSlice(values, v.frozen)
	v.assertTypes(&array{slice: va})
	return &array{slice: append(v.slice, va...), typ: v.typ, frozen: v.frozen}
}

// ReplaceNil performs an in-place replacement of nil interfaces with the NilValue
func ReplaceNil(vs []dgo.Value) {
	for i := range vs {
		if vs[i] == nil {
			vs[i] = Nil
		}
	}
}

// allAssignable returns true when all elements of slice vs are assignable to the given type t
func allAssignable(guard dgo.RecursionGuard, t dgo.Value, vs []dgo.Value) bool {
	if t == DefaultAnyType {
		return true
	}
	for i := range vs {
		if !Assignable(guard, t, vs[i]) {
			return false
		}
	}
	return true
}

// assignableToAll returns true when the given type t is assignable the type of all elements of slice vs
func assignableToAll(guard dgo.RecursionGuard, t dgo.Value, vs []dgo.Value) bool {
	for i := range vs {
		if !Assignable(guard, vs[i], t) {
			return false
		}
	}
	return true
}

func frozenArray(f string) error {
	return fmt.Errorf(`%s called on a frozen Array`, f)
}

func resolveSlice(ts []dgo.Value, ap dgo.AliasProvider) {
	for i := range ts {
		ts[i] = ap.Replace(ts[i])
	}
}
