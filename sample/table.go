package sample

// go:generate serialize -file=$GOFILE -table=true
var TypeIdTable = map[string]uint16{
	"Foo": 1,
	"Bar": 2,
}