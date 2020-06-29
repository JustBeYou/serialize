package encoders

import "encoding/binary"

// This piece of code was generated!!!
// DO NOT EDIT

func IntAsBytes(data int) []byte {
	a := Int64AsBytes(int64(data))
	return a
}

func UintAsBytes(data uint) []byte {
	a := Uint64AsBytes(uint64(data))
	return a
}

func Int8AsBytes(data int8) []byte {
	var output []byte = make([]byte, 1)
	output[0] = byte(data)
	return output
}

func Uint8AsBytes(data uint8) []byte {
	var output []byte = make([]byte, 1)
	output[0] = byte(data)
	return output
}


func Uint16AsBytes(data uint16) []byte {
	var output []byte = make([]byte, 2)
	binary.LittleEndian.PutUint16(output, data)
	return output
}

func Uint32AsBytes(data uint32) []byte {
	var output []byte = make([]byte, 4)
	binary.LittleEndian.PutUint32(output, data)
	return output
}

func Uint64AsBytes(data uint64) []byte {
	var output []byte = make([]byte, 8)
	binary.LittleEndian.PutUint64(output, data)
	return output
}

func Int16AsBytes(data int16) []byte {
	var output []byte = make([]byte, 2)
	binary.LittleEndian.PutUint16(output, uint16(data))
	return output
}

func Int32AsBytes(data int32) []byte {
	var output []byte = make([]byte, 4)
	binary.LittleEndian.PutUint32(output, uint32(data))
	return output
}

func Int64AsBytes(data int64) []byte {
	var output []byte = make([]byte, 8)
	binary.LittleEndian.PutUint64(output, uint64(data))
	return output
}
