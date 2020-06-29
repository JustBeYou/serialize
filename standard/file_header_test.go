package standard

import "testing"

func TestNewFieldHeader(t *testing.T) {
	_, err := NewFieldHeader(0xe1)
	if err != nil && err.Error() != "multiple size flags set" {
		t.Error("size validation failed")
	}
}

func TestFieldHeader_Serialize(t *testing.T) {
	h := FieldHeader{
		true,
		false,
		false,
	}
	output, err := h.Serialize()
	if err != nil || output[0] != 0x80 {
		t.Errorf("invalid serialized output: %x", output[0])
	}
}
