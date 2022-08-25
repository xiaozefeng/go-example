package container

import "fmt"

//go:generate ./gen.sh ./template/container.tmp container uint32 container
func generateUint32Example() {
	var u uint32 = 42
	c := NewUint32Container()
	c.Put(u)
	v := c.Get()
	fmt.Printf("generateExample: %d (%T)\n", v, v)
}

//go:generate ./gen.sh ./template/container.tmp container string container
func generateStringExample() {
	var s string = "Hello"
	c := NewStringContainer()
	c.Put(s)
	v := c.Get()
	fmt.Printf("generateExample: %s (%T)\n", v, v)
}
