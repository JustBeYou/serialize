Serialize
===
`serialize` should be used to easily define serializable structs for your already existing data structures, 
everything in native Go.

**Table of contents**
1. [Motivation](#motivation)
2. [Features](#features)
3. [Example](#example)
4. [To do's](#todos)
5. [Contributing & feedback](#contributing)

<a name="motivation"></a>
### Motivation
1. Using Protobuf, FlatBuffers, MessagePack or other related tools is hard to set up and 
implies using a separate domain specific language (DSL)
2. Using JSON or XML is not storage efficient
3. Most of the existing solutions do not support Go polymorphism (using 
`interface{}`)
4. There is little flexibility when it comes to treating different fields
in a different manner (for example, when serializing for hashing you would
like to exclude the hash field itself)
5. The standard way of serializing structures in Go is `gob`, but it is 
unintuitive to use and not suitable for processes like hashing (the output
is very environment dependent)

<a name="features"></a>
### Features
1. Serialize primitive types, user defined structures and arrays
2. Serialize polymorphic fields (ie. `interface{}`)
3. Options for each struct field
4. Multiple *serializers* interfaces with different options
5. Runtime information kept only about `interface{}` fields, structure schema is kept directly in code

<a name="example"></a>
### Example
This is a short example of `serialize`'s main features.

```go
package main

import "fmt"

//go:generate serialize -file=$GOFILE -type=Foo -serializer=hashing
type Foo struct {
	Value int
	Hash string `hashing:"ignore"` // This means that the "HashingSerializer" will ignore this field
	Custom Bar
	Poly interface{}
}

//go:generate serialize -file=$GOFILE -type=Bar -serializer=hashing
type Bar struct {
	Value int
}

// This table declaration is needed only if you use interface{} fields in your structs
//go:generate serialize -file=$GOFILE -table=true
var TypeIdTable = map[string]uint16{
	"Foo": 1,
	"Bar": 2,
}

func main() {
	foo := Foo{0, "this will be the hash", Bar{0}, Bar{1}}

	// Default serializer
	output, _ := foo.Serialize()
	unserializedData, consumedBytes, _ := Foo{}.Unserialize(output)
	newFoo := unserializedData.(Foo)
	fmt.Printf("%v -> %v -> %v (%d bytes)\n", foo, output, newFoo, consumedBytes)

	// Custom serializer that ignores the Hash field
	output, _ = foo.HashingSerialize()
	unserializedData, consumedBytes, _ = Foo{}.HashingUnserialize(output)
	newFoo = unserializedData.(Foo)
	fmt.Printf("%v -> %v -> %v (%d bytes)\n", foo, output, newFoo, consumedBytes)
}

// Interfaces for our custom serializer
type HashingSerializer interface {
	HashingSerialize() ([]byte, error)
}

type HashingUnserializer interface {
	HashingUnserialize([]byte) (interface{}, uint64, error)
}
```

<a name="todos"></a>
### To do's
1. Support `map`
2. Support arrays of `interface{}`
3. Benchmarks for time/memory
4. Better flow for polymorphism handling

<a name="contributing"></a>
### Contributing & feedback
If you find this library useful and you would like to share some thoughts or help developing it, drop me an email at
mihailferaru2000@gmail.com 

And don't forget to STAR it! :star2:
