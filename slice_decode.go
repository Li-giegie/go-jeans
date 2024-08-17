package go_jeans

import (
	"encoding/binary"
	"math"
)

// DecodeSlice 仅解码切片类型，slice 参数为切片类型的指针例如 *[]int
func DecodeSlice(buf []byte, slice ...interface{}) error {
	var index int
	var length, j uint32
	for i := 0; i < len(slice); i++ {
		switch sv := slice[i].(type) {
		case *[]uint:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]uint, length)
				for j = 0; j < length; j++ {
					(*sv)[j] = uint(binary.LittleEndian.Uint64(buf[index : index+8]))
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
				for j = 0; j < length; j++ {
					(*sv)[j] = binary.LittleEndian.Uint16(buf[index : index+2])
					index += 2
				}
			}
		case *[]uint32:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]uint32, length)
				for j = 0; j < length; j++ {
					(*sv)[j] = binary.LittleEndian.Uint32(buf[index : index+4])
					index += 4
				}
			}
		case *[]uint64:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]uint64, length)
				for j = 0; j < length; j++ {
					(*sv)[j] = binary.LittleEndian.Uint64(buf[index : index+8])
					index += 8
				}
			}

		case *[]int:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]int, length)
				for j = 0; j < length; j++ {
					(*sv)[j] = int(binary.LittleEndian.Uint64(buf[index : index+8]))
					index += 8
				}
			}

		case *[]int8:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]int8, length)
				for j = 0; j < length; j++ {
					(*sv)[j] = int8(buf[index])
					index++
				}
			}

		case *[]int16:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]int16, length)
				for j = 0; j < length; j++ {
					(*sv)[j] = int16(binary.LittleEndian.Uint16(buf[index : index+2]))
					index += 2
				}
			}

		case *[]int32:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]int32, length)
				for j = 0; j < length; j++ {
					(*sv)[j] = int32(binary.LittleEndian.Uint32(buf[index : index+4]))
					index += 4
				}
			}

		case *[]int64:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]int64, length)
				for j = 0; j < length; j++ {
					(*sv)[j] = int64(binary.LittleEndian.Uint64(buf[index : index+8]))
					index += 8
				}
			}

		case *[]float32:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]float32, length)
				for j = 0; j < length; j++ {
					n := binary.LittleEndian.Uint32(buf[index : index+4])
					(*sv)[j] = math.Float32frombits(n)
					index += 4
				}
			}
		case *[]float64:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]float64, length)
				for j = 0; j < length; j++ {
					n := binary.LittleEndian.Uint64(buf[index : index+8])
					(*sv)[j] = math.Float64frombits(n)
					index += 8
				}
			}
		case *[]bool:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]bool, length)
				for j = 0; j < length; j++ {
					if buf[index] == _true {
						(*sv)[j] = true
						index++
						continue
					}
					index++
					(*sv)[j] = false
				}
			}

		case *[]string:
			length = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
			if length > 0 {
				*sv = make([]string, length)
				var itemLen uint32
				for j = 0; j < length; j++ {
					itemLen = binary.LittleEndian.Uint32(buf[index : index+4])
					index += 4
					(*sv)[j] = *bytesToString(buf[index : index+int(itemLen)])
					index += int(itemLen)
				}
			}
		default:
			return &InvalidType{v: sv, i: i}
		}
	}
	return nil
}
