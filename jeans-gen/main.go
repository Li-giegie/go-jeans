package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	go_jeans "github.com/Li-giegie/go-jeans"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	x          = "{{x}}"
	structName = "{{structName}}"
	args       = "{{args}}"
)

const EncodeTemplate = `func ({{x}} *{{structName}}) Encode() ([]byte, error) {
	return go_jeans.Encode({{args}})
}`

const DecodeTemplate = `func ({{x}} *{{structName}}) Decode(data []byte) error {
	return go_jeans.Decode(data{{args}})
}`

type structCache struct {
	File string
	*StructInfo
}

func (s *structCache) gen() error {
	f, err := os.OpenFile(s.File+"-gen-jeans", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	var encodeT = EncodeTemplate
	var decodeT = DecodeTemplate
	//替换x
	encodeT = strings.Replace(encodeT, x, s.Name[:1], 1)
	decodeT = strings.Replace(decodeT, x, s.Name[:1], 1)
	//替换结构体名
	encodeT = strings.Replace(encodeT, structName, s.Name, 1)
	decodeT = strings.Replace(decodeT, structName, s.Name, 1)
	//替换参数列表
	var enArgs, deArgs bytes.Buffer
	for _, field := range s.Fields {
		if field.IsPointer {
			enArgs.WriteString(",*" + string(s.Name[0]) + "." + field.Name)
			deArgs.WriteString("," + string(s.Name[0]) + "." + field.Name)
			continue
		}
		enArgs.WriteString("," + string(s.Name[0]) + "." + field.Name)
		deArgs.WriteString(",&" + string(s.Name[0]) + "." + field.Name)
	}
	encodeT = strings.Replace(encodeT, args, enArgs.String()[1:], 1)
	decodeT = strings.Replace(decodeT, args, deArgs.String(), 1)
	_, err = f.WriteString(encodeT)
	_, err = f.WriteString("\n")
	_, err = f.WriteString(decodeT)
	return err
}

var (
	description = flag.String("description", "", "jeans-gen是一个生成go-jeans编解码函数的工具")
	Paths       = flag.String("p", "./", "可以是一个文件或者目录，多个文件或者目录用“,”分割")
	StructList  = flag.String("s", "", "结构体名称,多个名称用\",\"分割,取值为auto时会查找代码中标注“// todo:jeans-gen”的结构体")
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
func init() {
	flag.Parse()
}

func main() {
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
			conf.structCache[stu.Name] = structCache{
				File:       path,
				StructInfo: stu,
			}
		}
	}
	jeansStruct := make([]*structCache, 0, len(conf.structs))
	for _, s := range conf.structs {
		stu, ok := conf.structCache[s]
		if !ok {
			log.Fatalf("Struct \"%v\" is not found in the following files\n%v\n", s, strings.Join(conf.paths, "\n"))
		}
		tmpStu, err := statisticalField(CopyStructCache(conf.structCache), go_jeans.SupportList, &stu)
		if err != nil {
			log.Fatalln(err)
		}
		jeansStruct = append(jeansStruct, tmpStu)
	}
	for _, info := range jeansStruct {
		fmt.Println(info.Name, info.File)
		for _, field := range info.Fields {
			fmt.Println("统计结果：", *field)
		}
		if err = info.gen(); err != nil {
			log.Fatalln(err)
		}
	}
	return

}

type StructInfo struct {
	Name   string
	Fields []*FieldInfo
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
		File: _struct.File,
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
