package sample

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