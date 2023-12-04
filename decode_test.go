package go_jeans

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDecodeSlice(t *testing.T) {
	var u32list = []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	var u32list2 = []uint32{2, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	buf, err := EncodeSlice(u32list, u32list2)
	if err != nil {
		t.Error(err)
		return
	}
	var deu32List, deu32List2 []uint32

	err = DecodeSlice(buf, &deu32List, &deu32List2)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(deu32List, deu32List2)
}

// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkEnDecodeSlice-12       10091096               109.7 ns/op
func BenchmarkEnDecodeSlice(b *testing.B) {
	var u32list = []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	var result []uint32
	for i := 0; i < b.N; i++ {
		buf, err := EncodeSlice(u32list)
		if err != nil {
			b.Error(err)
			return
		}
		if err = DecodeSlice(buf, &result); err != nil {
			b.Error(err)
			return
		}
	}
}

// cpu: AMD Ryzen 5 5600H with Radeon Graphics
// BenchmarkJson-12          841549              1455 ns/op
func BenchmarkJson(b *testing.B) {
	var u32list = []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	var result []uint32
	for i := 0; i < b.N; i++ {
		buf, err := json.Marshal(u32list)
		if err != nil {
			b.Error(err)
			return
		}
		if err = json.Unmarshal(buf, &result); err != nil {
			b.Error(err)
			return
		}
	}
}
