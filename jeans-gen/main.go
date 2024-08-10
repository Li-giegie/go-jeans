package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type structCache struct {
	fileName     string
	structString string
	*StructInfo
}

func (s *structCache) gen() ([]byte, error) {
	var templates = []string{EncodeTemplate, DecodeTemplate, StringTemplate}
	//替换参数列表
	var buf, enArgs, deArgs, strArgs, strFormat bytes.Buffer
	for _, field := range s.Fields {
		strFormat.WriteString(field.Name + ": %v, ")
		shoutName := strings.ToLower(s.Name[:1])
		if field.IsPointer {
			enArgs.WriteString(",*" + shoutName + "." + field.Name)
			deArgs.WriteString("," + shoutName + "." + field.Name)
			strArgs.WriteString(",*" + shoutName + "." + field.Name)
			continue
		}
		enArgs.WriteString("," + shoutName + "." + field.Name)
		strArgs.WriteString("," + shoutName + "." + field.Name)
		deArgs.WriteString(",&" + shoutName + "." + field.Name)
	}
	for i, _ := range templates {
		//替换x
		templates[i] = strings.Replace(templates[i], x, strings.ToLower(s.Name[:1]), 1)
		//替换结构体名
		templates[i] = strings.Replace(templates[i], structName, s.Name, 1)
	}
	templates[0] = strings.Replace(templates[0], args, enArgs.String()[1:], 1)
	templates[1] = strings.Replace(templates[1], args, deArgs.String(), 1)
	templates[2] = strings.Replace(templates[2], args, "\""+s.Name+" {"+strFormat.String()[:strFormat.Len()-2]+"}\""+strArgs.String(), 1)

	buf.WriteString("\n\n//jeans-gen: struct name: " + s.Name + " time: " + time.Now().String() + "\n")
	for _, template := range templates {
		if _, err := buf.WriteString(template + "\n"); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

var (
	description = flag.String("description", "", "jeans-gen是一个生成go-jeans编解码函数的工具")
	Paths       = flag.String("p", "./", "可以是一个文件或者目录，多个文件或者目录用“,”分割")
	StructList  = flag.String("s", "", "结构体名称,多个名称用\",\"分割,取值为auto时会查找代码中标注“// todo:jeans-gen”的结构体")
	out         = flag.String("o", "%auto%", "保存生成代码的位置,可选值[ \"%auto%\" | \"%append%\" | 任意一个路径 ]\n取值为 \"%auto%\" 时自动生成到每个结构体同级目录下 \"文件名-gen-jeans\"文件中\n取值为 \"%append%\" 时插入生成代码到go源文件中，不会删除任何代码")
)

type configure struct {
	paths       []string
	structs     []string
	structCache map[string]structCache
}

func newConf(Paths, StructList *string) (*configure, error) {
	tmpPath := sliceTrimSpace(strings.Split(*Paths, ","))
	c := new(configure)
	c.paths = make([]string, 0, len(tmpPath))
	for _, _path := range tmpPath {
		ok, isDir := PathExist(_path)
		if !ok {
			return nil, errors.New("open path fail " + _path + " not exist or do not have permission")
		}
		if !isDir {
			c.paths = append(c.paths, _path)
			continue
		}
		err := filepath.Walk(_path, func(path string, info fs.FileInfo, err error) error {
			if err == nil && !info.IsDir() && filepath.Ext(path) == ".go" {
				c.paths = append(c.paths, path)
			}
			return err
		})
		if err != nil {
			return nil, err
		}
	}
	c.structs = sliceTrimSpace(strings.Split(*StructList, ","))
	if len(c.structs) == 0 {
		return nil, errors.New("not find struct")
	}
	c.structCache = make(map[string]structCache, len(c.structs))
	return c, nil
}
func main() {
	flag.Parse()
	conf, err := newConf(Paths, StructList)
	if err != nil {
		log.Fatalln(err)
	}
	p := newParser()
	for _, path := range conf.paths {
		buf, err := os.ReadFile(path)
		if err != nil {
			log.Fatalln(err)
		}
		for _, s := range p.FindAllStructString(string(buf), true) {
			stu := p.Parse(s, "jeans")
			stu.PackageName = p.parsePackageName(string(buf))

			conf.structCache[stu.Name] = structCache{
				fileName:     path,
				StructInfo:   stu,
				structString: s,
			}
		}
	}
	//for _, c := range conf.structCache {
	//	fmt.Println("解析到结构体：",c.StructInfo.PackageName)
	//}
	jeansStruct := make([]*structCache, 0, len(conf.structs))
	for _, s := range conf.structs {
		stu, ok := conf.structCache[s]
		if !ok {
			log.Fatalf("Struct \"%v\" is not found in the following files\n%v\n", s, strings.Join(conf.paths, "\n"))
		}
		tmpStu, err := statisticalField(CopyStructCache(conf.structCache), SupportList, &stu)
		if err != nil {
			log.Fatalln(err)
		}
		jeansStruct = append(jeansStruct, tmpStu)
	}
	var fileMap = make(map[string]*os.File)
	defer func() {
		for _, file := range fileMap {
			_ = file.Close()
		}
	}()
	for _, info := range jeansStruct {
		genBuf, err := info.gen()
		if err != nil {
			log.Fatalln(err)
		}
		switch *out {
		case "%append%":
			fBuf, err := os.ReadFile(info.fileName)
			if err != nil {
				log.Fatalln("read file fail ", err)
			}
			n := bytes.Index(fBuf, []byte(info.structString))
			if n == -1 {
				log.Fatalln("find struct fail")
			}
			err = InsertDataAtPosition(info.fileName, int64(n+len(info.structString)), genBuf)
			if err != nil {
				log.Fatalln(err)
			}
		default:
			fname := info.fileName + "-jeans-gen"
			if *out != "%auto%" {
				fname = *out
			}
			f, ok := fileMap[fname]
			if !ok {
				f, err = os.OpenFile(fname, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
				if err != nil {
					log.Fatalln(err)
				}
				fileMap[fname] = f
			}
			_, err = f.Write(genBuf)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
	return
}

type StructInfo struct {
	PackageName string
	Name        string
	Fields      []*FieldInfo
}

type FieldInfo struct {
	Name      string
	Type      string
	Tag       string
	IsAnon    bool
	IsPointer bool
}

type parser struct {
	structRegex *regexp.Regexp
}

func newParser() *parser {
	return &parser{structRegex: regexp.MustCompile(`type\s+(\w+)\s+struct\s*{([^}]+)}`)}
}

func (p *parser) FindAllStructString(src string, all bool, structName ...string) []string {
	reg := p.structRegex.FindAllStringSubmatch(src, -1)
	result := make([]string, 0, len(reg))
	ok := true
	for _, content := range reg {
		if !all {
			ok = false
			name := p.parseStructName(content[0])
			for _, s := range structName {
				if s == name {
					ok = true
					break
				}
			}
		}
		if ok {
			result = append(result, content[0])
		}
	}
	return result
}

func (p *parser) Parse(structString string, tagKey ...string) *StructInfo {
	stru := new(StructInfo)
	stru.Fields = make([]*FieldInfo, 0)
	stru.Name = p.parseStructName(structString)
	stru.Fields = p.parseField(stru.Name, structString, tagKey)
	return stru
}

func (p *parser) parsePackageName(src string) string {
	for _, s := range strings.Split(src, "\n") {
		n := strings.Index(s, "package ")
		if n == -1 {
			continue
		}
		return strings.TrimSpace(s[n+8:])
	}
	return ""
}

func (p *parser) parseStructName(structString string) string {
	typeIndex := strings.Index(structString, "type ")
	structIndex := strings.Index(structString, " struct ")
	if typeIndex == -1 || structIndex == -1 {
		log.Fatalln("struct invalid: not find struct name")
	}
	Name := strings.Trim(structString[typeIndex+5:structIndex], " ")
	if Name == "" {
		log.Fatalln("struct invalid: struct name is \"\"")
	}
	return Name
}

func (p *parser) parseField(structName, structString string, tagKey []string) []*FieldInfo {
	left := strings.Index(structString, "{")
	right := strings.Index(structString, "}")
	if left == -1 || right == -1 {
		log.Fatalf("struct %s invalid\n", structName)
	}
	fieldContent := structString[left+1 : right]
	if fieldContent == "" {
		log.Fatalf("struct %s invalid : not find field %s \n", structName, structString)
	}
	tmpFields := sliceTrimSpace(strings.Split(fieldContent, "\n"))
	fields := make([]*FieldInfo, 0, len(tmpFields))
	for _, field := range tmpFields {
		if field == "" || strings.Contains(field, "//") || strings.Contains(field, "/*") {
			continue
		}
		_field := sliceTrimSpace(strings.Split(field, " "))
		if len(_field) > 3 {
			tmpField := []string{_field[0], _field[1], strings.Join(_field[2:], "")}
			_field = tmpField
		}
		switch len(_field) {
		case 1:
			if !IsFirstLetterUpper(_field[0]) {
				continue
			}
			fields = append(fields, &FieldInfo{
				Name:      strings.TrimLeft(_field[0], "*"),
				Type:      _field[0],
				IsAnon:    true,
				IsPointer: _field[0][0] == '*',
			})
		case 2:
			//type+tag
			if strings.Contains(_field[1], "`") && strings.Contains(_field[1], ":") && strings.Contains(_field[1], "\"") {
				val := parseFieldTag(_field[1], tagKey)
				if !IsFirstLetterUpper(_field[0]) && val != "true" {
					continue
				}
				fields = append(fields, &FieldInfo{
					Name:      strings.TrimLeft(_field[0], "*"),
					Type:      _field[0],
					Tag:       val,
					IsAnon:    true,
					IsPointer: _field[0][0] == '*',
				})
				continue
			}
			if !IsFirstLetterUpper(_field[0]) {
				continue
			}
			//name+type
			fields = append(fields, &FieldInfo{
				Name:      _field[0],
				Type:      _field[1],
				IsPointer: _field[1][0] == '*',
			})
		case 3:
			val := parseFieldTag(_field[2], tagKey)
			if !IsFirstLetterUpper(_field[0]) && val != "true" {
				continue
			}
			fields = append(fields, &FieldInfo{
				Name:      _field[0],
				Type:      _field[1],
				Tag:       val,
				IsPointer: _field[1][0] == '*',
			})
		default:
			log.Fatalf("struct %s invalid : parse field fail %s\n", structName, structString)
		}
	}
	return fields
}

// 解析一个字段的tag值，key可以为多个满足一个即可获取value
func parseFieldTag(str string, key []string) (val string) {
	if len(key) == 0 {
		log.Fatalln("parseFieldTag key length = 0", str)
	}
	tags := strings.Split(strings.ReplaceAll(strings.ReplaceAll(str, "\"", ""), "`", ""), " ")
	var ok bool
	for _, s := range tags {
		kv := strings.Split(s, ":")
		if len(kv) != 2 {
			continue
		}
		ok = false
		for _, s2 := range key {
			if kv[0] == s2 {
				ok = true
				break
			}
		}
		if ok {
			return kv[1]
		}
	}
	return ""
}

func sliceTrimSpace(str []string) []string {
	var newStr = make([]string, 0, len(str))
	for _, s := range str {
		tsf := strings.TrimSpace(s)
		if tsf != "" {
			newStr = append(newStr, tsf)
		}
	}
	return newStr
}

func IsFirstLetterUpper(s string) bool {
	if len(s) == 0 {
		return false // 如果字符串为空，则没有首字母
	}
	switch s[0] {
	case '_':
		return false
	case '*':
		s = s[1:]
	}
	r, _ := utf8.DecodeRuneInString(s) // 解码字符串的第一个rune（字符）
	return unicode.IsUpper(r)          // 判断该rune是否为大写字母
}

func PathExist(_path string) (ok, isDir bool) {
	info, err := os.Stat(_path)
	if err != nil {
		return false, false
	}
	return true, info.IsDir()
}

// 统计一个结构体包含的所有字段信息
func statisticalField(allStruct map[string]structCache, supportType map[string]struct{}, _struct *structCache, isRecursion ...bool) (ret *structCache, err error) {
	ret = &structCache{
		fileName:     _struct.fileName,
		structString: _struct.structString,
		StructInfo: &StructInfo{
			Name:   _struct.Name,
			Fields: make([]*FieldInfo, 0, len(_struct.Fields)),
		},
	}
	ret.Fields = make([]*FieldInfo, 0, len(_struct.Fields))
	for _, field := range _struct.Fields {
		if _, ok := supportType[field.Type]; ok {
			ret.Fields = append(ret.Fields, field)
			continue
		}
		fieldType := strings.TrimLeft(field.Type, "*")
		resStruct, ok := allStruct[fieldType]
		if !ok {
			return nil, fmt.Errorf("field Type \"%v\" not find\n", field.Type)
		}
		for _, info := range resStruct.Fields {
			info.Name = field.Name + "." + info.Name
		}
		tmpRes, err := statisticalField(allStruct, supportType, &resStruct, true)
		if err != nil {
			return nil, err
		}
		for _, info := range tmpRes.Fields {
			ret.Fields = append(ret.Fields, info)
		}
	}
	return ret, nil
}

func CopyStructCache(s map[string]structCache) map[string]structCache {
	buf, err := json.Marshal(s)
	if err != nil {
		log.Fatalln(err)
	}
	ret := make(map[string]structCache)
	if err = json.Unmarshal(buf, &ret); err != nil {
		log.Fatalln("2", err)
	}
	return ret
}

// 获取文件名不包括后缀名
func GetFileName(filename string) string {
	// 使用strings.Split分割文件名和后缀
	extIndex := strings.LastIndex(filename, ".")
	if extIndex != -1 {
		return filename[:extIndex]
	}
	// 如果没有后缀，直接返回文件名
	return filename
}

func InsertDataAtPosition(filename string, position int64, data []byte) error {
	// 打开文件以供读取
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 读取从插入位置到文件末尾的数据
	restData, err := ioutil.ReadAll(io.NewSectionReader(file, position, -1))
	if err != nil {
		return err
	}

	// 关闭文件
	file.Close()

	// 重新打开文件以供写入
	file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// 移动文件指针到插入位置
	_, err = file.Seek(position, 0)
	if err != nil {
		return err
	}

	// 写入要插入的数据
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	// 将剩余的数据写回文件
	_, err = file.Write(restData)
	if err != nil {
		return err
	}

	return nil
}
