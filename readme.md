# go-jeans 是一个序列化、反序列化、拆封二进制字节流的包.

**场景**
- 用来解决TCP传输中 [粘包](https://blog.csdn.net/weixin_41047704/article/details/85340311) 的问题
- 高性能的序列化、反序列化数据结构（不使用反射，尽最大可能提高性能）

![golang](https://img.shields.io/badge/golang-v1.19-blue)
![simple](https://img.shields.io/badge/simple-extend-green)
![tcp-Pack](https://img.shields.io/badge/tcp-pack-yellowgreen)
![serve](https://img.shields.io/badge/network_transmission-pack-red)


## 使用教程

* ### 在项目中导入包
  go get -u github.com/Li-giegie/go-jeans

* ### 打包
```go
//使用如下函数 参数需打包的字节
Pack(buf []byte) []byte
//自定义包头长度 参数二可选这16、32、64
PackN(buf []byte, pLen PacketHerderLenType) ([]byte, error)
```

* ### 拆包
```go
//入参一般是connect对象，或是实现了reader的任何对象
Unpack(r io.Reader) (buf []byte, err error)
//参数二 包头长度
UnpackN(r io.Reader, pLen PacketHerderLenType) 
```
[使用例子](./example/tcp-demo/server.go)

* ### 序列化、反序列化使用
[序列化、反序列化使用](./utils_test.go)  
