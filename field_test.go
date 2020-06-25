package main

import "testing"

func TestNewFieldHeader(t *testing.T) {
	_, err := NewFieldHeader(0xff)
	if err != nil || err.Error() != "invalid field ID" {
		t.Error("ID validation failed")
	}
	_, err = NewFieldHeader(0xe1)
	if err != nil || err.Error() != "multiple size flags set" {
		t.Error("size validation failed")
	}
	_, err = NewFieldHeader(0x41)
	if err != nil || err.Error() != "size flag set for something that is not an array" {
		t.Error("array flag validation failed")
	}

	h, err := NewFieldHeader(0xc1)
	if err != nil || h.IsArray == false || h.Is16BitSize == false || h.TypeID != UintID {
		t.Error("invalid header creation")
	}
}

func TestFieldHeader_Serialize(t *testing.T) {
	h := FieldHeader{
		UintID,
		true,
		false,
		false,
		false,
	}
	output, err := h.Serialize()
	if err != nil || output[0] != 0x81 {
		t.Errorf("invalid serialized output: %x", output[0])
	}
}
