package encoding

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/toddyco/bacnet2go/specs"
	"reflect"
)

const (
	flag16bits byte = 0xFE
	flag32bits byte = 0xFF
)

// Encoder is the struct used to turn baetyl-bacnet types to byte arrays. All
// public methods of encoder can set the internal error value. If such
// error is set, all encoding methods will be no-ops. This allows to
// defer error checking after several encoding operations
type Encoder struct {
	buf *bytes.Buffer
	err error
}

func NewEncoder() Encoder {
	e := Encoder{
		buf: new(bytes.Buffer),
		err: nil,
	}
	return e
}

func (e *Encoder) Error() error {
	return e.err
}

func (e *Encoder) Bytes() []byte {
	return e.buf.Bytes()
}

// ContextUnsigned write a (context)tag / value pair where the value
// type is an unsigned int
func (e *Encoder) ContextUnsigned(tagNumber byte, value uint32) {
	if e.err != nil {
		return
	}
	length := valueLength(value)
	t := tag{
		ID:      tagNumber,
		Context: true,
		Value:   uint32(length),
		Opening: false,
		Closing: false,
	}
	encodeTag(e.buf, t)
	unsigned(e.buf, value)
	// binary.Write(e.buf, binary.BigEndian, value)
}

func (e *Encoder) ContextSigned(tagNumber byte, value int32) {
	if e.err != nil {
		return
	}
	var length uint32
	if (value >= -128) && (value < 128) {
		length = 1
	} else if (value >= -32768) && (value < 32768) {
		length = 2
	} else if (value > -8388608) && (value < 8388608) {
		length = 3
	} else {
		length = 4
	}
	t := tag{
		ID:      tagNumber,
		Context: true,
		Value:   length,
		Opening: false,
		Closing: false,
	}

	encodeTag(e.buf, t)
	if (value >= -128) && (value < 128) {
		e.buf.WriteByte(uint8(value))
	} else if (value >= -32768) && (value < 32768) {
		binary.Write(e.buf, binary.BigEndian, uint16(value))
	} else if (value > -8388608) && (value < 8388608) {
		e.buf.WriteByte(byte(value >> 16))
		binary.Write(e.buf, binary.BigEndian, uint16(value))
	} else {
		binary.Write(e.buf, binary.BigEndian, value)
	}
	// binary.Write(e.buf, binary.BigEndian, value)
}

func (e *Encoder) ContextNull(tagNumber byte) {
	t := tag{
		ID:    tagNumber,
		Value: 1,
	}

	encodeTag(e.buf, t)
}

func (e *Encoder) ContextBoolean(tagNumber byte, value bool) {
	var i uint32
	if value {
		i = 1
	}
	t := tag{
		ID:    tagNumber,
		Value: i,
	}

	encodeTag(e.buf, t)
}

func (e *Encoder) ContextTypeReal(tagNumber byte, value float32) {
	t := tag{
		ID:    tagNumber,
		Value: 4,
	}

	encodeTag(e.buf, t)
	binary.Write(e.buf, binary.BigEndian, value)
}

func (e *Encoder) ContextTypeDouble(tagNumber byte, value float32) {
	t := tag{
		ID:    tagNumber,
		Value: 8,
	}

	encodeTag(e.buf, t)
	binary.Write(e.buf, binary.BigEndian, value)
}

// TODO:长度的逻辑有点问题,tag的context值
func (e *Encoder) ContextTypeOctetString(tagNumber byte, value string) {
	len := stringLength(value)
	t := tag{
		ID:      tagNumber,
		Context: true,
		Value:   1,
		Opening: false,
		Closing: false,
	}
	encodeTag(e.buf, t)

	if len != 0 {
		e.buf.WriteString(value)
	}
}

// TODO:tag的context值
func (e *Encoder) ContextTypeTypeCharacterString(tagNumber byte, value string) {
	len := stringLength(value)
	t := tag{
		ID:      tagNumber,
		Context: true,
		Value:   1,
		Opening: false,
		Closing: false,
	}
	encodeTag(e.buf, t)

	if len != 0 {
		e.buf.WriteString(value)
	}
}

// TODO:长度的逻辑有点问题,tag的context值
func (e *Encoder) ContextTypeTypeBitString(tagNumber byte, value string) {
	len := stringLength(value)
	t := tag{
		ID:      tagNumber,
		Context: true,
		Value:   1,
		Opening: false,
		Closing: false,
	}
	encodeTag(e.buf, t)

	if len != 0 {
		e.buf.WriteString(value)
	}
}

func stringLength(value string) int {
	if len(value) > 1476 {
		return 0
	} else {
		return len(value)
	}
}

func (e *Encoder) ContextTypeEnumerated(tagNumber byte, value uint32) {
	if e.err != nil {
		return
	}
	length := valueLength(value)
	t := tag{
		ID:      tagNumber,
		Context: false,
		Value:   uint32(length),
		Opening: false,
		Closing: false,
	}
	encodeTag(e.buf, t)
	unsigned(e.buf, value)
}

func (e *Encoder) ContextTypeDate(tagNumber byte, value uint32) {
	// TODO:待完善
}

func (e *Encoder) ContextTypeTime(tagNumber byte, value uint32) {
	// TODO:待完善
}

// ContextObjectID write a (context)tag / value pair where the value
// type is an unsigned int
func (e *Encoder) ContextObjectID(tagNumber byte, objectID specs.ObjectID) {
	if e.err != nil {
		return
	}
	t := tag{
		ID:      tagNumber,
		Context: true,
		Value:   4, // length of objectID is 4
		Opening: false,
		Closing: false,
	}
	encodeTag(e.buf, t)
	v, err := objectID.Encode()
	if err != nil {
		e.err = err
		return
	}
	_ = binary.Write(e.buf, binary.BigEndian, v)
}

// ContextOpening
func (e *Encoder) ContextOpening(tagNumber byte) {
	t := tag{
		ID:      tagNumber,
		Context: true,
		Opening: true,
	}

	encodeTag(e.buf, t)
}

// ContextClosing
func (e *Encoder) ContextClosing(tagNumber byte) {
	t := tag{
		ID:      tagNumber,
		Context: true,
		Closing: true,
	}

	encodeTag(e.buf, t)
}

// AppData writes a tag and value of any standard baetyl-bacnet application
// data type. Returns an error if v if of a invalid type
func (e *Encoder) AppData(v interface{}) {
	if e.err != nil {
		return
	}
	if v == nil {
		t := tag{ID: applicationTagNull}
		encodeTag(e.buf, t)
		return
	}
	switch val := v.(type) {
	case float64, bool:
		e.err = fmt.Errorf("not implemented ")
	case float32:
		t := tag{ID: applicationTagReal, Value: 4}
		encodeTag(e.buf, t)
		_ = binary.Write(e.buf, binary.BigEndian, val)
	case string:
		//+1 because there will be one byte for the string encoding format
		t := tag{ID: applicationTagCharacterString, Value: uint32(len(val) + 1)}
		encodeTag(e.buf, t)
		_ = e.buf.WriteByte(utf8Encoding)
		_, _ = e.buf.Write([]byte(val))
	case uint32:
		length := valueLength(val)
		t := tag{ID: applicationTagUnsignedInt, Value: uint32(length)}
		encodeTag(e.buf, t)
		unsigned(e.buf, val)
	case specs.SegmentationSupport:
		v := uint32(val)
		length := valueLength(v)
		t := tag{ID: applicationTagEnumerated, Value: uint32(length)}
		encodeTag(e.buf, t)
		unsigned(e.buf, v)
	case specs.ObjectID:
		t := tag{ID: applicationTagObjectID, Value: 4}
		encodeTag(e.buf, t)
		v, err := val.Encode()
		if err != nil {
			e.err = err
			return
		}
		_ = binary.Write(e.buf, binary.BigEndian, v)
	default:
		e.err = fmt.Errorf("encodeAppdata: unknown type %T", v)
	}
}

func (e *Encoder) ContextAbstractType(tagNumber byte, v specs.PropertyValue) error {
	encodeTag(e.buf, tag{ID: tagNumber, Context: true, Opening: true})

	switch v.Type {
	case specs.TypeNull:
		e.ContextNull(byte(v.Type))
	case specs.TypeBoolean:
		val, ok := v.Value.(bool)
		if !ok {
			return fmt.Errorf("wrong value, value:[%v]", reflect.ValueOf(v.Value).Type())
		}
		e.ContextBoolean(byte(v.Type), val)
	case specs.TypeUnsignedInt:
		val, ok := v.Value.(uint32)
		if !ok {
			return fmt.Errorf("wrong value, value:[%v]", reflect.ValueOf(v.Value).Type())
		}
		e.ContextUnsigned(byte(v.Type), val)
	case specs.TypeSignedInt:
		val, ok := v.Value.(int32)
		if !ok {
			return fmt.Errorf("wrong value, value:[%v]", reflect.ValueOf(v.Value).Type())
		}
		e.ContextSigned(byte(v.Type), val)
	case specs.TypeReal:
		val, ok := v.Value.(float32)
		if !ok {
			return fmt.Errorf("wrong value, value:[%v]", reflect.ValueOf(v.Value).Type())
		}
		e.ContextTypeReal(byte(v.Type), val)
	case specs.TypeDouble:
		val, ok := v.Value.(float32)
		if !ok {
			return fmt.Errorf("wrong value, value:[%v]", reflect.ValueOf(v.Value).Type())
		}
		e.ContextTypeDouble(byte(v.Type), val)
	case specs.TypeOctetString:
	case specs.TypeCharacterString:
	case specs.TypeBitString:
	case specs.TypeEnumerated:
		val, ok := v.Value.(uint32)
		if !ok {
			return fmt.Errorf("wrong value, value:[%v]", reflect.ValueOf(v.Value).Type())
		}
		e.ContextTypeEnumerated(byte(v.Type), val)
	case specs.TypeDate:
	case specs.TypeTime:
	case specs.TypeObjectID:
	default:
		return fmt.Errorf("wrong Type. type:[%d]", v.Type)
	}
	encodeTag(e.buf, tag{ID: tagNumber, Context: true, Closing: true})
	return nil
	// encodeTag(e.buf, tag{ID: tagNumber, Context: true, Opening: true})
	// // length := valueLength(v.Value)
	// length := 4
	// t := tag{ID: byte(v.Type), Value: uint32(length)}
	// encodeTag(e.buf, t)
	// // unsigned(e.buf, v.Value)
	// binary.Write(e.buf, binary.BigEndian, v.Value)
	// encodeTag(e.buf, tag{ID: tagNumber, Context: true, Closing: true})
}

// valueLength caclulates how large the necessary value needs to be to fit in the appropriate
// packet length
func valueLength(value uint32) int {
	/* length of enumerated is variable, as per 20.2.11 */
	// return binary.Size(value)
	if value < 0x100 {
		return 1
	} else if value < 0x10000 {
		return 2
	} else if value < 0x1000000 {
		return 3
	}
	return 4
}

// unsigned writes the value in the buffer using a variabled-sized encoding
func unsigned(buf *bytes.Buffer, value uint32) int {
	switch {
	case value < 0x100:
		buf.WriteByte(uint8(value))
		return 1
	case value < 0x10000:
		_ = binary.Write(buf, binary.BigEndian, uint16(value))
		return 2
	case value < 0x1000000:
		// There is no default 24 bit integer in go, so we have to
		// write it manually (in big endian)
		buf.WriteByte(byte(value >> 16))
		_ = binary.Write(buf, binary.BigEndian, uint16(value))
		return 3
	default:
		_ = binary.Write(buf, binary.BigEndian, value)
		return 4
	}
}
