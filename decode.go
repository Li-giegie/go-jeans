package go_jeans

import (
	"encoding/binary"
	"math"
	"unsafe"
)

// Decode 解码一个切片字节，参数args为可变参数，支持类型基本类型、切片类型
func Decode(buf []byte, args ...interface{}) error {
	var index int
	var length, j uint32
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
		case *[]uint:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*v = make([]uint, length)
				for j = 0; j < length; j++ {
					(*v)[j] = uint(binary.LittleEndian.Uint64(buf[index : index+8]))
					index += 8
				}
			}
		case *[]uint8:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			if length > 0 {
				*v = buf[index+4 : index+4+int(length)]
			}
			index += 4 + int(length)
		case *[]uint16:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*v = make([]uint16, length)
				for j = 0; j < length; j++ {
					(*v)[j] = binary.LittleEndian.Uint16(buf[index : index+2])
					index += 2
				}
			}
		case *[]uint32:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*v = make([]uint32, length)
				for j = 0; j < length; j++ {
					(*v)[j] = binary.LittleEndian.Uint32(buf[index : index+4])
					index += 4
				}
			}
		case *[]uint64:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*v = make([]uint64, length)
				for j = 0; j < length; j++ {
					(*v)[j] = binary.LittleEndian.Uint64(buf[index : index+8])
					index += 8
				}
			}

		case *[]int:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*v = make([]int, length)
				for j = 0; j < length; j++ {
					(*v)[j] = int(binary.LittleEndian.Uint64(buf[index : index+8]))
					index += 8
				}
			}

		case *[]int8:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*v = make([]int8, length)
				for j = 0; j < length; j++ {
					(*v)[j] = int8(buf[index])
					index++
				}
			}

		case *[]int16:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*v = make([]int16, length)
				for j = 0; j < length; j++ {
					(*v)[j] = int16(binary.LittleEndian.Uint16(buf[index : index+2]))
					index += 2
				}
			}

		case *[]int32:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*v = make([]int32, length)
				for j = 0; j < length; j++ {
					(*v)[j] = int32(binary.LittleEndian.Uint32(buf[index : index+4]))
					index += 4
				}
			}

		case *[]int64:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*v = make([]int64, length)
				for j = 0; j < length; j++ {
					(*v)[j] = int64(binary.LittleEndian.Uint64(buf[index : index+8]))
					index += 8
				}
			}

		case *[]float32:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*v = make([]float32, length)
				for j = 0; j < length; j++ {
					n := binary.LittleEndian.Uint32(buf[index : index+4])
					if float32(n) > math.MaxFloat32 {
						return ErrOfBytesToBaseType_float
					}
					(*v)[j] = math.Float32frombits(n)
					index += 4
				}
			}
		case *[]float64:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*v = make([]float64, length)
				for j = 0; j < length; j++ {
					n := binary.LittleEndian.Uint64(buf[index : index+8])
					if float64(n) > math.MaxFloat64 {
						return ErrOfBytesToBaseType_float
					}
					(*v)[j] = math.Float64frombits(n)
					index += 8
				}
			}
		case *[]bool:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*v = make([]bool, length)
				for j = 0; j < length; j++ {
					if buf[index] == TRUE {
						(*v)[j] = true
						index++
						continue
					}
					index++
					(*v)[j] = false
				}
			}

		case *[]string:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*v = make([]string, length)
				var itemLen uint32
				for j = 0; j < length; j++ {
					itemLen = binary.LittleEndian.Uint32(buf[index : index+4])
					index += 4
					(*v)[j] = *bytesToString(buf[index : index+int(itemLen)])
					index += int(itemLen)
				}
			}
		default:
			return decodeError(i)
		}
	}
	return nil
}
