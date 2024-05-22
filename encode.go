package go_jeans

import (
	"math"
)

// Encode 编码基本类型或切片类型
func Encode(args ...interface{}) ([]byte, error) {
	buf, err := EncodeFaster(make([]byte, 0, BufferSize), args...)
	if err == nil {
		BufferSize = len(buf)
	}
	return buf, err
}

// EncodeFaster 编码基本类型或切片类型，buf参数为一个缓冲字节切片，长度为0容量为预估值，例如 make([]byte,0,128)
func EncodeFaster(buf []byte, args ...interface{}) ([]byte, error) {
	var length, j int
	for i := 0; i < len(args); i++ {
		switch v := args[i].(type) {
		case string:
			buf = littleAppendUint32(buf, uint32(len(v)))
			buf = append(buf, *stringToBytes(&v)...)
		case int:
			buf = littleAppendUint64(buf, uint64(v))
		case []byte:
			buf = littleAppendUint32(buf, uint32(len(v)))
			buf = append(buf, v...)
		case int8:
			buf = append(buf, uint8(v))
		case int16:
			buf = littleAppendUint16(buf, uint16(v))
		case int32:
			buf = littleAppendUint32(buf, uint32(v))
		case int64:
			buf = littleAppendUint64(buf, uint64(v))
		case uint:
			buf = littleAppendUint64(buf, uint64(v))
		case uint8:
			buf = append(buf, []byte{v}...)
		case uint16:
			buf = littleAppendUint16(buf, v)
		case uint32:
			buf = littleAppendUint32(buf, v)
		case uint64:
			buf = littleAppendUint64(buf, v)
		case float32:
			buf = littleAppendUint32(buf, math.Float32bits(v))
		case float64:
			buf = littleAppendUint64(buf, math.Float64bits(v))
		case bool:
			if v {
				buf = append(buf, TRUE)
				continue
			}
			buf = append(buf, FALSE)
		case []uint:
			length = len(v)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for j = 0; j < length; j++ {
					buf = littleAppendUint64(buf, uint64(v[j]))
				}
			}
		case []int:
			length = len(v)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for j = 0; j < length; j++ {
					buf = littleAppendUint64(buf, uint64(v[j]))
				}
			}
		case []uint16:
			length = len(v)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for j = 0; j < length; j++ {
					buf = littleAppendUint16(buf, v[j])
				}
			}
		case []uint32:
			length = len(v)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for j = 0; j < length; j++ {
					buf = littleAppendUint32(buf, v[j])
				}
			}
		case []uint64:
			length = len(v)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for j = 0; j < length; j++ {
					buf = littleAppendUint64(buf, v[j])
				}
			}
		case []int8:
			length = len(v)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for j = 0; j < length; j++ {
					buf = append(buf, uint8(v[j]))
				}
			}
		case []int16:
			length = len(v)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for j = 0; j < length; j++ {
					buf = littleAppendUint16(buf, uint16(v[j]))
				}
			}
		case []int32:
			length = len(v)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for j = 0; j < length; j++ {
					buf = littleAppendUint32(buf, uint32(v[j]))
				}
			}
		case []int64:
			length = len(v)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for j = 0; j < length; j++ {
					buf = littleAppendUint64(buf, uint64(v[j]))
				}
			}
		case []float32:
			length = len(v)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for j = 0; j < length; j++ {
					buf = littleAppendUint32(buf, math.Float32bits(v[j]))
				}
			}
		case []float64:
			length = len(v)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for j = 0; j < length; j++ {
					buf = littleAppendUint64(buf, math.Float64bits(v[j]))
				}
			}
		case []bool:
			length = len(v)
			buf = littleAppendUint32(buf, uint32(length))
			if length != 0 {
				for j = 0; j < length; j++ {
					if v[j] {
						buf = append(buf, TRUE)
						continue
					}
					buf = append(buf, FALSE)
				}
			}
		case []string:
			length = len(v)
			buf = littleAppendUint32(buf, uint32(length))
			if length > 0 {
				var itemLen int
				for j = 0; j < length; j++ {
					itemLen = len(v[j])
					buf = littleAppendUint32(buf, uint32(itemLen))
					if itemLen > 0 {
						buf = append(buf, *stringToBytes(&v[j])...)
					}
				}
			}
		//case []interface{}:
		//	tmp, err := EncodeFaster(make([]byte, 0, SliceBufferSize), v...)
		//	if err != nil {
		//		return nil, err
		//	}
		//	buf = append(buf, tmp...)
		default:
			return nil, encodeError(i)
		}
	}
	return buf, nil
}
