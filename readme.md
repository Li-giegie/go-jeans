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
* [MessageA](#消息类型) \ [MessageA_proto](#消息类型) 由消息ID、消息组成 适用范围（个人见解）：客户端、服务端简单交互
* [MessageB](#消息类型) \ [MessageB_proto](#消息类型) 由消息ID、消息、源地址、目的api、目的地址 组成 适用范围（个人见解）：客户端和服务端、客户端请求服务端转发到指定客户端。
* [MessageC](#消息类型) \ [MessageC_proto](#消息类型) 与MessageB不同的是消息组成部分使用了非int类型的string类型，为什么这样做呢，在序列化和反序列化中，uint32类型占已知4个字节，比起string不确定长度不是更快速吗，更有优势？答案是速度的确如此但不能保证多用户之间的消息ID不重复这是完全使用消息ID区分的情况，如果想一些其他的办法当然也能解决ID重复问题，比如客户端之间隔离...，但想要保证每一个用户的消息的ID都不相同使用int就显得力不从心了，如果是大并发情况下产生一个重复的ID是一件非常严重和棘手的事情。

## 使用教程
* ### 在项目中导入包
  go get -u github.com/Li-giegie/go-jeans

* ### 打包
#### 选择适合的消息类型（MessageA、MessageB、MessageC），每个消息对象都挂在两个方法，Marshal(打包)、Unmarshal(拆包)
```go
//选择一种消息类型
msgA := go_jeans.NewMsgA([]byte("hello ? i'm the client !"))
//打包 返回*bytes.Buffer, error
buf,err := msgA.Marshal()
if err!= nil {
    log.Fatalln("pack err：",err)
}
//buf.Bytes()打包后的字节
```

* ### 拆包
#### 根据打包的消息类型（MessageA、MessageB、MessageC）进行拆包 

```go
//根据发送的消息类型选择接收的类型
//创建了一个指针MessageA，并拆包
msgA,err := new(go_jeans.MessageA).Unmarshal(*conn)
if err != nil {
log.Fatalln("read msg err:",err)
}

log.Println(msgA)

```
