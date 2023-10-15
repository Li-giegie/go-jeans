package go_jeans

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"unsafe"
)

var ErrOfBytesToBaseType_float = errors.New("float err: of Decode float bounds out of max or min value")
var ErrOfBytesToBaseType_String = errors.New("string err: of Decode resolution length is greater than the remaining length")
var ErrOfBytesToBaseType_SliceBytes = errors.New("slice byte err: of Decode  resolution length is greater than the remaining length")

// BaseTypeToBytesBufferSize 定义用于存放序列化基础类型后的字节切片的缓冲区大小
var BaseTypeToBytesBufferSize = 128

// PacketHerderLenType 包头长度
type PacketHerderLenType byte

const (
	PacketHerderLen_16 = iota
	PacketHerderLen_32
	PacketHerderLen_64
)

func Write(w io.Writer, buf []byte) error {
	_, err := w.Write(Pack(buf))
	return err
}

// WriteN 自定义消息头长度模式支持2字节、4字节、8字节
func WriteN(w io.Writer, buf []byte, whl PacketHerderLenType) error {
	buf, err := PackN(buf, whl)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

func Read(r io.Reader) ([]byte, error) {
	return Unpack(r)
}

func ReadN(r io.Reader, whl PacketHerderLenType) ([]byte, error) {
	return UnpackN(r, whl)
}

// Pack 将一个字节切片重新封装成：4个字节长度+buf，的新buf（数据包）
func Pack(buf []byte) []byte {
	var hl = make([]byte, 4)
	binary.LittleEndian.PutUint32(hl, uint32(len(buf)))
	return append(hl, buf...)
}

// PackN 将一个字节切片重新封装成：自定义长度长度（plen（入参））+ +buf，的新buf（数据包）
func PackN(buf []byte, pLen PacketHerderLenType) ([]byte, error) {
	var bufBinayLen []byte
	switch pLen {
	case PacketHerderLen_16:
		if len(buf)+2 > math.MaxUint16 {
			return nil, errors.New("WriteHerderLen_16 overflow")
		}
		bufBinayLen = make([]byte, 2)
		binary.LittleEndian.PutUint16(bufBinayLen, uint16(len(buf)))
	case PacketHerderLen_32:
		if uint32(len(buf)+4) > uint32(math.MaxUint32) {
			return nil, errors.New("WriteHerderLen_32 overflow")
		}
		bufBinayLen = make([]byte, 4)
		binary.LittleEndian.PutUint32(bufBinayLen, uint32(len(buf)))
	case PacketHerderLen_64:
		if uint64(len(buf)+8) > uint64(math.MaxUint64) {
			return nil, errors.New("WriteHerderLen_64 overflow")
		}
		bufBinayLen = make([]byte, 8)
		binary.LittleEndian.PutUint64(bufBinayLen, uint64(len(buf)))
	}

	return append(bufBinayLen, buf...), nil
}

// Unpack 入参一个reader，返回一个有由Pack、PackN打包的完整的
func Unpack(r io.Reader) (buf []byte, err error) {
	var packHeaderLen = make([]byte, 4, 4)
	_, err = io.ReadFull(r, packHeaderLen)
	if err != nil {
		return nil, err
	}
	pl := binary.LittleEndian.Uint32(packHeaderLen)
	buf = make([]byte, pl)
	_, err = io.ReadFull(r, buf)
	return buf, err
}

// UnpackN 入参一个reader，返回一个有由Pack、PackN打包的完整的
func UnpackN(r io.Reader, pLen PacketHerderLenType) (buf []byte, err error) {
	var packHeaderLen uint
	switch pLen {
	case PacketHerderLen_16:
		lenBuf, err := read(r, 2)
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint(binary.LittleEndian.Uint16(lenBuf))
	case PacketHerderLen_32:
		lenBuf, err := read(r, 4)
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint(binary.LittleEndian.Uint32(lenBuf))
	case PacketHerderLen_64:
		lenBuf, err := read(r, 8)
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint(binary.LittleEndian.Uint64(lenBuf))
	}

	return read(r, uint32(packHeaderLen))
}

func read(r io.Reader, length uint32) ([]byte, error) {
	var tmp = make([]byte, length, length)
	_, err := io.ReadFull(r, tmp)
	return tmp, err
}

// Encode 将go中的基本值类型进行编码，编码的参数和解码的参数顺序必须一致，只有在传递的类型不支持时会返回错误，其他情况不会，注意这一步并不打包返回的切片
//	func TestEncode(t *testing.T) {
//		var s struct{
//			A int
//			B string
//			C bool
//		}
//
//		buf ,err := Encode(s.A,s.B,s.C)
//		if err != nil {
//			return
//		}
//		fmt.Println(buf)
//	}
func Encode(args ...interface{}) ([]byte, error) {
	var buf = make([]byte, 0, BaseTypeToBytesBufferSize)
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			var hl = make([]byte, 4)
			binary.LittleEndian.PutUint32(hl, uint32(len(v)))
			var tmpBuffer []byte
			*(*string)(unsafe.Pointer(&tmpBuffer)) = v
			*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&tmpBuffer)) + 2*unsafe.Sizeof(&tmpBuffer))) = len(v)
			buf = append(buf, hl...)
			buf = append(buf, tmpBuffer...)
		case int:
			tmpBuffer := make([]byte, binary.MaxVarintLen64)
			binary.PutVarint(tmpBuffer, int64(v))
			buf = append(buf, tmpBuffer...)
		case []byte:
			var hl = make([]byte, 4)
			binary.LittleEndian.PutUint32(hl, uint32(len(v)))
			buf = append(buf, hl...)
			buf = append(buf, v...)
		case int8:
			buf = append(buf, uint8(v))
		case int16:
			tmpBuffer := make([]byte, binary.MaxVarintLen16)
			binary.PutVarint(tmpBuffer, int64(v))
			buf = append(buf, tmpBuffer...)
		case int32:
			tmpBuffer := make([]byte, binary.MaxVarintLen32)
			binary.PutVarint(tmpBuffer, int64(v))
			buf = append(buf, tmpBuffer...)
		case int64:
			tmpBuffer := make([]byte, binary.MaxVarintLen64)
			binary.PutVarint(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
		case uint:
			tmpBuffer := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmpBuffer, uint64(v))
			buf = append(buf, tmpBuffer...)
		case uint8:
			buf = append(buf, []byte{v}...)
		case uint16:
			tmpBuffer := make([]byte, 2)
			binary.LittleEndian.PutUint16(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
		case uint32:
			tmpBuffer := make([]byte, 4)
			binary.LittleEndian.PutUint32(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
		case uint64:
			tmpBuffer := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
		case float32:
			tmpBuffer := make([]byte, 4)
			binary.LittleEndian.PutUint32(tmpBuffer, math.Float32bits(v))
			buf = append(buf, tmpBuffer...)
		case float64:
			tmpBuffer := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmpBuffer, math.Float64bits(v))
			buf = append(buf, tmpBuffer...)
		case bool:
			if v {
				buf = append(buf, 1)
			} else {
				buf = append(buf, 0)
			}
		default:
			return nil, errors.New("conversion of this type is not supported")
		}
	}

	return buf, nil
}

// 同Encode相同，不同的是每当解码一个字段时会记录其长度，方便在内容发生变化后不重新编码的情况下修改
func EncodeWithLenByteItem(args ...interface{}) (buf []byte, itemLen []int32, err error) {
	buf = make([]byte, 0, BaseTypeToBytesBufferSize)
	itemLen = make([]int32, 0, len(args))
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			var hl = make([]byte, 4)
			binary.LittleEndian.PutUint32(hl, uint32(len(v)))
			var tmpBuffer []byte
			*(*string)(unsafe.Pointer(&tmpBuffer)) = v
			*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&tmpBuffer)) + 2*unsafe.Sizeof(&tmpBuffer))) = len(v)
			buf = append(buf, hl...)
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, int32(4+len(v)))
		case []byte:
			var hl = make([]byte, 4)
			binary.LittleEndian.PutUint32(hl, uint32(len(v)))
			buf = append(buf, hl...)
			buf = append(buf, v...)
			itemLen = append(itemLen, int32(4+len(v)))
		case int:
			tmpBuffer := make([]byte, binary.MaxVarintLen64)
			binary.PutVarint(tmpBuffer, int64(v))
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, int32(binary.MaxVarintLen64))
		case int8:
			buf = append(buf, uint8(v))
			itemLen = append(itemLen, 1)
		case int16:
			tmpBuffer := make([]byte, binary.MaxVarintLen16)
			binary.PutVarint(tmpBuffer, int64(v))
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, int32(binary.MaxVarintLen16))
		case int32:
			tmpBuffer := make([]byte, binary.MaxVarintLen32)
			binary.PutVarint(tmpBuffer, int64(v))
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, int32(binary.MaxVarintLen32))
		case int64:
			tmpBuffer := make([]byte, binary.MaxVarintLen64)
			binary.PutVarint(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, int32(binary.MaxVarintLen64))
		case uint:
			tmpBuffer := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmpBuffer, uint64(v))
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, 8)
		case uint8:
			buf = append(buf, []byte{v}...)
			itemLen = append(itemLen, 1)
		case uint16:
			tmpBuffer := make([]byte, 2)
			binary.LittleEndian.PutUint16(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, 2)
		case uint32:
			tmpBuffer := make([]byte, 4)
			binary.LittleEndian.PutUint32(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, 4)
		case uint64:
			tmpBuffer := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, 8)
		case float32:
			tmpBuffer := make([]byte, 4)
			binary.LittleEndian.PutUint32(tmpBuffer, math.Float32bits(v))
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, 4)
		case float64:
			tmpBuffer := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmpBuffer, math.Float64bits(v))
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, 8)
		case bool:
			if v {
				buf = append(buf, 1)
			} else {
				buf = append(buf, 0)
			}
			itemLen = append(itemLen, 1)
		default:
			return nil, nil, errors.New("conversion of this type is not supported")
		}
	}

	return buf, itemLen, nil
}

// Decode 将一个字节切片序列化成入参的值，参数要求是GO的基本类型，指针形式传递，编码的参数和解码的参数顺序必须一致
//
//	func TestDecode(t *testing.T) {
//		var s struct{
//			A int
//			B string
//			C bool
//		}
//		buf, _ := Encode(s.A,s.B,s.C)
//		err := Decode(buf,&s.A,&s.B,&s.C)
//		if err != nil {
//			return
//		}
//		fmt.Println(s)
//	}
func Decode(buf []byte, args ...interface{}) error {
	var index int
	for i := 0; i < len(args); i++ {
		switch v := args[i].(type) {
		case *string:
			l := binary.LittleEndian.Uint32(buf[index : index+4])
			if int(l) > len(buf[index+4:]) {
				return ErrOfBytesToBaseType_String
			}
			b := buf[index+4 : int(l)+index+4]
			*v = *(*string)(unsafe.Pointer(&b))
			index += 4 + int(l)
		case *int:
			n, _ := binary.Varint(buf[index : index+binary.MaxVarintLen64])
			*v = int(n)
			index += binary.MaxVarintLen64
		case *[]byte:
			l := binary.LittleEndian.Uint32(buf[index : index+4])
			if int(l) > len(buf[index+4:]) {
				return ErrOfBytesToBaseType_SliceBytes
			}
			*v = buf[index+4 : int(l)+index+4]
			index += 4 + int(l)
		case *bool:
			if buf[index] == 1 {
				*v = true
			} else {
				*v = false
			}
			index++
		case *int8:
			*v = int8(buf[index])
			index++
		case *int16:
			n, _ := binary.Varint(buf[index : index+binary.MaxVarintLen16])
			*v = int16(n)
			index += binary.MaxVarintLen16
		case *int32:
			n, _ := binary.Varint(buf[index : index+binary.MaxVarintLen32])
			*v = int32(n)
			index += binary.MaxVarintLen32
		case *int64:
			n, _ := binary.Varint(buf[index : index+binary.MaxVarintLen64])
			*v = n
			index += binary.MaxVarintLen64
		case *uint:
			*v = uint(binary.LittleEndian.Uint64(buf[index : index+8]))
			index += 8
		case *uint8:
			*v = buf[index]
			index += 1
		case *uint16:
			*v = binary.LittleEndian.Uint16(buf[index : index+2])
			index += 2
		case *uint32:
			*v = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
		case *uint64:
			*v = binary.LittleEndian.Uint64(buf[index : index+8])
			index += 8
		case *float32:
			n := binary.LittleEndian.Uint32(buf[index : index+4])
			if float32(n) > math.MaxFloat32 {
				return ErrOfBytesToBaseType_float
			}
			*v = math.Float32frombits(n)
			index += 4
		case *float64:
			n := binary.LittleEndian.Uint64(buf[index : index+8])
			if float64(n) > math.MaxUint64 {
				return ErrOfBytesToBaseType_float
			}
			*v = math.Float64frombits(n)
			index += 8
		default:
			return errors.New("conversion of this type is not supported")
		}

	}
	return nil
}
