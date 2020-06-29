package sample

import (
	"bytes"
	"testing"
)

func TestFoo(t *testing.T) {
	foo := Foo{
		-10,
		20,
		"foobar",
		true,
	}

	output, _ := foo.Serialize()

	expected := []byte{246, 255, 255, 255, 255, 255, 255, 255, 20, 0, 0, 0, 0, 0, 0, 0, 0, 6, 102, 111, 111, 98, 97, 114, 1}
	if !bytes.Equal(output, expected) {
		t.Errorf("Bad serialization: %v != %v", output, expected)
	}

	unserializedData, err := Foo{}.Unserialize(output)
	newFoo := unserializedData.(Foo)

	if err != nil {
		t.Errorf("Could not unserialize: %s\n", err.Error())
	}

	if newFoo.Bar != -10 || newFoo.Fizz != 20 || newFoo.Buzz != "foobar" || newFoo.FizzBuzz != true {
		t.Errorf("Bad deserialization: %v\n", newFoo)
	}
}