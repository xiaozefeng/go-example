package memory

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestAlignof(t *testing.T) {
	fmt.Printf("string alignof is %d\n", unsafe.Alignof(string("a")))
	fmt.Printf("complex128 alignof is %d\n", unsafe.Alignof(complex128(0)))
	fmt.Printf("int32 alignof is %d\n", unsafe.Alignof(int32(0)))
	fmt.Printf("bool alignof is %d\n", unsafe.Alignof(false))
	fmt.Printf("int alignof is %d\n", unsafe.Alignof(int(0)))
	fmt.Printf("int64 alignof is %d\n", unsafe.Alignof(int64(0)))
	fmt.Printf("float64 alignof is %d\n", unsafe.Alignof(float64(0)))
	fmt.Printf("slice int alignof is %d\n", unsafe.Alignof(make([]int32, 0)))

}

type Args struct {
	b1 bool  //1
	n1 int32 // 4
	n3 int32 // 4
	n2 int   // 8
	b2 bool  //1
}

type User struct {
	A int32   // 4
	B []int32 // 24  go中slice就是占用 24个固定的
	C string  // 16
	D bool    // 1
}

func TestSizeOf(t *testing.T) {
	// 基础
	fmt.Println(unsafe.Sizeof(""))         // 16
	fmt.Println(unsafe.Sizeof(int(0)))     // 8
	fmt.Println(unsafe.Sizeof(int32(0)))   // 4
	fmt.Println(unsafe.Sizeof(false))      // 1
	fmt.Println(unsafe.Sizeof([]string{})) // 24
	fmt.Println(unsafe.Sizeof([]int32{}))  // 24
	var a interface{}
	fmt.Println(unsafe.Sizeof(a))                   // 16
	fmt.Println(unsafe.Sizeof(map[string]string{})) //8

	fmt.Println()
	fmt.Println(unsafe.Sizeof(Args{}))
	fmt.Println(unsafe.Sizeof(User{}))
}

//string alignof is 8
//complex128 alignof is 8
//int32 alignof is 4
//bool alignof is 1
//int alignof is 8
//int64 alignof is 8
//float64 alignof is 8
//slice int alignof is 8
