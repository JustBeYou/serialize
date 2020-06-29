package decoders

import "encoding/binary"

// This piece of code was generated!!!
// DO NOT EDIT

func IntFromBytes(data []byte) (int, uint64) {
	a, b := Int64FromBytes(data)
	return int(a), b
}

func UintFromBytes(data []byte) (uint, uint64) {
	a, b := Uint64FromBytes(data)
	return uint(a), b
}

func Int8FromBytes(data []byte) (int8, uint64) {
	return int8(data[0]), 1
}

func Uint8FromBytes(data []byte) (uint8, uint64) {
	return uint8(data[0]), 1
}


func Uint16FromBytes(data []byte) (uint16, uint64) {
	return binary.LittleEndian.Uint16(data), 2
}

func Uint32FromBytes(data []byte) (uint32, uint64) {
	return binary.LittleEndian.Uint32(data), 4
}

func Uint64FromBytes(data []byte) (uint64, uint64) {
	return binary.LittleEndian.Uint64(data), 8
}

func Int16FromBytes(data []byte) (int16, uint64) {
	return int16(binary.LittleEndian.Uint16(data)), 2
}

func Int32FromBytes(data []byte) (int32, uint64) {
	return int32(binary.LittleEndian.Uint32(data)), 4
}

func Int64FromBytes(data []byte) (int64, uint64) {
	return int64(binary.LittleEndian.Uint64(data)), 8
}