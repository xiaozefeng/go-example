package web

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestRegexp_Compile(t *testing.T) {
	var s = `:name(^.+$)`
	containsAny := strings.ContainsAny(s, `()`)
	fmt.Println(containsAny)
	//seg := strings.Split(s, `(`)
	//fmt.Println("seg[0]", seg[0][1:])
	//fmt.Println("seg[1]", `(`+seg[1])
	reg := regexp.MustCompile(s)
	fmt.Println("reg", reg)

}
