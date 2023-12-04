package go_jeans

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"unsafe"
)

// Encode 将go中的基本值类型进行编码，编码的参数和解码的参数顺序必须一致，只有在传递的类型不支持时会返回错误，其他情况不会，注意这一步并不打包返回的切片
//
//	func TestEncode(t *testing.T) {
//		var s struct{
//			A int
//			B string
//			C bool
//		}
//
//		buf ,err := Encode(s.A,s.B,s.C)
//		if err != nil {
//			return
//		}
//		fmt.Println(buf)
//	}
func Encode(args ...interface{}) ([]byte, error) {
	var buf = make([]byte, 0, BaseTypeToBytesBufferSize)
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			var hl = make([]byte, 4)
			binary.LittleEndian.PutUint32(hl, uint32(len(v)))
			var tmpBuffer []byte
			*(*string)(unsafe.Pointer(&tmpBuffer)) = v
			*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&tmpBuffer)) + 2*unsafe.Sizeof(&tmpBuffer))) = len(v)
			buf = append(buf, hl...)
			buf = append(buf, tmpBuffer...)
		case int:
			tmpBuffer := make([]byte, binary.MaxVarintLen64)
			binary.PutVarint(tmpBuffer, int64(v))
			buf = append(buf, tmpBuffer...)
		case []byte:
			var hl = make([]byte, 4)
			binary.LittleEndian.PutUint32(hl, uint32(len(v)))
			buf = append(buf, hl...)
			buf = append(buf, v...)
		case int8:
			buf = append(buf, uint8(v))
		case int16:
			tmpBuffer := make([]byte, binary.MaxVarintLen16)
			binary.PutVarint(tmpBuffer, int64(v))
			buf = append(buf, tmpBuffer...)
		case int32:
			tmpBuffer := make([]byte, binary.MaxVarintLen32)
			binary.PutVarint(tmpBuffer, int64(v))
			buf = append(buf, tmpBuffer...)
		case int64:
			tmpBuffer := make([]byte, binary.MaxVarintLen64)
			binary.PutVarint(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
		case uint:
			tmpBuffer := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmpBuffer, uint64(v))
			buf = append(buf, tmpBuffer...)
		case uint8:
			buf = append(buf, []byte{v}...)
		case uint16:
			tmpBuffer := make([]byte, 2)
			binary.LittleEndian.PutUint16(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
		case uint32:
			tmpBuffer := make([]byte, 4)
			binary.LittleEndian.PutUint32(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
		case uint64:
			tmpBuffer := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
		case float32:
			tmpBuffer := make([]byte, 4)
			binary.LittleEndian.PutUint32(tmpBuffer, math.Float32bits(v))
			buf = append(buf, tmpBuffer...)
		case float64:
			tmpBuffer := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmpBuffer, math.Float64bits(v))
			buf = append(buf, tmpBuffer...)
		case bool:
			if v {
				buf = append(buf, 1)
			} else {
				buf = append(buf, 0)
			}
		default:
			return nil, errors.New("conversion of this type is not supported")
		}
	}

	return buf, nil
}

// 同Encode相同，不同的是每当解码一个字段时会记录其长度，方便在内容发生变化后不重新编码的情况下修改
func EncodeWithLenByteItem(args ...interface{}) (buf []byte, itemLen []int32, err error) {
	buf = make([]byte, 0, BaseTypeToBytesBufferSize)
	itemLen = make([]int32, 0, len(args))
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			var hl = make([]byte, 4)
			binary.LittleEndian.PutUint32(hl, uint32(len(v)))
			var tmpBuffer []byte
			*(*string)(unsafe.Pointer(&tmpBuffer)) = v
			*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&tmpBuffer)) + 2*unsafe.Sizeof(&tmpBuffer))) = len(v)
			buf = append(buf, hl...)
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, int32(4+len(v)))
		case []byte:
			var hl = make([]byte, 4)
			binary.LittleEndian.PutUint32(hl, uint32(len(v)))
			buf = append(buf, hl...)
			buf = append(buf, v...)
			itemLen = append(itemLen, int32(4+len(v)))
		case int:
			tmpBuffer := make([]byte, binary.MaxVarintLen64)
			binary.PutVarint(tmpBuffer, int64(v))
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, int32(binary.MaxVarintLen64))
		case int8:
			buf = append(buf, uint8(v))
			itemLen = append(itemLen, 1)
		case int16:
			tmpBuffer := make([]byte, binary.MaxVarintLen16)
			binary.PutVarint(tmpBuffer, int64(v))
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, int32(binary.MaxVarintLen16))
		case int32:
			tmpBuffer := make([]byte, binary.MaxVarintLen32)
			binary.PutVarint(tmpBuffer, int64(v))
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, int32(binary.MaxVarintLen32))
		case int64:
			tmpBuffer := make([]byte, binary.MaxVarintLen64)
			binary.PutVarint(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, int32(binary.MaxVarintLen64))
		case uint:
			tmpBuffer := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmpBuffer, uint64(v))
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, 8)
		case uint8:
			buf = append(buf, []byte{v}...)
			itemLen = append(itemLen, 1)
		case uint16:
			tmpBuffer := make([]byte, 2)
			binary.LittleEndian.PutUint16(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, 2)
		case uint32:
			tmpBuffer := make([]byte, 4)
			binary.LittleEndian.PutUint32(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, 4)
		case uint64:
			tmpBuffer := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmpBuffer, v)
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, 8)
		case float32:
			tmpBuffer := make([]byte, 4)
			binary.LittleEndian.PutUint32(tmpBuffer, math.Float32bits(v))
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, 4)
		case float64:
			tmpBuffer := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmpBuffer, math.Float64bits(v))
			buf = append(buf, tmpBuffer...)
			itemLen = append(itemLen, 8)
		case bool:
			if v {
				buf = append(buf, 1)
			} else {
				buf = append(buf, 0)
			}
			itemLen = append(itemLen, 1)
		default:
			return nil, nil, errors.New("conversion of this type is not supported")
		}
	}

	return buf, itemLen, nil
}

// EncodeSlice 将入参切片编码：支持列表[[]uint32]
func EncodeSlice(slice ...interface{}) ([]byte, error) {
	var buf = make([]byte, 0, BaseTypeToBytesBufferSize)
	for _, item := range slice {
		switch sv := item.(type) {
		case []uint32:
			buf = append(buf, []byte{0, 0, 0, 0}...)
			binary.LittleEndian.PutUint32(buf[len(buf)-4:], uint32(len(sv)))
			for _, v := range sv {
				buf = append(buf, []byte{0, 0, 0, 0}...)
				binary.LittleEndian.PutUint32(buf[len(buf)-4:], v)
			}
		default:
			return nil, fmt.Errorf("[%T] conversion of this type is not supported", sv)
		}
	}
	return buf, nil
}
