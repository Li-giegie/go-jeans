package go_jeans

import (
	"math"
)

// EncodeBaseFaster 仅编码基本类型
func EncodeBaseFaster(buf []byte, args ...interface{}) ([]byte, error) {
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
				buf = append(buf, _true)
				continue
			}
			buf = append(buf, _false)
		default:
			return nil, &InvalidType{v: v, i: i}
		}
	}
	return buf, nil
}

// EncodeBase 仅编码基本类型
func EncodeBase(args ...interface{}) ([]byte, error) {
	buf, err := EncodeBaseFaster(make([]byte, 0, BaseBufferSize), args...)
	if err == nil {
		BaseBufferSize = len(buf)
	}
	return buf, err
}
