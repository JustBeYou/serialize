Serialize
===
`serialize` should be used to easily define serializable formats for your data structures.

By default, the serializer can be used by adding a `go:generate` directive:
```go
package main

//go:generate serialize -type=Foo
type Foo struct {
    Bar int
    Fizz uint
} 
``` 
In this example, for structure `Foo` will be generated two methods `Serialize()` and `Unserialize()`.

In some cases, you would like to have multiple ways of serializing your data structures, for example when you have a structure
containing its own hash, you would like to ignore this field when serializing for hash calculation.
For that, you should create a custom serializer, and a tag on the target field.
```go
package main

//go:generate serialize -type=Foo -serializer=Hashing
type Foo struct {
    Bar int
    Fizz uint
    Hash string `serialize:"Hashing" Hashing:"ignore"`
} 
```
For the code above, two interfaces will be generated, with the following naming convetion:
```
// generic
interface <serializer>Serializer - method <serializer>Serialize()
interface <serializer>Unserializer - method <serializer>Unserialize()

// for our example
interface HashingSerializer - method HashingSerialize()
interface HashingUnserializer - method HashingUnserialize()
```