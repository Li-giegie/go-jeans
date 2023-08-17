package go_jeans

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"unsafe"
)

const intSize = unsafe.Sizeof(int(1))

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

// 打包头的长度：32字节
func Pack(buf []byte) []byte {
	var hl = make([]byte, 4)
	binary.LittleEndian.PutUint32(hl, uint32(len(buf)))
	return append(hl, buf...)
}

// 自定义打包头的长度：16、32、64字节
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

func Unpack(r io.Reader) (buf []byte, err error) {
	var packHeaderLen = make([]byte, 4, 4)
	_, err = io.ReadFull(r, packHeaderLen)
	if err != nil {
		return nil, err
	}
	buf = make([]byte, binary.LittleEndian.Uint32(packHeaderLen))
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
func BaseTypeToBytes(args ...interface{}) ([]byte,error) {
	var buf = new(bytes.Buffer)
	var err error
	for _,arg :=range  args{
		switch v := arg.(type) {
		case string:
			var hl = make([]byte, 4)
			binary.LittleEndian.PutUint32(hl, uint32(len(v)))
			_,err = buf.Write(append(hl,[]byte(v)...))
		case int:
			tmpBuffer := make([]byte,intSize)
			if intSize == 4 {
				binary.BigEndian.PutUint32(tmpBuffer, uint32(v))
			}else {
				binary.BigEndian.PutUint64(tmpBuffer, uint64(v))
			}
			_,err = buf.Write(tmpBuffer)
		case int8:
			err = buf.WriteByte(uint8(v))
		case int16:
			tmpBuffer := make([]byte,2)
			binary.BigEndian.PutUint16(tmpBuffer, uint16(v))
			_,err = buf.Write(tmpBuffer)
		case int32:
			tmpBuffer := make([]byte, 4)
			binary.BigEndian.PutUint32(tmpBuffer, uint32(v))
			_,err = buf.Write(tmpBuffer)
		case int64:
			tmpBuffer := make([]byte, 8)
			binary.BigEndian.PutUint64(tmpBuffer, uint64(v))
			_,err = buf.Write(tmpBuffer)
		case uint:
			tmpBuffer := make([]byte,intSize)
			if intSize == 4 {
				binary.LittleEndian.PutUint32(tmpBuffer,uint32(v))
			}else {
				binary.LittleEndian.PutUint64(tmpBuffer,uint64(v))
			}
			_,err = buf.Write(tmpBuffer)
		case uint8:
			err = buf.WriteByte(v)
		case uint16:
			tmpBuffer := make([]byte,2)
			binary.LittleEndian.PutUint16(tmpBuffer,v)
			_,err = buf.Write(tmpBuffer)
		case uint32:
			tmpBuffer := make([]byte,4)
			binary.LittleEndian.PutUint32(tmpBuffer,v)
			_,err = buf.Write(tmpBuffer)
		case uint64:
			tmpBuffer := make([]byte,8)
			binary.LittleEndian.PutUint64(tmpBuffer,v)
			_,err = buf.Write(tmpBuffer)
		case float32:
			tmpBuffer := make([]byte, 4)
			binary.LittleEndian.PutUint32(tmpBuffer, math.Float32bits(v))
			_,err = buf.Write(tmpBuffer)
		case float64:
			tmpBuffer := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmpBuffer, math.Float64bits(v))
			_,err = buf.Write(tmpBuffer)
		case bool:
			if v {
				buf.WriteByte(1)
			}else {
				buf.WriteByte(0)
			}
			continue
		default:
			return nil,errors.New("conversion of this type is not supported")
		}
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(),nil
}

func BaseTypeToBytesV2(args ...interface{}) ([]byte,error) {
	//var buf = bytes.NewBuffer(nil)
	var bufv2 = make([]byte,0)
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			var hl = make([]byte, 4)
			binary.LittleEndian.PutUint32(hl, uint32(len(v)))

			bufv2 = append(bufv2, hl...)
			bufv2 = append(bufv2, []byte(v)...)
		case int:
			tmpBuffer := make([]byte,intSize)
			if intSize == 4 {
				binary.BigEndian.PutUint32(tmpBuffer, uint32(v))
			}else {
				binary.BigEndian.PutUint64(tmpBuffer, uint64(v))
			}
			bufv2 = append(bufv2, tmpBuffer...)
		case int8:
			bufv2 = append(bufv2, uint8(v))
		case int16:
			tmpBuffer := make([]byte,2)
			binary.BigEndian.PutUint16(tmpBuffer, uint16(v))
			bufv2 = append(bufv2, tmpBuffer...)
		case int32:
			tmpBuffer := make([]byte, 4)
			binary.BigEndian.PutUint32(tmpBuffer, uint32(v))
			bufv2 = append(bufv2, tmpBuffer...)
		case int64:
			tmpBuffer := make([]byte, 8)
			binary.BigEndian.PutUint64(tmpBuffer, uint64(v))
			bufv2 = append(bufv2, tmpBuffer...)
		case uint:
			tmpBuffer := make([]byte,intSize)
			if intSize == 4 {
				binary.LittleEndian.PutUint32(tmpBuffer,uint32(v))
			}else {
				binary.LittleEndian.PutUint64(tmpBuffer,uint64(v))
			}
			bufv2 = append(bufv2, tmpBuffer...)
		case uint8:
			bufv2 = append(bufv2, []byte{v}...)
		case uint16:
			tmpBuffer := make([]byte,2)
			binary.LittleEndian.PutUint16(tmpBuffer,v)
			bufv2 = append(bufv2, tmpBuffer...)
		case uint32:
			tmpBuffer := make([]byte,4)
			binary.LittleEndian.PutUint32(tmpBuffer,v)
			bufv2 = append(bufv2, tmpBuffer...)
		case uint64:
			tmpBuffer := make([]byte,8)
			binary.LittleEndian.PutUint64(tmpBuffer,v)
			bufv2 = append(bufv2, tmpBuffer...)
		case float32:
			tmpBuffer := make([]byte, 4)
			binary.LittleEndian.PutUint32(tmpBuffer, math.Float32bits(v))
			bufv2 = append(bufv2, tmpBuffer...)
		case float64:
			tmpBuffer := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmpBuffer, math.Float64bits(v))
			bufv2 = append(bufv2, tmpBuffer...)
		case bool:
			if v {
				bufv2 = append(bufv2, 1)
			}else {
				bufv2 = append(bufv2, 0)
			}
		default:
			return nil,errors.New("conversion of this type is not supported")
		}
	}


	return bufv2,nil
}


// 将一个字节切片序列化成入参的值，参数要求是GO的基本类型，指针形式传递
func BytesToBaseType(buf []byte,args ...interface{})  error {
	var index int
	for i:=0;i< len(args);i++ {
		switch v := args[i].(type) {
		case *bool:
			if buf[index] == 1{
				*v = true
			}else{
				*v=false
			}
			index++
		case *int:
			if intSize == 4 {
				*v = int(int32(binary.BigEndian.Uint32(buf[index:index+4])))
				index+=4
			}else {
				*v = int(int64(binary.BigEndian.Uint64(buf[index:index+8])))
				index+=8
			}
		case *int8:
			*v = int8(buf[index])
			index++
		case *int16:
			*v = int16(binary.BigEndian.Uint16(buf[index:index+2]))
			index+=2
		case *int32:
			*v = int32(binary.BigEndian.Uint32(buf[index:index+4]))
			index+=4
		case *int64:
			*v = int64(binary.BigEndian.Uint64(buf[index:index+8]))
			index+=8
		case *uint:
			if intSize == 4 {
				*v =uint( binary.LittleEndian.Uint32(buf[index:index+4]))
				index+=4
			}else {
				*v =uint(binary.LittleEndian.Uint64(buf[index:index+8]))
				index+=8
			}
		case *uint8:
			*v = buf[index]
			index+=1
		case *uint16:
			*v = binary.LittleEndian.Uint16(buf[index:index+2])
			index+=2
		case *uint32:
			*v = binary.LittleEndian.Uint32(buf[index:index+4])
			index+=4
		case *uint64:
			*v = binary.LittleEndian.Uint64(buf[index:index+8])
			index+=8
		case *float32:
			*v = math.Float32frombits(binary.LittleEndian.Uint32(buf[index:index+4]))
			index+=4
		case *float64:
			*v = math.Float64frombits(binary.LittleEndian.Uint64(buf[index:index+8]))
			index+=8
		case *string:
			l:=binary.LittleEndian.Uint32(buf[index:index+4])
			*v=string(buf[index+4:int(l)+index+4])
			index+=4+int(l)
		default:
			return errors.New("conversion of this type is not supported")
		}

	}

	return nil
}
