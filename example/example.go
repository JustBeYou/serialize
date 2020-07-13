package example

//go:generate serialize -file=$GOFILE -type=Foo -serializer=hashing
type Foo struct {
	Value int
	Hash string `hashing:"ignore"`
	Custom Bar
	Poly interface{}
}

//go:generate serialize -file=$GOFILE -type=Bar -serializer=hashing
type Bar struct {
	Value int
}

//go:generate serialize -file=$GOFILE -table=true
var TypeIdTable = map[string]uint16{
	"Foo": 1,
	"Bar": 2,
}

func example() {
	foo := Foo{0, "this will be the hash", Bar{0}, Bar{1},
	}

	// Default serializer
	output, _ := foo.Serialize()
	unserializedData, _, err := Foo{}.Unserialize(output)
	newFoo := unserializedData.(Foo)

	// Custom serializer that ignores the Hash field
	output, _ := foo.HashingSerialize()
	unserializedData, _, err := Foo{}.HashingUnserialize(output)
	newFoo := unserializedData.(Foo)
}

// Interfaces for our custom serializer
type HashingSerializer interface {
	HashingSerialize() ([]byte, error)
}

type HashingUnserializer interface {
	HashingUnserialize([]byte) (interface{}, uint64, error)
}