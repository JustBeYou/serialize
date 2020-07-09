package sample

//go:generate serialize -file=$GOFILE -type=Foo -serializer=hashing
type Foo struct {
	Bar int
	Fizz uint
	Buzz string
	FizzBuzz bool
	BarArray []bool
	Hash string `hashing:"ignore"`
	Custom Bar
	CustomInterface interface{}
}

//go:generate serialize -file=$GOFILE -type=Bar -serializer=hashing
type Bar struct {
	Value int
}

type OnlyBar interface {
	OnlyBar() error
}

func (b Bar) OnlyBar() error {
	return nil
}

type HashingSerializer interface {
	HashingSerialize() ([]byte, error)
}

type HashingUnserializer interface {
	HashingUnserialize([]byte) (interface{}, uint64, error)
}