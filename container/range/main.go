package main

import "fmt"

func main() {
	// range array
	fmt.Println("range array")
	var a = [5]int{1, 3, 5, 7, 9}
	rangeArray(a)

	// range slice
	fmt.Println("range slice s1")
	var s1 = []int{1, 3, 5, 7, 9}
	rangeSlice(s1)
	fmt.Println("range slice s2")
	var s2 = make([]int, 6)
	rangeSlice(s2)
	var s3 = s1[1:]
	fmt.Println("range slice s3")
	rangeSlice(s3)

	// range map
	fmt.Println("range map: ")
	var m1 = map[string]int{
		"jackie": 99,
		"mickey": 98,
		"luna":   97,
	}
	rangeMap(m1)
}

func rangeMap(m map[string]int) {
	for k, v := range m {
		fmt.Printf("key: %s, value: %d \n", k, v)
	}
}

func rangeSlice(s1 []int) {
	for i, v := range s1 {
		fmt.Printf("index:%d ,value: %d \n ", i, v)
	}
}

func rangeArray(a [5]int) {
	for i, v := range a {
		fmt.Printf("index:%d ,value: %d \n ", i, v)
	}
}
