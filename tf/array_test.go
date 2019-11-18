package tf_test

import (
	"fmt"

	"github.com/lyraproj/dgo/tf"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
)

func ExampleArray() {
	tp := tf.Array()
	fmt.Println(tp.Equals(typ.Array))
	fmt.Println(tp.Assignable(vf.Values(`hello`, 42)))
	fmt.Println(tp.Assignable(42))
	// Output:
	// true
	// true
	// false
}

func ExampleArray_min() {
	tp := tf.Array(2)
	fmt.Println(tp.Assignable(vf.Values(`hello`, 42)))
	fmt.Println(tp.Assignable(vf.Values(`hello`)))
	// Output:
	// true
	// false
}

func ExampleArray_type() {
	tp := tf.Array(typ.String)
	fmt.Println(tp.Assignable(vf.Values(`hello`)))
	fmt.Println(tp.Assignable(vf.Values(42)))
	// Output:
	// true
	// false
}

func ExampleArray_min_max() {
	tp := tf.Array(1, 2)
	fmt.Println(tp.Assignable(vf.Values(`hello`, 42)))
	fmt.Println(tp.Assignable(vf.Values(`hello`, 42, `word`)))
	// Output:
	// true
	// false
}

func ExampleArray_type_min() {
	// Create a new array type with a minimum size of 2
	tp := tf.Array(typ.String, 2)
	fmt.Println(tp.Assignable(vf.Values(`hello`, `word`)))
	fmt.Println(tp.Assignable(vf.Values(`hello`)))
	// Output:
	// true
	// false
}

func ExampleArray_type_min_max() {
	tp := tf.Array(typ.String, 2, 3)
	fmt.Println(tp.Assignable(vf.Values(`hello`, `word`)))
	// Output: true
}
