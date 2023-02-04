package object

import "testing"

func TestStringHashing(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "Hello Alien Planet"}
	diff2 := &String{Value: "Hello Alien Planet"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Fatalf("Strings with the same value have different hash values")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Fatalf("Strings with the same value have different hash values")
	}

	if hello1.HashKey() == diff2.HashKey() {
		t.Fatalf("Strings with the different value have same hash values")
	}
}
