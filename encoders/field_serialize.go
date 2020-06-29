package encoders

import (
	"encoding/binary"
	"serialize/standard"
)

func BoolAsBytes(isTrue bool) []byte {
	// isArray is used to store the boolean value
	header, _ := standard.FieldHeader{}.Serialize()
	if isTrue {
		header[0] |= 1
	}
	return header
}

func StringAsBytes(data string) []byte {
	header := standard.FieldHeader{}

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