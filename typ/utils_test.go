package typ

import (
	"fmt"

	"github.com/lyraproj/dgo/vf"
)

func ExampleGeneric() {
	vt := vf.Strings(`hello`, `world`).Type()
	fmt.Println(vt)
	fmt.Println(Generic(vt))

	// Output:
	// {"hello","world"}
	// []string
}

func ExampleAsType() {
	v := vf.String(`hello`)
	fmt.Println(AsType(v))

	// Output:
	// "hello"
}
