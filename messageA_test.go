package go_jeans

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"testing"
)
var w sync.WaitGroup

func Test_MSGA(t *testing.T) {
	msg := NewMsgA([]byte("hello word~"))
	buf,err := msg.Marshal()
	if err != nil {
		log.Fatalln(err)
	}

	var _msg = new(MessageA)
	_msg,err  = _msg.Unmarshal(buf)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(_msg.Msg))
}

func Test_MSGA_Proto(t *testing.T) {
	msg := NewMsgA_Proto([]byte("hello word~"))
	buf,err := msg.Marshal()
	if err != nil {
		log.Fatalln(err)
	}

	var _msg = new(MessageA)

	_msg,err  = _msg.Unmarshal(buf)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(_msg.Msg))
}

func _msgA_Marshal(str string) *bytes.Buffer {

	msg := NewMsgA([]byte(str))
	buf,err := msg.Marshal()
	if err != nil {
		log.Fatalln("Marshal err ",err)
	}

	return buf
}

func _msgA_Unmarshal(buf *bytes.Buffer) *MessageA{

	msgb,err := (&MessageA{}).Unmarshal(buf)
	if err != nil {
		log.Fatalln("Marshal err ",err)
	}
	return msgb
}

func Test_msgA(t *testing.T) {
	fmt.Println(_msgA_Marshal("hello word"))
	fmt.Println(_msgA_Unmarshal(_msgA_Marshal("hello word")))
}

func Test_msgA_(t *testing.T) {
	var cc = make([]uint32,0)
	for i := 0; i < 1000000; i++ {

		w.Add(1)
		go func() {
			defer w.Done()

			cc = append(cc,NewMsgA([]byte("hello word")).MsgId )
		}()
	}
	w.Wait()

	for f, u := range cc {

		for g, u2 := range cc {
			if u == u2 {
				if f == g { continue }
				log.Fatalln("重复出现先重复值",f,g,u,u2,cc[f],cc[g],cc[f-1])
			}
		}
	}

}

func BenchmarkA(b *testing.B) {

	//log.Fatalln(NewMsgA_Proto([]byte("hello word!")).Marshal())
// goos: windows
// goarch: amd64
// pkg: github.com/Li-giegie/go_jeans
// cpu: 11th Gen Intel(R) Core(TM) i5-11400 @ 2.60GHz

	for i := 0; i < b.N; i++ {
		// { _msgA_Marshal
		// BenchmarkName-12    	 6367321	       178.1 ns/op
		// BenchmarkName-12    	 6736392	       177.8 ns/op
		// BenchmarkName-12    	 6727437	       185.0 ns/op
		// }


		_msgA_Unmarshal(_msgA_Marshal("hello word"))
		// { _msgA_Unmarshal
		// BenchmarkName-12    	10385392	       111.7 ns/op
		// BenchmarkName-12    	10385392	       111.7 ns/op
		// BenchmarkName-12    	10828760	       112.0 ns/op
		// }

		// { _msgA_Marshal \ _msgA_Unmarshal
		//	BenchmarkName-12    	 4723406	       253.7 ns/op
		//	BenchmarkName-12    	 4115900	       267.6 ns/op
		//	BenchmarkName-12    	 4559139	       253.0 ns/op
		// }


		//if _,err := NewMsgA_Proto([]byte("hello word!")).Marshal();err != nil {
		//	log.Fatalln(err)
		//}
		// BenchmarkName-12    	 5014975	       234.7 ns/op
		// BenchmarkName-12    	 4950943	       232.1 ns/op
		// BenchmarkName-12    	 4970926	       235.0 ns/op

		//NewMsgA_Proto([]byte("hello word!")).Marshal()
		// BenchmarkName-12    	 5236460	       233.5 ns/op
		// BenchmarkName-12    	 5132373	       233.9 ns/op
		// BenchmarkName-12    	 5176333	       234.1 ns/op

		//if err := proto.Unmarshal([]byte{8,1,18,11,104,101,108,108,111,32,119,111,114,100,33},&MessageA_Proto{}); err != nil {
		//	log.Fatalln(err)
		//}
		// BenchmarkName-12    	 7274049	       164.0 ns/op
		// BenchmarkName-12    	 7486509	       162.3 ns/op
		// BenchmarkName-12    	 7452726	       171.2 ns/op
		// BenchmarkName-12    	 7413456	       162.3 ns/op

		//proto.Unmarshal([]byte{8,1,18,11,104,101,108,108,111,32,119,111,114,100,33},&MessageA_Proto{})
		// BenchmarkName-12    	 7469208	       165.0 ns/op
		// BenchmarkName-12    	 7471105	       160.4 ns/op
		// BenchmarkName-12    	 7453050	       162.3 ns/op
	}
}