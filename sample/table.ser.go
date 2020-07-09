package sample

var TypeIdTable = map[string]uint16{
	"Foo": 1,
	"Bar": 2,
}

var IdTypeTable map[uint16]func([]byte) (interface{}, uint64, error)
func init() {
	IdTypeTable = map[uint16]func([]byte) (interface{}, uint64, error) {
		1: func (data []byte) (interface{}, uint64, error) {
			return Foo{}.Unserialize(data)
		},
		2: func (data []byte) (interface{}, uint64, error) {
			return Bar{}.Unserialize(data)
		},
	}
}