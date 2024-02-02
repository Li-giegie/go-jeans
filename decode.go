package go_jeans

import (
	"encoding/binary"
	"errors"
	"math"
	"strconv"
	"unsafe"
)

func DecodeV2(buf []byte, args ...interface{}) error {
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
			*v = int(binary.LittleEndian.Uint64(buf[index : index+8]))
			index += 8
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
			*v = int16(binary.LittleEndian.Uint16(buf[index : index+2]))
			index += 2
		case *int32:
			*v = int32(binary.LittleEndian.Uint32(buf[index : index+4]))
			index += 4
		case *int64:
			*v = int64(binary.LittleEndian.Uint64(buf[index : index+8]))
			index += 8
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
			return decodeError(i)
		}
	}
	return nil
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
			return decodeError(i)
		}

	}
	return nil
}

// DecodeSlice 将切片类型的编码还原成切片，支持列表[*[]uint32]
func DecodeSlice(buf []byte, slice ...interface{}) error {
	var index int
	for _i, item := range slice {
		switch sv := item.(type) {
		case *[]uint32:
			l := binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			var tmp = make([]uint32, l)
			for i := 0; i < int(l); i++ {
				tmp[i] = binary.LittleEndian.Uint32(buf[index : index+4])
				index += 4
			}
			*sv = tmp
		default:
			return decodeError(_i)
		}
	}
	return nil
}

// CountLength 统计字段的长度，可用于定义缓冲区容量
// 例如：
// var a,b,c string
// n,_ := CountLength(a,b,c)
// buf := make([]byte,0,n)
// buf,_= EncodeV2(buf,a,b,c)
func CountLength(args ...interface{}) (length int, err error) {
	for i := 0; i < len(args); i++ {
		switch v := args[i].(type) {
		case string:
			length += 4 + len(v)
		case []byte:
			length += 4 + len(v)
		case int8, uint8, bool:
			length++
		case int16, uint16:
			length += 2
		case int32, uint32, float32:
			length += 4
		case int, uint, int64, uint64, float64:
			length += 8
		default:
			return 0, decodeError(i)
		}
	}
	return
}

func decodeError(i int) error {
	return errors.New("decode err: index is " + strconv.Itoa(i) + "unsupported type ")
}
