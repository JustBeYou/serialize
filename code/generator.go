package code

import (
	"fmt"
	"strings"
)

type FieldOptions map[string]SerializerOptions
type SerializerOptions map[string]bool

type StructField struct {
	name string
	typeName string
	isArray bool
	options FieldOptions
	isCustomType bool
}

func GenPackageHeaderAndImports(name string) string {
	return fmt.Sprintf(`
		package %s
		import (
			"github.com/JustBeYou/serialize/encoders"
			"github.com/JustBeYou/serialize/decoders"
			"github.com/JustBeYou/serialize/standard"
			"errors"
			"fmt"
		)
	
		// THIS FILE WAS GENERATED BY serialize
		// PLEASE DO NOT EDIT
	`, name)
}

func GenSerializationHeader(serializerName, typeName string) string {
	return fmt.Sprintf(`
		func (self %s)%sSerialize() ([]byte, error) {
			var output, bytesTemp []byte
			var tempHeader standard.FieldHeader
			var tempLen uint64
	`, typeName, strings.Title(serializerName))
}

func GenSerializationFooter() string {
	return "return output, nil }"
}

func GenFieldSerialization(serializerName string, fieldInfo StructField) string {
	if serializerOptions, ok := fieldInfo.options[serializerName]; ok {
		if _, ok := serializerOptions["ignore"]; ok {
			return ""
		}
	}

	if fieldInfo.isArray {
		generated := fmt.Sprintf(`
			tempLen = uint64(len(self.%s))
			tempHeader = standard.NewArrayHeader(tempLen)
			bytesTemp, _ = tempHeader.Serialize()

			if tempHeader.Is16BitSize {
				bytesTemp = append(bytesTemp, encoders.Uint16AsBytes(uint16(tempLen))...)
			} else if tempHeader.Is32BitSize {
				bytesTemp = append(bytesTemp, encoders.Uint32AsBytes(uint32(tempLen))...)
			} else if tempHeader.Is64BitSize {
				bytesTemp = append(bytesTemp, encoders.Uint64AsBytes(uint64(tempLen))...)
			} else {
				bytesTemp = append(bytesTemp, encoders.Uint8AsBytes(uint8(tempLen))...)
			}`, fieldInfo.name)

		if fieldInfo.isCustomType {
			return "// CUSTOM []TYPE DETECTED!\n"
		} else {
			generated += fmt.Sprintf(`
			for _, v := range self.%s {
				bytesTemp = append(bytesTemp, encoders.%sAsBytes(v)...)
			}
			output = append(output, bytesTemp...)
			`, fieldInfo.name, strings.Title(fieldInfo.typeName))
		}

		return generated
	}

	if fieldInfo.isCustomType {
		return fmt.Sprintf(`
			bytesTemp, _ = self.%s.%sSerialize()
			output = append(output, bytesTemp...)
		`, fieldInfo.name, strings.Title(serializerName))
	}

	return fmt.Sprintf(`
		bytesTemp = encoders.%sAsBytes(self.%s)
		output = append(output, bytesTemp...)
	`, strings.Title(fieldInfo.typeName), fieldInfo.name)
}

func GenUnserializationHeader(serializerName, typeName string) string {
	return fmt.Sprintf(`
		func (self %s) %sUnserialize(data []byte) (interface{}, uint64, error) {
			var output %s
			var index uint64 = 0
			var consumed uint64 = 0
			var err error
			var tempHeader standard.FieldHeader
			var tempLen uint64
			var tempCustom interface{}
	`, typeName, strings.Title(serializerName), typeName)
}

func GenFieldUnserialization(serializerName string, fieldInfo StructField) string {
	if serializerOptions, ok := fieldInfo.options[serializerName]; ok {
		if _, ok := serializerOptions["ignore"]; ok {
			return ""
		}
	}

	if fieldInfo.isArray {
		generated := `
			tempHeader, err = standard.FieldHeaderFromBytes(data[index])
			if err != nil {
				return output, index, errors.New(fmt.Sprintf("Could not decode"))
			}
			index += 1

			if tempHeader.Is16BitSize {
				var tempLen2 uint16
				tempLen2, consumed, err = decoders.Uint16FromBytes(data[index:])
				tempLen = uint64(tempLen2)
			} else if tempHeader.Is32BitSize {
				var tempLen2 uint32
				tempLen2, consumed, err = decoders.Uint32FromBytes(data[index:])
				tempLen = uint64(tempLen2)
			} else if tempHeader.Is64BitSize {
				var tempLen2 uint64
				tempLen2, consumed, err = decoders.Uint64FromBytes(data[index:])
				tempLen = uint64(tempLen2)
			} else {
				var tempLen2 uint8
				tempLen2, consumed, err = decoders.Uint8FromBytes(data[index:])
				tempLen = uint64(tempLen2)
			}
			index += consumed
			`

		if fieldInfo.isCustomType {
			return "// CUSTOM []TYPE DETECTED!\n"
		} else {
			generated += fmt.Sprintf(`
			for i := uint64(0); i < tempLen; i++  {
				var tempValue %s
				tempValue, consumed, err = decoders.%sFromBytes(data[index:])
				output.%s = append(output.%s, tempValue)
				if err != nil {
					return output, index, errors.New(fmt.Sprintf("Could not decode at %%d: %%s\n", index, err.Error()))
				}
				index += consumed
			}
			`, fieldInfo.typeName, strings.Title(fieldInfo.typeName), fieldInfo.name, fieldInfo.name)
		}

		return generated
	}

	if fieldInfo.isCustomType {
		return fmt.Sprintf(`
			tempCustom, consumed, err = %s{}.%sUnserialize(data[index:])
			if err != nil {
				return output, index, errors.New(fmt.Sprintf("Could not decode"))
			}
			output.%s = tempCustom.(%s)
			index += consumed
		`, fieldInfo.typeName, strings.Title(serializerName), fieldInfo.name, fieldInfo.typeName)
	}

	return fmt.Sprintf(`
		output.%s, consumed, err = decoders.%sFromBytes(data[index:])
		if err != nil {
			return output, index, errors.New(fmt.Sprintf("Could not decode"))
		}
		index += consumed
	`, fieldInfo.name, strings.Title(fieldInfo.typeName))
}

func GenUnserializationFooter() string {
	return "return output, index, nil }"
}