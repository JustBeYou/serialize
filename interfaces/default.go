package interfaces

type Serializer interface {
	Serialize() ([]byte, error)
}

type Unserializer interface {
	Unserialize([]byte) (interface{}, error)
}
