package go_jeans

import (
	"encoding/binary"
	"math"
)

// DecodeBase 将一个字节切片序列化成入参的值，参数要求是GO的基本类型，指针形式传递，编码的参数和解码的参数顺序必须一致
func DecodeBase(buf []byte, args ...interface{}) error {
	var index int
	for i := 0; i < len(args); i++ {
		switch v := args[i].(type) {
		case *string:
			l := binary.LittleEndian.Uint32(buf[index : index+4])
			if int(l) > len(buf[index+4:]) {
				return ErrOfBytesToBaseType_String
			}
			*v = *bytesToString(buf[index+4 : int(l)+index+4])
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
