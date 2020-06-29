package sample

//go:generate serialize -file=$GOFILE -type=Foo
type Foo struct {
	Bar int
	Fizz uint
	Buzz string
	FizzBuzz bool
}