package go_jeans

import (
	"strconv"
	"testing"
)

func TestNewMessage(t *testing.T) {
	for i := 0; i < 3; i++ {
		//创建一个消息对象
		msg := NewMsg([]byte("hello word - " + strconv.Itoa(i)))
		//序列化成字节流
		buf, _ := msg.Marshal()
		//还原成新的消息对象	new(Message).Unmarshal(buf).debug()
		var msg2 = new(Message)
		msg2.Unmarshal(buf)
		//msg2.debug()
	}
}

//620162409
func BenchmarkNewMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewMsg([]byte("hello word"))
	}
}
