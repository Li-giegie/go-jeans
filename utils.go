package go_jeans

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"io"
)

type A struct {
	Str string
}

func _pack_proto(m protoreflect.ProtoMessage) ([]byte,error) {
	var buf = new(bytes.Buffer)

	var mbuf []byte

	var err error

	mbuf,err = proto.Marshal(m)

	if err != nil {
		return nil,newErr("Marshal proto err:",err)
	}

	if err = binary.Write(buf,binary.LittleEndian,uint32(len(mbuf)));err != nil {
		return nil, newErr("write msg len err :",err)
	}
	_,err = buf.Write(mbuf)

	return buf.Bytes(),nil
}

func Read(conn io.Reader) ([]byte,error) {

	packLen,err := ReadN(conn,4)
	if err != nil {
		return nil, newErr("read data err -1:",err)
	}

	return ReadN(conn,binary.LittleEndian.Uint32(packLen))
}

func ReadN(conn io.Reader,length uint32) ([]byte,error) {
	var tmp = make([]byte,length,length)
	_,err := io.ReadFull(conn,tmp)
	return tmp,err
}

func newErr(textOrErr ...interface{}) error {
	var errBuf = new(bytes.Buffer)
	for _, i := range textOrErr {
		errBuf.WriteString(fmt.Sprint(i))
	}
	return errors.New(errBuf.String())
}
