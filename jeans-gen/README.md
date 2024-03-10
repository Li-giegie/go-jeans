# jeans-gen
jeans-gen是一个将结构体成员字段生成go-jeans的Encode、Decode函数的代码生成工具

### 下载安装
``` 
go install github.com/Li-giegie/go-jeans/jeans-gen@latest
```
安装在GOPATH\bin目录 可执行文件名：jeans-gen
### 使用方法
#### 生成结构体成员字段规则
1) .成员字段添加tag为jeans:"true"表示允许jeans-gen生成（可选值true|false）
2) .成员字段为公开
3) .如果成员字段为私有但想被序列化、反序列化，添加tag即可
4) .如果结构体成员字段的类型也是一个结构体，在生成的时候把另一个文件的路径加入在-p参数中
5) .不支持成员类型为自定义类型 (type ui32 uint32)

对于公开字段默认允许生成，对于私有字段不加tag则不允许生成，公开字段不想生成设置tag为false
#### 可执行程序帮助
```
jeans-gen -h    
  -description string
        jeans-gen是一个生成go-jeans编解码函数的工具
  -o string
        保存生成代码的位置,可选值[ "%auto%" | "%append%" | 任意一个路径 ] (default "%auto%")
        取值为 "%auto%" 时自动生成到每个结构体同级目录下 "文件名-gen-jeans"文件中
        取值为 "%append%" 时插入生成代码到go源文件中，不会删除任何代码 
  -p string
        可以是一个文件或者目录，多个文件或者目录用“,”分割 (default "./")
  -s string
        结构体名称,多个名称用","分割,取值为auto时会查找代码中标注“// todo:jeans-gen”的结构体
```
#### 示例：生成main.go文件中的 “User” 结构体
```go
//main.go
type User struct{
    Id int
	Name string
	Age uint8
	Sex bool
	birth Birth
}

type Birth struct {
    year int
	month int
	day int 
}

```
1) .执行命令 jeans-gen -p="main.go" -s="User"
2) .在main.go同级目录内生成main.go-gen-jeans文件里面包含Encode、Decode两个方法如下
3) .并没有在生成的结构中序列化User.birth字段,因为该字段为非公开字段，且没有加tag
```go
func (U *User) Encode() ([]byte, error) {
	return go_jeans.Encode(U.Id,U.Name,U.Age,U.Sex)
}
func (U *User) Decode(data []byte) error {
	return go_jeans.Decode(data,&U.Id,&U.Name,&U.Age,&U.Sex)
}

```
### 注意事项
仅适用与go基本类型

不支持自定义类型

不支持[]struct、map等

不支持一个结构体包含其他包中的类型

工具可能存在BUG，如果造成的不便还请理解，同时提出宝贵的建议

