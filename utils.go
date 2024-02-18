package go_jeans

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"unsafe"
)

var ErrOfBytesToBaseType_float = errors.New("float err: of Decode float bounds out of max or min value")
var ErrOfBytesToBaseType_String = errors.New("string err: of Decode resolution length is greater than the remaining length")
var ErrOfBytesToBaseType_SliceBytes = errors.New("slice byte err: of Decode  resolution length is greater than the remaining length")

// PacketHerderLenType 头长度
type PacketHerderLenType uint8

const (
	PacketHerderLenType_uint16 PacketHerderLenType = iota
	PacketHerderLenType_uint32
	PacketHerderLenType_uint64
)

// Pack 将一个字节切片重新封装成：4个字节长度+buf，的新buf（数据包）
func Pack(buf []byte) []byte {
	var hl = make([]byte, 4)
	binary.LittleEndian.PutUint32(hl, uint32(len(buf)))
	return append(hl, buf...)
}

var PacketHerderLenErr = errors.New("invalid PacketHerderLenType")

// PackN 将一个字节切片重新封装成：自定义长度长度（plen（入参））+ +buf，的新buf（数据包）
func PackN(buf []byte, pLen PacketHerderLenType) ([]byte, error) {
	var bufBinayLen []byte
	switch pLen {
	case PacketHerderLenType_uint16:
		if len(buf)+2 > math.MaxUint16 {
			return nil, errors.New("WriteHerderLen_16 overflow")
		}
		bufBinayLen = make([]byte, 2)
		binary.LittleEndian.PutUint16(bufBinayLen, uint16(len(buf)))
	case PacketHerderLenType_uint32:
		if uint32(len(buf)+4) > uint32(math.MaxUint32) {
			return nil, errors.New("WriteHerderLen_32 overflow")
		}
		bufBinayLen = make([]byte, 4)
		binary.LittleEndian.PutUint32(bufBinayLen, uint32(len(buf)))
	case PacketHerderLenType_uint64:
		if uint64(len(buf)+8) > uint64(math.MaxUint64) {
			return nil, errors.New("WriteHerderLen_64 overflow")
		}
		bufBinayLen = make([]byte, 8)
		binary.LittleEndian.PutUint64(bufBinayLen, uint64(len(buf)))
	default:
		return nil, PacketHerderLenErr
	}

	return append(bufBinayLen, buf...), nil
}

// Unpack 入参一个reader，返回一个有由Pack、PackN打包的字节切片
func Unpack(r io.Reader) (buf []byte, err error) {
	var packHeaderLen = make([]byte, 4)
	_, err = io.ReadFull(r, packHeaderLen)
	if err != nil {
		return packHeaderLen, err
	}
	pl := binary.LittleEndian.Uint32(packHeaderLen)
	buf = make([]byte, pl)
	_, err = io.ReadFull(r, buf)
	return buf, err
}

// UnpackN 入参一个reader，返回一个有由Pack、PackN打包的完整的
func UnpackN(r io.Reader, pLen PacketHerderLenType) (buf []byte, err error) {
	var packHeaderLen uint64
	switch pLen {
	case PacketHerderLenType_uint16:
		lenBuf, err := read(r, 2)
		if err != nil {
			return lenBuf, err
		}
		packHeaderLen = uint64(binary.LittleEndian.Uint16(lenBuf))
	case PacketHerderLenType_uint32:
		lenBuf, err := read(r, 4)
		if err != nil {
			return lenBuf, err
		}
		packHeaderLen = uint64(binary.LittleEndian.Uint32(lenBuf))
	case PacketHerderLenType_uint64:
		lenBuf, err := read(r, 8)
		if err != nil {
			return lenBuf, err
		}
		packHeaderLen = binary.LittleEndian.Uint64(lenBuf)
	default:
		return nil, PacketHerderLenErr
	}
	return read(r, packHeaderLen)
}

func read(r io.Reader, length uint64) ([]byte, error) {
	var tmp = make([]byte, length)
	_, err := io.ReadFull(r, tmp)
	return tmp, err
}

// CheckField 是否支持编码 返回值 > -1 表示又不支持的字段
func CheckField(args ...interface{}) (index int, err error) {
	for i := 0; i < len(args); i++ {
		switch t := args[i].(type) {
		case string, int8, uint8, bool, int16, uint16, int32, uint32, float32, int, uint, int64, uint64, float64:
		default:
			return i, fmt.Errorf("field type: %T val: %v nonsupport", t, t)
		}
	}
	return -1, nil
}

// CheckFieldSlice 用于判断切片是否支持编码 返回值 > -1 表示又不支持的字段
func CheckFieldSlice(args ...interface{}) (index int, err error) {
	for i := 0; i < len(args); i++ {
		switch t := args[i].(type) {
		case []uint32:
		default:
			return i, fmt.Errorf("field type: %T val: %v nonsupport", t, t)
		}
	}
	return -1, nil
}

// CountLength 统计字段的长度，可用于定义缓冲区容量
// 例如：
// var a,b,c string
// n,_ := CountLength(a,b,c)
// buf := make([]byte,0,n)
// buf,_= EncodeV2(buf,a,b,c)
func CountLength(args ...interface{}) (length int) {
	for i := 0; i < len(args); i++ {
		switch v := args[i].(type) {
		case string:
			length += 4 + len(v)
		case []byte:
			length += 4 + len(v)
		case int8, uint8, bool:
			length++
		case int16, uint16:
			length += 2
		case int32, uint32, float32:
			length += 4
		case int, uint, int64, uint64, float64:
			length += 8
		default:
			panic("field nonsupport Call the CheckField or CheckFieldSlice function to get details")
		}
	}
	return
}

func stringToBytes(str *string) []byte {
	var tmpBuffer []byte
	*(*string)(unsafe.Pointer(&tmpBuffer)) = *str
	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&tmpBuffer)) + 2*unsafe.Sizeof(&tmpBuffer))) = len(*str)
	return tmpBuffer
}

func bytesToString(buf []byte) *string {
	return (*string)(unsafe.Pointer(&buf))
}
