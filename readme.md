# go-jeans
go-jeans是一个高效的数据序列化与反序列化库，解决手动封装数据结构所带来的繁琐问题，通过将数据结构快速编码为二进制格式，比主流编码格式更优越的性能表现，使得数据的传输和存储更加高效可靠。
### 特性
- 脱离反射机制
- 资源消耗小
- 堪比手动封装
- 支持Go中的基本数据类型和切片类型

### 下载安装
``` 
go get -u github.com/Li-giegie/go-jeans
```

### 使用方法
[编解码支持列表](#支持列表)  
编码方法：  
基础类型 EncodeBaseFaster(buf []byte, args ...interface{}) ([]byte, error)  
- buf 为存放结果的缓冲区，编码结果为函数最终返回值并非该参数
- args 为要编码的字段，字段类型参考支持列表，

基础类型 EncodeBase(args ...interface{}) ([]byte, error)
- args 为要编码的字段，字段类型参考支持列表，

切片类型 EncodeSlice(slice ...interface{}) ([]byte, error)
- slice 为要编码的字段，字段类型参考支持列表

切片类型 EncodeSliceFaster(buf []byte, slice ...interface{}) ([]byte, error)
- buf 为存放结果的缓冲区，编码结果为函数最终返回值并非该参数
- slice 为要编码的字段，字段类型参考支持列表

混合类型 Encode、EncodeFaster参数为基本类型切片类型

解码方法：

基本类型 DecodeBase(buf []byte, args ...interface{}) error
- buf 为编码的字节切片
- args 为要还原的指针参数类型列表

切片类型 DecodeSlice(buf []byte, slice ...interface{}) error
- buf 为编码的字节切片
- slice 为要还原的指针切片类型参数列表

混合类型 Decode参数为基本指针类型切片指针类型

```go
// 编码（序列化） 将go中的基本值类型进行编码，编码的参数和解码的参数顺序必须一致，只有在传递的类型不支持时会返回错误，其他情况不会，注意这一步并不打包返回的切片
func TestEncode(t *testing.T) {
    var s struct{
        A int
        B string
        C bool
    }

    buf ,err := Encode(s.A,s.B,s.C)
    if err != nil {
        return
    }
    fmt.Println(buf)
}

// 解码 将一个字节切片序列化成入参的值，参数要求是GO的基本类型，指针形式传递，编码的参数和解码的参数顺序必须一致
func TestDecode(t *testing.T) {
    var s struct{
        A int
        B string
        C bool
    }
    buf, _ := Encode(s.A,s.B,s.C)
    err := Decode(buf,&s.A,&s.B,&s.C)
    if err != nil {
        return
    }
    fmt.Println(s)
}

// 快速编码
func TestEncodeFaster(t *testing.T) {
    type base struct {
        I  int
        Ui uint
        Bo bool
        B  byte
        Bs []byte
        S  string
    }
    var encodeBase = new(base)
    //伪造数据
    err := faker.FakeData(encodeBase)
    if err != nil {
        t.Error(err)
        return
    }
    // 1.创建一个缓冲切片，容量我们要预估一下，尽量不要让内存不够多次分配会占用可观的性能
    // 结构体中变量的大概占用大小，处了[]byte、string外其他的内存占用是确定的，[]byte、string处了本身的长度外还包活一个长度字段占4个字节，只有知道长度信息才能够还原
    bufCap := 8 + 8 + 1 + 1 + (4 + 100) + (4 + 100)
    // 切片的长度一定是0，容量是我们预估的,CountLength()函数可以计算容量
    buf := make([]byte, 0, bufCap)
    buf, err = EncodeFaster(buf, encodeBase.I, encodeBase.Ui, encodeBase.B, encodeBase.Bs, encodeBase.Bo, encodeBase.S)
    if err != nil {
        t.Error(err)
        return
    }
    // 解码验证
    decodeBase := new(base)
    if err = Decode(buf, &decodeBase.I, &decodeBase.Ui, &decodeBase.B, &decodeBase.Bs, &decodeBase.Bo, &decodeBase.S); err != nil {
        t.Error(err)
        return
    }
    if !reflect.DeepEqual(encodeBase, decodeBase) {
        t.Error("decode fail")
        return
    }
    fmt.Printf("encodeBase: %v \ndecodeBase: %v\n", encodeBase, decodeBase)
}
```
### 注意事项
编解码支持的类型仅为Go中的基本类型：int、int8 ~ int64、uint、uint8 ~ uint64、bool、string、float32~64、byte、[]byte(比较常用)，后续考录支持更多的类似

[Encode 编码|序列化：](#) 入参必须为Go中的基本类型，如果确认入参全部被支持，可忽略错误，如果出现错误会返回入参的顺序，例如如果s.A不被支持即返回信息中包含index 0

[Decode 解码|反序列化：](#) 入参必须为Go中的<span style="color: pink">指针基本类型</span>，如果不是指针返回错误，编码的顺序和解码的顺序需要保持一致，可参考使用方法中的示例

### TODO
优化代码


### Benchmark测试
数据来自我的联想笔记本电脑（仅供参考）
```
go test -run=none -cpu 1 -benchmem -bench=Benchmark
goos: windows
goarch: amd64
pkg: github.com/Li-giegie/go-jeans
cpu: AMD Ryzen 5 5600H with Radeon Graphics
BenchmarkEncode                         12100116               101.4 ns/op           176 B/op          1 allocs/op
BenchmarkDecode                         32333073                36.67 ns/op            0 B/op          0 allocs/op
BenchmarkEncodeAndDecode                 7679233               149.2 ns/op           128 B/op          1 allocs/op
BenchmarkEncodeFaster                   20996823                52.05 ns/op            0 B/op          0 allocs/op
BenchmarkJsonMarshal                     1348072               994.1 ns/op           448 B/op          2 allocs/op
BenchmarkJsonUnmarshal                    287700              5054 ns/op             408 B/op         18 allocs/op
BenchmarkJsonMarshalAndUnmarshal          236295              5114 ns/op             888 B/op         21 allocs/op
BenchmarkProtoBufMarshal                 5213646               238.8 ns/op            80 B/op          1 allocs/op
BenchmarkProtoBufUnmarshal               4922341               247.2 ns/op            80 B/op          2 allocs/op
BenchmarkProtoBufMarshalAndUnmarshal     1991695               579.8 ns/op           384 B/op          4 allocs/op
BenchmarkMsgPackMarshal                  1504008               803.2 ns/op           496 B/op          4 allocs/op
BenchmarkMsgPackUnmarshal                1000000              1188 ns/op              80 B/op          2 allocs/op
BenchmarkGobEncode                       2202207               459.9 ns/op           304 B/op          0 allocs/op
BenchmarkGobDecode                         52582             22945 ns/op            8672 B/op        253 allocs/op
```

### 支持列表
编码支持字段类型：

EncodeBase、EncodeBaseFaster
- string, int8, uint8, bool, int16, uint16, int32, uint32, float32, int, uint, int64, uint64, float64

EncodeSlice、EncodeSliceFaster
- []uint, []uint8, []uint16, []uint32, []uint64, []int, []int8, []int16, []int32, []int64, []float32, []float64, []bool, []string

Encode、EncodeFaster 为上面两种编码的混合版，参数只要是上面两种类型的都支持

解码支持字段类型为编码类型的指针：

DecodeBase
- *string, *int8, *uint8, *bool, *int16, *uint16, *int32, *uint32, *float32, *int, *uint, *int64, *uint64, *float64

DecodeSlice
- *[]uint, *[]uint8, *[]uint16, *[]uint32, *[]uint64, *[]int, *[]int8, *[]int16, *[]int32, *[]int64, *[]float32, *[]float64, *[]bool, *[]string

Decode 为上面两种解码的混合版