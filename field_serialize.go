package main

import (
	"encoding/binary"
	"errors"
)

/**
 * Field encoding format
 * 1st byte (upper) - flags (left to right)
 *		- 1st bit - 16 bit array/string/map size?
 *		- 2nd bit - 32 bit array/string/map size?
 *		- 3rd bit - 64 bit array/string/map size?
 *		- 4th ... 7th - RESERVED
 * 		- 8th bit - if boolean and true -> 1, else 0
 * 1 to 8 bytes for size (depending on flags)
 * Value bytes
 *
 * If it is not an array/string/map/bool, then the header is omitted
 */
type FieldHeader struct {
	Is16BitSize bool
	Is32BitSize bool
	Is64BitSize bool
}

func NewFieldHeader(data byte) (FieldHeader, error) {
	header := FieldHeader{}
	header.Is16BitSize = (data>>7 & 0x1) == 0x1;
	header.Is32BitSize = (data>>6 & 0x1) == 0x1;
	header.Is64BitSize = (data>>5 & 0x1) == 0x1;

	return header, FieldHeaderValidator(header)
}

func (h FieldHeader) Serialize() ([]byte, error) {
	var output byte = 0
	if h.Is16BitSize {
		output |= 1 << 7
	} else if h.Is32BitSize {
		output |= 1 << 6
	} else if h.Is64BitSize {
		output |= 1 << 5
	}

	return []byte{output}, FieldHeaderValidator(h)
}

func FieldHeaderValidator(header FieldHeader) error {
	var err error = nil
	if !oneOrLess(
		header.Is16BitSize,
		header.Is32BitSize,
		header.Is64BitSize,
	) {
		err = errors.New("multiple size flags set")
	}
	return err
}

func oneOrLess(a, b, c bool) bool {
	return (((a != b) != c) && !(a && b && c)) || (!a && !b && !c)
}

var serializeTemplates = map[string]string{
	"bool": "output = output.append(BoolAsBytes(self.{{ .name }}))",
	"string": "output = output.append(StringAsBytes(self.{{ .name }}))",
}

func BoolAsBytes(isTrue bool) []byte {
	// isArray is used to store the boolean value
	header, _ := FieldHeader{}.Serialize()
	if isTrue {
		header[0] |= 1
	}
	return header
}

func BoolFromBytes(data []byte) (bool, uint64) {
	return data[0] & 1 == 1, 1
}

func StringAsBytes(data string) []byte {
	header := FieldHeader{}

	dataAsBytes := []byte(data)
	var sizeAsBytes []byte
	if len(dataAsBytes) > 0xffffffff {
		header.Is64BitSize = true
		sizeAsBytes = make([]byte, 8)
		binary.LittleEndian.PutUint64(sizeAsBytes, uint64(len(dataAsBytes)))
	} else if len(dataAsBytes) > 0xffff {
		header.Is32BitSize = true
		sizeAsBytes = make([]byte, 4)
		binary.LittleEndian.PutUint32(sizeAsBytes, uint32(len(dataAsBytes)))
	} else if len(dataAsBytes) > 0xff {
		header.Is16BitSize = true
		sizeAsBytes = make([]byte, 2)
		binary.LittleEndian.PutUint16(sizeAsBytes, uint16(len(dataAsBytes)))
	} else {
		sizeAsBytes = make([]byte, 1)
		sizeAsBytes[0] = byte(len(dataAsBytes))
	}

	headerAsBytes, _ := header.Serialize()

	var output []byte
	output = append(output, headerAsBytes...)
	output = append(output, sizeAsBytes...)
	output = append(output, dataAsBytes...)
	return output
}

func StringFromBytes(data []byte) (string, uint64, error) {
	var consumed uint64 = 1
	header, err := NewFieldHeader(data[0])
	if err != nil {
		return "", 0, err
	}

	var size uint64
	if header.Is16BitSize {
		size = uint64(binary.LittleEndian.Uint16(data[1:3]))
		consumed += 2
	} else if header.Is32BitSize {
		size = uint64(binary.LittleEndian.Uint16(data[1:5]))
		consumed += 4
	} else if header.Is64BitSize {
		size = uint64(binary.LittleEndian.Uint16(data[1:9]))
		consumed += 8
	} else {
		consumed += 1
		size = uint64(data[1])
	}

	output := string(data[consumed:consumed+size])
	consumed += size
	return output, consumed, nil
}