# go-jeans 是一个打包套接字，字节流的包，主要用来解决TCP传输中 [粘包](https://blog.csdn.net/weixin_41047704/article/details/85340311) 的问题

![golang](https://img.shields.io/badge/golang-v1.19-blue)
![simple](https://img.shields.io/badge/simple-extend-green)
![tcp-Pack](https://img.shields.io/badge/tcp-pack-yellowgreen)
![serve](https://img.shields.io/badge/network_transmission-pack-red)
## 一个消息包含 发送、回复、接收

* ### go-jeans 消息打包方式
  * 基于 [Protobuff](https://zhuanlan.zhihu.com/p/401958878) 的字节流
  * 基于 封装消息头字节流
* ### 基于三类消息结构传输 ([MessageA](#消息结构)、[MessageB](#消息结构)、[MessageC](#消息结构))

## 本框架封装了常见的三类消息格式 `_proto标识代表使用Protobuff打包格式`
* [MessageA](#消息结构) \ [MessageA_proto](#消息结构) 由消息ID、消息组成 适用范围（个人见解）：客户端、服务端简单交互
* [MessageB](#消息结构) \ [MessageB_proto](#消息结构) 由消息ID、消息、源地址、目的api、目的地址 组成 适用范围（个人见解）：客户端和服务端、客户端请求服务端转发到指定客户端。
* [MessageC](#消息结构) \ [MessageC_proto](#消息结构) 与MessageB不同的是消息组成部分使用了非int类型的string类型，为什么这样做呢，在序列化和反序列化中，uint32类型占已知4个字节，比起string不确定长度不是更快速吗，更有优势？答案是速度的确如此但不能保证多用户之间的消息ID不重复这是完全使用消息ID区分的情况，如果想一些其他的办法当然也能解决ID重复问题，比如客户端之间隔离...，但想要保证每一个用户的消息的ID都不相同使用int就显得力不从心了，如果是大并发情况下产生一个重复的ID是一件非常严重和棘手的事情。


## 不是最新的文档 完整的实例在test文件里面
## 消息结构
#### MessageA 由消息ID、消息组成 适用范围（个人见解）：客户端、服务端简单交互
```go
//由消息ID、消息组成 适用范围（个人见解）：客户端、服务端简单交互
type MessageA struct {
    state         protoimpl.MessageState
    sizeCache     protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields
    
    MsgId uint32 `protobuf:"varint,1,opt,name=MsgId,proto3" json:"MsgId,omitempty"`
    Msg   []byte `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
}
```
#### MessageB 由消息ID、消息、源地址（或标识）、目的接口、目的地址组成 适用范围（个人见解）：客户端和服务端、客户端请求服务端转发到指定客户端。
```go
//由消息ID、消息、源地址（或标识）、目的接口、目的地址组成 适用范围（个人见解）：客户端和服务端、客户端请求服务端转发到指定客户端。
type MessageB struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MsgId    uint32 `protobuf:"varint,1,opt,name=MsgId,proto3" json:"MsgId,omitempty"`
	Msg      []byte `protobuf:"bytes,2,opt,name=Msg,proto3" json:"Msg,omitempty"`
	SrcAddr  uint32 `protobuf:"varint,3,opt,name=SrcAddr,proto3" json:"SrcAddr,omitempty"`
	DestApi  uint32 `protobuf:"varint,4,opt,name=DestApi,proto3" json:"DestApi,omitempty"`
	DestAddr uint32 `protobuf:"varint,5,opt,name=DestAddr,proto3" json:"DestAddr,omitempty"`
}


```
## 使用教程
* ### 在项目中导入包
    go get -u github.com/Li-giegie/go_jeans
* ### 打包
#### 选择适合的消息结构（MessageA、MessageB） 打包提供三种方法
```go
// 1. 创建一个消息
//入参 打包的消息内容
NewMsgA(msg []byte) *MessageA

//打包字符串
NewMsgA_String(msg string) ([]byte,error)

////打包Json对象
NewMsgA_JSON(obj interface{}) ([]byte,error)


// 1.创建一个消息
msgA := NewMsgA_String("hello i'm client !")
// 2.获取字节流 
buf,err := msgA.Bytes()
if err != nil {
	painc(any(err))
}


```

* ### 拆包
#### 根据打包的消息结构（MessageA、MessageB）进行拆包 

```go
//入参一个实现了io.Reader接口的对象 在go中一般情况下为socket的 connect对象
// 拆解MessageA 结构的包
UnpackA(conn io.Reader) (*MessageA,error)

// 拆解MessageA 结构的包
UnpackB(conn io.Reader) (*MessageB,error)

```

* ### 回复消息
#### 请求过来的消息回复 非常有必要这样做 因为需要保证MsgId一致性
```go
// 用法一、
msgA,err := go_jeans.UnpackA(*conn)
msgA.Reply([]byte(""))

```
