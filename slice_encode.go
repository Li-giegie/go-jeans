package go_jeans

import (
	"errors"
	"math"
	"strconv"
)

var (
	BufferSize      = 128
	SliceBufferSize = 256
)

const (
	FALSE = iota
	TRUE
)

// EncodeSlice 仅编码切片类型
func EncodeSlice(slice ...interface{}) ([]byte, error) {
	buf, err := EncodeSliceFaster(make([]byte, 0, SliceBufferSize), slice...)
	if err == nil {
		SliceBufferSize = len(buf)
	}
	return buf, err
}

// EncodeSliceFaster 仅编码切片类型
func EncodeSliceFaster(buf []byte, slice ...interface{}) ([]byte, error) {
	var length, i int
	for index, item := range slice {
		switch sv := item.(type) {
		case []uint:
			length = len(sv)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for i = 0; i < length; i++ {
					buf = littleAppendUint64(buf, uint64(sv[i]))
				}
			}
		case []int:
			length = len(sv)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for i = 0; i < length; i++ {
					buf = littleAppendUint64(buf, uint64(sv[i]))
				}
			}
		case []uint8:
			length = len(sv)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				buf = append(buf, sv...)
			}
		case []uint16:
			length = len(sv)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for i = 0; i < length; i++ {
					buf = littleAppendUint16(buf, sv[i])
				}
			}
		case []uint32:
			length = len(sv)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for i = 0; i < length; i++ {
					buf = littleAppendUint32(buf, sv[i])
				}
			}
		case []uint64:
			length = len(sv)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for i = 0; i < length; i++ {
					buf = littleAppendUint64(buf, sv[i])
				}
			}
		case []int8:
			length = len(sv)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for i = 0; i < length; i++ {
					buf = append(buf, uint8(sv[i]))
				}
			}
		case []int16:
			length = len(sv)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for i = 0; i < length; i++ {
					buf = littleAppendUint16(buf, uint16(sv[i]))
				}
			}
		case []int32:
			length = len(sv)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for i = 0; i < length; i++ {
					buf = littleAppendUint32(buf, uint32(sv[i]))
				}
			}
		case []int64:
			length = len(sv)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for i = 0; i < length; i++ {
					buf = littleAppendUint64(buf, uint64(sv[i]))
				}
			}
		case []float32:
			length = len(sv)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for i = 0; i < length; i++ {
					buf = littleAppendUint32(buf, math.Float32bits(sv[i]))
				}
			}
		case []float64:
			length = len(sv)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for i = 0; i < length; i++ {
					buf = littleAppendUint64(buf, math.Float64bits(sv[i]))
				}
			}
		case []bool:
			length = len(sv)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for i = 0; i < length; i++ {
					if sv[i] {
						buf = append(buf, TRUE)
						continue
					}
					buf = append(buf, FALSE)
				}
			}
		case []string:
			length = len(sv)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				var itemLen int
				for i = 0; i < length; i++ {
					itemLen = len(sv[i])
					buf = littleAppendUint32(buf, uint32(itemLen))
					if itemLen > 0 {
						buf = append(buf, stringToBytes(&sv[i])...)
					}
				}
			}
		default:
			return nil, encodeError(index)
		}
	}
	return buf, nil
}

func encodeError(i int) error {
	return errors.New("encode err: index [" + strconv.Itoa(i) + "] unsupported type ")
}
