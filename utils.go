package go_jeans

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"unsafe"
)

const intSize = unsafe.Sizeof(int(1))

// 定义用于存放序列化基础类型后的字节切片的缓冲区大小
var BaseTypeToBytesBufferSize = 128

// 包头长度
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

// 自定义消息头长度模式支持2字节、4字节、8字节
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

// 将一个字节切片重新封装成：4个字节长度+buf，的新buf（数据包）
func Pack(buf []byte) []byte {
	var hl = make([]byte, 4)
	binary.LittleEndian.PutUint32(hl, uint32(len(buf)))
	return append(hl, buf...)
}

// 将一个字节切片重新封装成：自定义长度长度（plen（入参））+ +buf，的新buf（数据包）
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

// 入参一个reader，返回一个有由Pack、PackN打包的完整的
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

// 将一个go中的基本类型转成字节切片，参数中包含非基本类型返回空切片和错误，注意这一步并不打包返回的切片
func BaseTypeToBytes(args ...interface{}) ([]byte, error) {
	var buf = make([]byte, 0, BaseTypeToBytesBufferSize)
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			var hl = make([]byte, 4)
			binary.LittleEndian.PutUint32(hl, uint32(len(v)))
			buf = append(buf, hl...)
			buf = append(buf, []byte(v)...)
		case int:
			tmpBuffer := make([]byte, intSize, intSize)
			if intSize == 4 {
				binary.BigEndian.PutUint32(tmpBuffer, uint32(v))
			} else {
				binary.BigEndian.PutUint64(tmpBuffer, uint64(v))
			}
			buf = append(buf, tmpBuffer...)
		case int8:
			buf = append(buf, uint8(v))
		case int16:
			tmpBuffer := make([]byte, 2)
			binary.BigEndian.PutUint16(tmpBuffer, uint16(v))
			buf = append(buf, tmpBuffer...)
		case int32:
			tmpBuffer := make([]byte, 4)
			binary.BigEndian.PutUint32(tmpBuffer, uint32(v))
			buf = append(buf, tmpBuffer...)
		case int64:
			tmpBuffer := make([]byte, 8)
			binary.BigEndian.PutUint64(tmpBuffer, uint64(v))
			buf = append(buf, tmpBuffer...)
		case uint:
			tmpBuffer := make([]byte, intSize)
			if intSize == 4 {
				binary.LittleEndian.PutUint32(tmpBuffer, uint32(v))
			} else {
				binary.LittleEndian.PutUint64(tmpBuffer, uint64(v))
			}
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
		case []byte:
			var hl = make([]byte, 4)
			binary.LittleEndian.PutUint32(hl, uint32(len(v)))
			buf = append(buf, hl...)
			buf = append(buf, v...)
		default:
			return nil, errors.New("conversion of this type is not supported")
		}
	}

	return buf, nil
}

// 将一个字节切片序列化成入参的值，参数要求是GO的基本类型，指针形式传递
func BytesToBaseType(buf []byte, args ...interface{}) error {
	var index int
	for i := 0; i < len(args); i++ {
		switch v := args[i].(type) {
		case *bool:
			if buf[index] == 1 {
				*v = true
			} else {
				*v = false
			}
			index++
		case *int:
			if intSize == 4 {
				*v = int(int32(binary.BigEndian.Uint32(buf[index : index+4])))
				index += 4
			} else {
				*v = int(int64(binary.BigEndian.Uint64(buf[index : index+8])))
				index += 8
			}
		case *int8:
			*v = int8(buf[index])
			index++
		case *int16:
			*v = int16(binary.BigEndian.Uint16(buf[index : index+2]))
			index += 2
		case *int32:
			*v = int32(binary.BigEndian.Uint32(buf[index : index+4]))
			index += 4
		case *int64:
			*v = int64(binary.BigEndian.Uint64(buf[index : index+8]))
			index += 8
		case *uint:
			if intSize == 4 {
				*v = uint(binary.LittleEndian.Uint32(buf[index : index+4]))
				index += 4
			} else {
				*v = uint(binary.LittleEndian.Uint64(buf[index : index+8]))
				index += 8
			}
		case *uint8:
			*v = buf[index]
			index += 1
		case *uint16:
			*v = binary.LittleEndian.Uint16(buf[index : index+2])
			index += 2
		case *uint32:
			*v = binary.LittleEndian.Uint32(buf[index : index+4])
			index += 4
		case *uint64:
			*v = binary.LittleEndian.Uint64(buf[index : index+8])
			index += 8
		case *float32:
			n := binary.LittleEndian.Uint32(buf[index : index+4])
			if float32(n) > math.MaxFloat32 {
				return ErrOfBytesToBaseType_float
			}
			*v = math.Float32frombits(n)
			index += 4
		case *float64:
			n := binary.LittleEndian.Uint64(buf[index : index+8])
			if float64(n) > math.MaxUint64 {
				return ErrOfBytesToBaseType_float
			}
			*v = math.Float64frombits(n)
			index += 8
		case *string:
			l := binary.LittleEndian.Uint32(buf[index : index+4])
			if int(l) > len(buf[index+4:])  {
				return ErrOfBytesToBaseType_String
			}
			*v = string(buf[index+4 : int(l)+index+4])
			index += 4 + int(l)
		case *[]byte:
			l := binary.LittleEndian.Uint32(buf[index : index+4])
			if int(l) > len(buf[index+4:])  {
				return ErrOfBytesToBaseType_SliceBytes
			}
			*v = buf[index+4 : int(l)+index+4]
			index += 4 + int(l)
		default:
			return errors.New("conversion of this type is not supported")
		}

	}

	return nil
}
