package go_jeans

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

var ErrOfBytesToBaseType_float = errors.New("float err: of Decode float bounds out of max or min value")
var ErrOfBytesToBaseType_String = errors.New("string err: of Decode resolution length is greater than the remaining length")
var ErrOfBytesToBaseType_SliceBytes = errors.New("slice byte err: of Decode  resolution length is greater than the remaining length")

// BaseTypeToBytesBufferSize 定义用于存放序列化基础类型后的字节切片的缓冲区大小
var BaseTypeToBytesBufferSize = 128

// PacketHerderLenType 头长度
type PacketHerderLenType byte

const (
	PacketHerderLen_16 = iota
	PacketHerderLen_32
	PacketHerderLen_64
)

func Write(w io.Writer, buf []byte) error {
	_, err := w.Write(Pack(buf))
	return err
}

// WriteN 自定义消息头长度模式支持2字节、4字节、8字节
func WriteN(w io.Writer, buf []byte, whl PacketHerderLenType) error {
	buf, err := PackN(buf, whl)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

func Read(r io.Reader) ([]byte, error) {
	return Unpack(r)
}

func ReadN(r io.Reader, whl PacketHerderLenType) ([]byte, error) {
	return UnpackN(r, whl)
}

// Pack 将一个字节切片重新封装成：4个字节长度+buf，的新buf（数据包）
func Pack(buf []byte) []byte {
	var hl = make([]byte, 4)
	binary.LittleEndian.PutUint32(hl, uint32(len(buf)))
	return append(hl, buf...)
}

// PackN 将一个字节切片重新封装成：自定义长度长度（plen（入参））+ +buf，的新buf（数据包）
func PackN(buf []byte, pLen PacketHerderLenType) ([]byte, error) {
	var bufBinayLen []byte
	switch pLen {
	case PacketHerderLen_16:
		if len(buf)+2 > math.MaxUint16 {
			return nil, errors.New("WriteHerderLen_16 overflow")
		}
		bufBinayLen = make([]byte, 2)
		binary.LittleEndian.PutUint16(bufBinayLen, uint16(len(buf)))
	case PacketHerderLen_32:
		if uint32(len(buf)+4) > uint32(math.MaxUint32) {
			return nil, errors.New("WriteHerderLen_32 overflow")
		}
		bufBinayLen = make([]byte, 4)
		binary.LittleEndian.PutUint32(bufBinayLen, uint32(len(buf)))
	case PacketHerderLen_64:
		if uint64(len(buf)+8) > uint64(math.MaxUint64) {
			return nil, errors.New("WriteHerderLen_64 overflow")
		}
		bufBinayLen = make([]byte, 8)
		binary.LittleEndian.PutUint64(bufBinayLen, uint64(len(buf)))
	}

	return append(bufBinayLen, buf...), nil
}

// Unpack 入参一个reader，返回一个有由Pack、PackN打包的完整的
func Unpack(r io.Reader) (buf []byte, err error) {
	var packHeaderLen = make([]byte, 4, 4)
	_, err = io.ReadFull(r, packHeaderLen)
	if err != nil {
		return nil, err
	}
	pl := binary.LittleEndian.Uint32(packHeaderLen)
	buf = make([]byte, pl)
	_, err = io.ReadFull(r, buf)
	return buf, err
}

// UnpackN 入参一个reader，返回一个有由Pack、PackN打包的完整的
func UnpackN(r io.Reader, pLen PacketHerderLenType) (buf []byte, err error) {
	var packHeaderLen uint
	switch pLen {
	case PacketHerderLen_16:
		lenBuf, err := read(r, 2)
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint(binary.LittleEndian.Uint16(lenBuf))
	case PacketHerderLen_32:
		lenBuf, err := read(r, 4)
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint(binary.LittleEndian.Uint32(lenBuf))
	case PacketHerderLen_64:
		lenBuf, err := read(r, 8)
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint(binary.LittleEndian.Uint64(lenBuf))
	}

	return read(r, uint32(packHeaderLen))
}

func read(r io.Reader, length uint32) ([]byte, error) {
	var tmp = make([]byte, length, length)
	_, err := io.ReadFull(r, tmp)
	return tmp, err
}
