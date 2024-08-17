package go_jeans

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"unsafe"
)

// PacketType 头长度 + data
type PacketType uint8

const (
	// PacketType8 1字节 + data
	PacketType8 PacketType = iota
	PacketType16
	PacketType24
	PacketType32
	PacketType64
)

var pLimitErr = errors.New("PackN err: 'data' exceeds the limit length")
var ptErr = errors.New("PacketType err: invalid packageType")

var (
	BufferSize      = 256
	BaseBufferSize  = 128
	SliceBufferSize = 256
)

const (
	_false uint8 = iota
	_true
)

type slice struct {
	ptr unsafe.Pointer
	len int
	cap int
}

// PackN 将一个字节切片重新封装成：长度 + 数据
func PackN(data []byte, pLen PacketType) ([]byte, error) {
	var buf []byte
	var headerLen uint64
	switch pLen {
	case PacketType8:
		headerLen = math.MaxUint8
		buf = []byte{byte(len(data))}
	case PacketType16:
		headerLen = math.MaxUint16
		buf = make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, uint16(len(data)))
	case PacketType24:
		headerLen = 0x00FFFFFF
		lb := len(data)
		buf = []byte{
			byte(lb),
			byte(lb >> 8),
			byte(lb >> 16),
		}
	case PacketType32:
		headerLen = math.MaxUint32
		buf = make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(len(data)))
	case PacketType64:
		headerLen = math.MaxUint64
		buf = make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, uint64(len(data)))
	default:
		return nil, ptErr
	}
	if uint64(len(data)) > headerLen {
		return nil, pLimitErr
	}
	return append(buf, data...), nil
}

// LittleEndianPutUint24 使用小端实现对u转换成3字节的正整数切片
func LittleEndianPutUint24(n uint32, b []byte) {
	b[0] = byte(n)
	b[1] = byte(n >> 8)
	b[2] = byte(n >> 16)
}

// LittleEndianUint24 转换为3字节的uint32
func LittleEndianUint24(b []byte) (n uint32) {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}

// UnpackN 入参一个reader，返回一个有由Pack、PackN打包的完整的切片
func UnpackN(r io.Reader, pLen PacketType) (data []byte, err error) {
	return UnpackBufferN(r, make([]byte, 128), pLen)
}

func UnpackBufferN(r io.Reader, buf []byte, pLen PacketType) (data []byte, err error) {
	var packHeaderLen uint64
	switch pLen {
	case PacketType8:
		_, err = io.ReadFull(r, buf[:1])
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint64(buf[0])
	case PacketType16:
		_, err = io.ReadFull(r, buf[:2])
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint64(binary.LittleEndian.Uint16(buf[:2]))
	case PacketType24:
		_, err = io.ReadFull(r, buf[:3])
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint64(uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16)
	case PacketType32:
		_, err = io.ReadFull(r, buf[:4])
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint64(binary.LittleEndian.Uint32(buf[:4]))
	case PacketType64:
		_, err := io.ReadFull(r, buf[:8])
		if err != nil {
			return nil, err
		}
		packHeaderLen = binary.LittleEndian.Uint64(buf[:8])
	default:
		return nil, pLimitErr
	}
	if v := uint64(len(buf)); v < packHeaderLen {
		buf = append(buf, make([]byte, packHeaderLen-v)...)
	}
	_, err = io.ReadFull(r, buf[:packHeaderLen])
	return buf[:packHeaderLen], err
}

type str struct {
	s   string
	cap int
}

func stringToBytes(s *string) *[]byte {
	return (*[]byte)(unsafe.Pointer(&str{*s, len(*s)}))
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
