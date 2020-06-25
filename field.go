package main

import "errors"

type StructField struct {
	name string
	typeName string
}

/**
 * Field encoding format
 * 1st byte (lower) - FieldTypeId
 * 1st byte (upper) - flags (left to right)
 *		- 1st bit - is array? / is true? if BooleanID
 *		- 2nd bit - 16 bit array size?
 *		- 3rd bit - 32 bit array size?
 *		- 4th bit - 64 bit array size?
 * If array, 1 to 8 bytes for size (depending on flags)
 * Value bytes
 */
type FieldHeader struct {
	TypeID FieldTypeId
	IsArray bool
	Is16BitSize bool
	Is32BitSize bool
	Is64BitSize bool
}

func NewFieldHeader(data byte) (FieldHeader, error) {
	header := FieldHeader{}
	header.TypeID = FieldTypeId(data & 0xf);
	header.IsArray = (data>>7 & 0x1) == 0x1;
	header.Is16BitSize = (data>>6 & 0x1) == 0x1;
	header.Is32BitSize = (data>>5 & 0x1) == 0x1;
	header.Is64BitSize = (data>>4 & 0x1) == 0x1;

	return header, FieldHeaderValidator(header)
}

func (h FieldHeader) Serialize() ([]byte, error) {
	var output = byte(h.TypeID)
	if h.IsArray {
		output |= 1 << 7
		if h.Is16BitSize {
			output |= 1 << 6
		} else if h.Is32BitSize {
			output |= 1 << 5
		} else if h.Is64BitSize {
			output |= 1 << 4
		}
	}

	return []byte{output}, FieldHeaderValidator(h)
}

func FieldHeaderValidator(header FieldHeader) error {
	var err error = nil
	if header.TypeID >= InvalidID {
		err = errors.New("invalid field ID")
	} else if !oneOrLess(
		header.Is16BitSize,
		header.Is32BitSize,
		header.Is64BitSize,
	) {
		err = errors.New("multiple size flags set")
	} else if header.IsArray == false && (header.Is16BitSize || header.Is32BitSize || header.Is64BitSize) {
		err = errors.New("size flag set for something that is not an array")
	}
	return err
}

func oneOrLess(a, b, c bool) bool {
	return (((a != b) != c) && !(a && b && c)) || (!a && !b && !c)
}

type FieldTypeId uint8;
const (
	BooleanID FieldTypeId = iota
	UintID
	IntID
	Uint8ID
	Uint16ID
	Uint32ID
	Uint64ID
	Int8ID
	Int16ID
	Int32ID
	Int64ID
	StringID
	InvalidID // invalid field id
)