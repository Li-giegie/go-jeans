package go_jeans

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"unsafe"
)

var (
	ErrOfBytesToBaseType_float      = errors.New("float err: of Decode float bounds out of max or min value")
	ErrOfBytesToBaseType_String     = errors.New("string err: of Decode resolution length is greater than the remaining length")
	ErrOfBytesToBaseType_SliceBytes = errors.New("slice byte err: of Decode  resolution length is greater than the remaining length")
)

// PacketHerderLenType 头长度
type PacketHerderLenType uint8

const (
	PacketHerderLenType_uint16 PacketHerderLenType = iota
	PacketHerderLenType_uint32
	PacketHerderLenType_uint64
)

var SupportList = map[string]struct{}{
	"int":        {},
	"int8":       {},
	"int16":      {},
	"int32":      {},
	"int64":      {},
	"uint":       {},
	"uint8":      {},
	"uint16":     {},
	"uint32":     {},
	"uint64":     {},
	"byte":       {},
	"float32":    {},
	"float64":    {},
	"string":     {},
	"bool":       {},
	"*int":       {},
	"*int8":      {},
	"*int16":     {},
	"*int32":     {},
	"*int64":     {},
	"*uint":      {},
	"*uint8":     {},
	"*uint16":    {},
	"*uint32":    {},
	"*uint64":    {},
	"*byte":      {},
	"*float32":   {},
	"*float64":   {},
	"*string":    {},
	"*bool":      {},
	"[]int":      {},
	"[]int8":     {},
	"[]int16":    {},
	"[]int32":    {},
	"[]int64":    {},
	"[]uint":     {},
	"[]uint8":    {},
	"[]uint16":   {},
	"[]uint32":   {},
	"[]uint64":   {},
	"[]byte":     {},
	"[]float32":  {},
	"[]float64":  {},
	"[]string":   {},
	"[]bool":     {},
	"*[]int":     {},
	"*[]int8":    {},
	"*[]int16":   {},
	"*[]int32":   {},
	"*[]int64":   {},
	"*[]uint":    {},
	"*[]uint8":   {},
	"*[]uint16":  {},
	"*[]uint32":  {},
	"*[]uint64":  {},
	"*[]byte":    {},
	"*[]float32": {},
	"*[]float64": {},
	"*[]string":  {},
	"*[]bool":    {},
}

// Pack 将一个字节切片重新封装成：4个字节长度+buf，的新buf（数据包）
func Pack(buf []byte) []byte {
	hl := make([]byte, 4)
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
	packHeaderLen := make([]byte, 4)
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
	tmp := make([]byte, length)
	_, err := io.ReadFull(r, tmp)
	return tmp, err
}

type fieldType uint8

const (
	fieldType_nonsupport = iota
	fieldType_base
	fieldType_base_pointer
	fieldType_slice
	fieldType_slice_pointer
)

func getFieldType(arg interface{}) fieldType {
	switch arg.(type) {
	case string, int8, uint8, bool, int16, uint16, int32, uint32, float32, int, uint, int64, uint64, float64:
		return fieldType_base
	case []uint, []uint8, []uint16, []uint32, []uint64, []int, []int8, []int16, []int32, []int64, []float32, []float64, []bool, []string:
		return fieldType_slice
	case *string, *int8, *uint8, *bool, *int16, *uint16, *int32, *uint32, *float32, *int, *uint, *int64, *uint64, *float64:
		return fieldType_base
	case *[]uint, *[]uint8, *[]uint16, *[]uint32, *[]uint64, *[]int, *[]int8, *[]int16, *[]int32, *[]int64, *[]float32, *[]float64, *[]bool, *[]string:
		return fieldType_slice_pointer
	default:
		return fieldType_nonsupport
	}
}

func stringToBytes(str *string) *[]byte {
	return (*[]byte)(unsafe.Pointer(str))
}

func bytesToString(buf []byte) *string {
	return (*string)(unsafe.Pointer(&buf))
}

func littleAppendUint64(b []byte, v uint64) []byte {
	return append(b,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
		byte(v>>32),
		byte(v>>40),
		byte(v>>48),
		byte(v>>56),
	)
}

func littleAppendUint32(b []byte, v uint32) []byte {
	return append(b,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
	)
}

func littleAppendUint16(b []byte, v uint16) []byte {
	return append(b,
		byte(v),
		byte(v>>8),
	)
}
