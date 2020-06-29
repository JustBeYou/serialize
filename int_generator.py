template = """
func {FuncName}FromBytes(data []byte) ({DataType}, uint64) {
	return binary.LittleEndian.{FuncName}(data), {Size}
}
"""

types = [
    ("uint16", 2),
    ("uint32", 4),
    ("uint64", 8),
]

code = """
// This piece of code was generated!!!
// DO NOT EDIT

func IntFromBytes(data []byte) (int, uint64) {
        a, b := Int64FromBytes(data)
        return a, b
}

func UintFromBytes(data []byte) (uint, uint64) {
        a, b := Uint64FromBytes(data)
        return a, b
}

func Int8FromBytes(data []byte) (int8, uint64) {
	return int8(data[0]), 1
}

func Uint8FromBytes(data []byte) (uint8, uint64) {
	return uint8(data[0]), 1
}

"""

for t in types:
    code += template.replace("{FuncName}", t[0].capitalize()).replace("{DataType}", t[0]).replace("{Size}", str(t[1]))

print (code)
