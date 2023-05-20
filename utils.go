package go_jeans

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

type PacketHerderLenType byte

const (
	PacketHerderLen_16 = iota
	PacketHerderLen_32
	PacketHerderLen_64
)

func Write(w io.Writer,buf []byte) error {
	_,err := w.Write(Pack(buf))
	return err
}

// 自定义消息头长度模式支持2字节、4字节、8字节
func WriteN(w io.Writer,buf []byte,whl PacketHerderLenType) error {
	buf,err := PackN(buf,whl)
	if err != nil {
		return err
	}
	_,err = w.Write(buf)
	return err
}

func Read(r io.Reader) ([]byte,error) {
	return Unpack(r)
}

func ReadN(r io.Reader,whl PacketHerderLenType) ([]byte,error) {
	return UnpackN(r,whl)
}

// 打包头的长度：32字节
func Pack(buf []byte) []byte {
	var hl = make([]byte,4)
	binary.LittleEndian.PutUint32(hl,uint32(len(buf)))
	return append(hl,buf...)
}

// 自定义打包头的长度：16、32、64字节
func PackN(buf []byte,pLen PacketHerderLenType) ([]byte,error) {
	var bufBinayLen []byte
	switch pLen {
	case PacketHerderLen_16:
		if len(buf)+2 > math.MaxUint16 {
			return nil,errors.New("WriteHerderLen_16 overflow")
		}
		bufBinayLen = make([]byte,2)
		binary.LittleEndian.PutUint16(bufBinayLen,uint16(len(buf)))
	case PacketHerderLen_32:
		if uint32(len(buf)+4) > uint32(math.MaxUint32) {
			return nil,errors.New("WriteHerderLen_32 overflow")
		}
		bufBinayLen = make([]byte,4)
		binary.LittleEndian.PutUint32(bufBinayLen,uint32(len(buf)))
	case PacketHerderLen_64:
		if uint64(len(buf)+8) > uint64(math.MaxUint64) {
			return nil,errors.New("WriteHerderLen_64 overflow")
		}
		bufBinayLen = make([]byte,8)
		binary.LittleEndian.PutUint64(bufBinayLen,uint64(len(buf)))
	}

	return append(bufBinayLen,buf...),nil
}

func Unpack(r io.Reader) (buf []byte,err error) {
	var packHeaderLen = make([]byte,4,4)
	_,err = io.ReadFull(r,packHeaderLen)
	if err != nil {
		return nil, err
	}
	buf = make([]byte,binary.LittleEndian.Uint32(packHeaderLen))
	_,err = io.ReadFull(r,buf)
	return buf,err
}

func UnpackN(r io.Reader,pLen PacketHerderLenType) (buf []byte,err error) {
	var packHeaderLen uint
	switch pLen {
	case PacketHerderLen_16:
		lenBuf,err := read(r,2)
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint(binary.LittleEndian.Uint16(lenBuf))
	case PacketHerderLen_32:
		lenBuf,err := read(r,4)
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint(binary.LittleEndian.Uint32(lenBuf))
	case PacketHerderLen_64:
		lenBuf,err := read(r,8)
		if err != nil {
			return nil, err
		}
		packHeaderLen = uint(binary.LittleEndian.Uint64(lenBuf))
	}

	return read(r,uint32(packHeaderLen))
}

func read(r io.Reader,length uint32) ([]byte,error) {
	var tmp = make([]byte,length,length)
	_,err := io.ReadFull(r,tmp)
	return tmp,err
}