package main

//go:generate serialize $GOFILE -type=Foo
type Foo struct {
	Bar int
	Fizz uint
}