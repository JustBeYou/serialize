package sample

import (
	"bytes"
	"testing"
)

// TODO: clean up tests code
func TestFoo(t *testing.T) {
	foo := Foo{
		-10,
		20,
		"foobar",
		true,
		[]bool{true, false, true},
		"AAAA",
		Bar{100},
		Bar{100},
	}

	output, _ := foo.Serialize()

	expected := []byte{246, 255, 255, 255, 255, 255, 255, 255, 20, 0, 0, 0, 0, 0, 0, 0, 0, 6, 102, 111, 111, 98, 97, 114, 1, 0, 3, 1, 0, 1, 0, 4, 65, 65, 65, 65, 100, 0, 0, 0, 0, 0, 0, 0, 0, 2, 100, 0, 0, 0, 0, 0, 0, 0}
	if !bytes.Equal(output, expected) {
		t.Errorf("Bad serialization: %v != %v", output, expected)
	}

	unserializedData, _, err := Foo{}.Unserialize(output)
	newFoo := unserializedData.(Foo)

	if err != nil {
		t.Errorf("Could not unserialize: %s\n", err.Error())
	}

	if newFoo.Bar != -10 || newFoo.Fizz != 20 || newFoo.Buzz != "foobar" || newFoo.FizzBuzz != true || len(newFoo.BarArray) != 3 || newFoo.BarArray[0] != true || newFoo.BarArray[1] != false || newFoo.BarArray[2] != true || newFoo.Custom.Value != 100 {
		t.Errorf("Bad deserialization: %v != %v\n", newFoo, foo)
	}

	_, ok := newFoo.CustomInterface.(OnlyBar)
	if !ok {
		t.Error("Polymorphism check failed")
	}
}