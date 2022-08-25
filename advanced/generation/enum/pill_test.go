package enum

import (
	"fmt"
	"testing"
)

func TestPrint(t *testing.T) {
	var p1 = Placebo
	var p2 = Aspirin
	var p3 = Ibuprofen
	var p4 = Paracetamol
	var p5 = Acetaminophen

	fmt.Println("p1=", p1)
	fmt.Println("p2=", p2)
	fmt.Println("p3=", p3)
	fmt.Println("p4=", p4)
	fmt.Println("p5=", p5)
}
