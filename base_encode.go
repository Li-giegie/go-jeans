package go_jeans

import (
	"encoding/binary"
	"math"
	"unsafe"
)

// EncodeBaseFaster 仅编码基本类型
func EncodeBaseFaster(buf []byte, args ...interface{}) ([]byte, error) {
	for i := 0; i < len(args); i++ {
		switch v := args[i].(type) {
		case string:
			var tmpBuffer []byte
			*(*string)(unsafe.Pointer(&tmpBuffer)) = v
			*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&tmpBuffer)) + 2*unsafe.Sizeof(&tmpBuffer))) = len(v)
			buf = binary.LittleEndian.AppendUint32(buf, uint32(len(v)))
			buf = append(buf, tmpBuffer...)
		case int:
			buf = binary.LittleEndian.AppendUint64(buf, uint64(v))
		case []byte:
			buf = binary.LittleEndian.AppendUint32(buf, uint32(len(v)))
			buf = append(buf, v...)
		case int8:
			buf = append(buf, uint8(v))
		case int16:
			buf = binary.LittleEndian.AppendUint16(buf, uint16(v))
		case int32:
			buf = binary.LittleEndian.AppendUint32(buf, uint32(v))
		case int64:
			buf = binary.LittleEndian.AppendUint64(buf, uint64(v))
		case uint:
			buf = binary.LittleEndian.AppendUint64(buf, uint64(v))
		case uint8:
			buf = append(buf, []byte{v}...)
		case uint16:
			buf = binary.LittleEndian.AppendUint16(buf, v)
		case uint32:
			buf = binary.LittleEndian.AppendUint32(buf, v)
		case uint64:
			buf = binary.LittleEndian.AppendUint64(buf, v)
		case float32:
			buf = binary.LittleEndian.AppendUint32(buf, math.Float32bits(v))
		case float64:
			buf = binary.LittleEndian.AppendUint64(buf, math.Float64bits(v))
		case bool:
			if v {
				buf = append(buf, TRUE)
				continue
			}
			buf = append(buf, FALSE)
		default:
			return nil, encodeError(i)
		}
	}
	return buf, nil
}

// EncodeBase 仅编码基本类型
func EncodeBase(args ...interface{}) ([]byte, error) {
	buf, err := EncodeBaseFaster(make([]byte, 0, BufferSize), args...)
	if err == nil {
		BufferSize = len(buf)
	}
	return buf, err
}
