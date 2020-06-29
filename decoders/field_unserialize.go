package decoders

import (
	"encoding/binary"
	"serialize/standard"
)

func BoolFromBytes(data []byte) (bool, uint64) {
	return data[0] & 1 == 1, 1
}

func StringFromBytes(data []byte) (string, uint64, error) {
	var consumed uint64 = 1
	header, err := standard.NewFieldHeader(data[0])
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