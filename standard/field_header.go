package standard

import "errors"

// TODO: support maps
// TODO: support array of interface{}

/**
 * Field encoding format
 * 1st byte (upper) - flags (left to right)
 *		- 1st bit - 16 bit array/string/map size?
 *		- 2nd bit - 32 bit array/string/map size?
 *		- 3rd bit - 64 bit array/string/map size?
 *		- 4th bit - 16 bit dynamic type identifier?
 *      - 5th .. 8th bit - RESERVED FOR FUTURE USE
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
	Is16BitIdentifier bool
}

func NewArrayHeader(size uint64) FieldHeader {
	header := FieldHeader{}
	if size > 0xffffffff {
		header.Is64BitSize = true
	} else if size > 0xffff {
		header.Is32BitSize = true
	} else if size > 0xff {
		header.Is16BitSize = true
	}

	return header
}

func NewIdentifierHeader(identifier uint16) FieldHeader {
	header := FieldHeader{}
	if identifier > 0xff {
		header.Is16BitIdentifier = true
	}

	return header
}

func FieldHeaderFromBytes(data byte) (FieldHeader, error) {
	header := FieldHeader{}
	header.Is16BitSize = (data>>7 & 0x1) == 0x1;
	header.Is32BitSize = (data>>6 & 0x1) == 0x1;
	header.Is64BitSize = (data>>5 & 0x1) == 0x1;
	header.Is16BitIdentifier = (data>>4 & 0x1) == 0x1;

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

	if h.Is16BitIdentifier {
		output |= 1 << 4
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

var DefaultTypes = map[string]bool{
	"bool": true,
	"string": true,
	"int": true,
	"uint": true,
	"int8": true,
	"int16": true,
	"int32": true,
	"int64": true,
	"uint8": true,
	"uint16": true,
	"uint32": true,
	"uint64": true,
}