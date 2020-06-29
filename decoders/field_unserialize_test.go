package decoders

import "testing"

func TestBoolFromBytes(t *testing.T) {
	value, consumed, _ := BoolFromBytes([]byte{1})
	if value != true || consumed != 1 {
		t.Error("bad decoding")
	}
}

func TestStringFromBytes(t *testing.T) {
	expected := "AAAA"
	input := []byte{0x0, 0x4, 0x41, 0x41, 0x41, 0x41}
	value, consumed, err := StringFromBytes(input)
	if value != expected || consumed != uint64(len(input)) || err != nil{
		t.Errorf("invalid dencoding (%d): %s != %s, %v", consumed, value, expected, err)
	}
}
