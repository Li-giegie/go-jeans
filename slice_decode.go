package go_jeans

import (
	"encoding/binary"
	"errors"
	"math"
	"strconv"
)

// DecodeSlice 将切片类型的编码还原成切片，支持列表[*[]uint32]
func DecodeSlice(buf []byte, slice ...interface{}) error {
	var index int
	var length, i uint32
	for _i, item := range slice {
		switch sv := item.(type) {
		case *[]uint:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]uint, length)
				for i = 0; i < length; i++ {
					(*sv)[i] = uint(binary.LittleEndian.Uint64(buf[index : index+8]))
					index += 8
				}
			}
		case *[]uint8:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			if length > 0 {
				*sv = buf[index+4 : index+4+int(length)]
			}
			index += 4 + int(length)
		case *[]uint16:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]uint16, length)
				for i = 0; i < length; i++ {
					(*sv)[i] = binary.LittleEndian.Uint16(buf[index : index+2])
					index += 2
				}
			}
		case *[]uint32:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]uint32, length)
				for i = 0; i < length; i++ {
					(*sv)[i] = binary.LittleEndian.Uint32(buf[index : index+4])
					index += 4
				}
			}
		case *[]uint64:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]uint64, length)
				for i = 0; i < length; i++ {
					(*sv)[i] = binary.LittleEndian.Uint64(buf[index : index+8])
					index += 8
				}
			}

		case *[]int:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]int, length)
				for i = 0; i < length; i++ {
					(*sv)[i] = int(binary.LittleEndian.Uint64(buf[index : index+8]))
					index += 8
				}
			}

		case *[]int8:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]int8, length)
				for i = 0; i < length; i++ {
					(*sv)[i] = int8(buf[index])
					index++
				}
			}

		case *[]int16:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]int16, length)
				for i = 0; i < length; i++ {
					(*sv)[i] = int16(binary.LittleEndian.Uint16(buf[index : index+2]))
					index += 2
				}
			}

		case *[]int32:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]int32, length)
				for i = 0; i < length; i++ {
					(*sv)[i] = int32(binary.LittleEndian.Uint32(buf[index : index+4]))
					index += 4
				}
			}

		case *[]int64:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]int64, length)
				for i = 0; i < length; i++ {
					(*sv)[i] = int64(binary.LittleEndian.Uint64(buf[index : index+8]))
					index += 8
				}
			}

		case *[]float32:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]float32, length)
				for i = 0; i < length; i++ {
					n := binary.LittleEndian.Uint32(buf[index : index+4])
					if float32(n) > math.MaxFloat32 {
						return ErrOfBytesToBaseType_float
					}
					(*sv)[i] = math.Float32frombits(n)
					index += 4
				}
			}
		case *[]float64:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]float64, length)
				for i = 0; i < length; i++ {
					n := binary.LittleEndian.Uint64(buf[index : index+8])
					if float64(n) > math.MaxFloat64 {
						return ErrOfBytesToBaseType_float
					}
					(*sv)[i] = math.Float64frombits(n)
					index += 8
				}
			}
		case *[]bool:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]bool, length)
				for i = 0; i < length; i++ {
					if buf[index] == TRUE {
						(*sv)[i] = true
						index++
						continue
					}
					index++
					(*sv)[i] = false
				}
			}

		case *[]string:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]string, length)
				var itemLen uint32
				for i = 0; i < length; i++ {
					itemLen = binary.LittleEndian.Uint32(buf[index : index+4])
					index += 4
					(*sv)[i] = *bytesToString(buf[index : index+int(itemLen)])
					index += int(itemLen)
				}
			}
		default:
			return decodeError(_i)
		}
	}
	return nil
}

func decodeError(i int) error {
	return errors.New("decode err: index is " + strconv.Itoa(i) + "unsupported type ")
}
