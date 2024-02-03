package go_jeans

import (
	"encoding/binary"
	"errors"
	"math"
	"strconv"
	"unsafe"
)

var BufferSize = 128

// EncodeFaster 在测试文件encode_decode_test.go中查看使用例子
func EncodeFaster(buf []byte, args ...interface{}) ([]byte, error) {
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
				buf = append(buf, 1)
				continue
			}
			buf = append(buf, 0)
		default:
			return nil, encodeError(i)
		}
	}
	return buf, nil
}

// Encode 在测试文件encode_decode_test.go中查看使用例子
func Encode(args ...interface{}) ([]byte, error) {
	buf, err := EncodeFaster(make([]byte, 0, BufferSize), args...)
	if err == nil {
		BufferSize = len(buf)
	}
	return buf, err
}

// EncodeSlice 在测试文件encode_decode_test.go中查看使用例子
func EncodeSlice(slice ...interface{}) ([]byte, error) {
	var buf = make([]byte, 0, BufferSize)
	for index, item := range slice {
		switch sv := item.(type) {
		case []uint32:
			buf = append(buf, []byte{0, 0, 0, 0}...)
			binary.LittleEndian.PutUint32(buf[len(buf)-4:], uint32(len(sv)))
			for _, v := range sv {
				buf = append(buf, []byte{0, 0, 0, 0}...)
				binary.LittleEndian.PutUint32(buf[len(buf)-4:], v)
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
