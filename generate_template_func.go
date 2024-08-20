package go_jeans

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// GenerateJeansFuncs 生成结构体编解码函数，obj 生成的结构体，只是简化了手动输入结构体成员字段名称这一过程
func GenerateJeansFuncs(w io.Writer, obj interface{}, mode ModeType) error {
	s := new(struc)
	typ := reflect.TypeOf(obj)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("obj type %s not struct", typ.Name())
	}
	s.Name = typ.Name()
	field, err := getStructFields(nil, nil, typ, mode)
	if err != nil {
		return err
	}
	s.Fields = field
	s.genEncodeFunc(w)
	s.genDecodeFunc(w)
	s.genStringFunc(w)
	s.genFieldsFunc(w)
	s.genFieldPointersFunc(w)
	return nil
}

type struc struct {
	Name   string
	Fields []*fieldInfo
}

func (s *struc) genEncodeFunc(w io.Writer) {
	w.Write([]byte(s.gen(encodeTemplate, "*")))
}

func (s *struc) genDecodeFunc(w io.Writer) {
	w.Write([]byte(s.gen(decodeTemplate, "&")))
}

func (s *struc) genFieldsFunc(w io.Writer) {
	w.Write([]byte(s.gen(fieldsTemplate, "*")))
}

func (s *struc) genFieldPointersFunc(w io.Writer) {
	w.Write([]byte(s.gen(fieldPointerTemplate, "&")))
}

func (s *struc) gen(tmpl string, c string) string {
	txt := strings.Replace(tmpl, x_tpl, strings.ToLower(s.Name)[:1], 1)
	txt = strings.Replace(txt, name_tpl, s.Name, 1)
	fields := make([]string, len(s.Fields))
	x := strings.ToLower(s.Name)[:1]
	pty := false
	for i, field := range s.Fields {
		if c == "&" {
			pty = !field.IsPtr
		} else {
			pty = field.IsPtr
		}
		if pty {
			fields[i] = c + x + "." + field.Name
			continue
		}
		fields[i] = x + "." + field.Name
	}
	txt = strings.Replace(txt, fields_tpl, strings.Join(fields, ","), 1)
	return txt
}

func (s *struc) genStringFunc(w io.Writer) {
	txt := s.gen(stringTemplate, "*")
	fieldNames := make([]string, len(s.Fields))
	for i, field := range s.Fields {
		fieldNames[i] = field.Name + ": %v"
	}
	txt = strings.Replace(txt, fieldNames_tpl, strings.Join(fieldNames, ", "), 1)
	w.Write([]byte(txt))
}

type fieldInfo struct {
	Name  string
	Type  string
	Kind  string
	Tag   string
	IsPtr bool
}

func (f *fieldInfo) String() string {
	return fmt.Sprintf("name: %s, Type: %s, Kind: %s, Tag %s, IsPtr: %v", f.Name, f.Type, f.Kind, f.Tag, f.IsPtr)
}

const reflectTagName = "jeans"

type ModeType uint8

const (
	ModeType_All ModeType = iota
	ModeType_TagOrDefaultBehavior
)

func getStructFields(cache map[string][]*fieldInfo, parentField *fieldInfo, obj interface{}, mode ModeType) ([]*fieldInfo, error) {
	typ, ok := obj.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(obj)
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%s nnt is struct type", typ.String())
	}
	if parentField == nil {
		parentField = &fieldInfo{
			Name: "",
			Type: typ.Name(),
			Kind: typ.Kind().String(),
		}
	}
	if cache == nil {
		cache = make(map[string][]*fieldInfo)
	}
	result := make([]*fieldInfo, 0, typ.NumField())
	var isPtr bool
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
			isPtr = true
		} else {
			isPtr = false
		}
		fieldName := selNotEmptyStr(field.Name, fieldType.Name())
		if parentField.Name != "" {
			fieldName = parentField.Name + "." + fieldName
		}
		fieldTypeStr := fieldType.String()
		if mode == ModeType_TagOrDefaultBehavior {
			tag := field.Tag.Get(reflectTagName)
			switch tag {
			case "enable":
			case "disable":
				continue
			default:
				if !field.IsExported() {
					continue
				}
			}
		}
		switch fieldType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.String, reflect.Bool, reflect.Slice:
			result = append(result, &fieldInfo{
				Name:  fieldName,
				Type:  fieldType.Name(),
				Kind:  fieldType.Kind().String(),
				IsPtr: isPtr,
				Tag:   field.Tag.Get(reflectTagName),
			})
		case reflect.Struct:
			// 防止链表结构
			if parentField.Type == fieldType.Name() {
				continue
			}
			// 优先走缓存
			if fields, ok := cache[fieldTypeStr]; ok {
				for _, info := range fields {
					tmp := new(fieldInfo)
					tmp.Type = info.Type
					tmp.Kind = info.Kind
					tmp.Tag = info.Tag
					tmp.IsPtr = info.IsPtr
					if index := strings.Index(info.Name, "."); index == -1 {
						tmp.Name = fieldName + "." + info.Name
					} else {
						tmp.Name = fieldName + info.Name[index:]
					}
					result = append(result, tmp)
				}
				continue
			}
			newParentField := new(fieldInfo)
			newParentField.Name = fieldName
			newParentField.Type = fieldType.Name()
			fields, err := getStructFields(cache, newParentField, fieldType, mode)
			if err != nil {
				return nil, err
			}
			result = append(result, fields...)
			cache[fieldTypeStr] = fields
		default:
			return result, fmt.Errorf("not supper type %v field name %s", fieldType.String(), field.Name)
		}
	}
	return result, nil
}

func selNotEmptyStr(s1, s2 string) string {
	if s1 == "" {
		return s2
	}
	return s1
}

const (
	x_tpl          = "{{x}}"
	name_tpl       = "{{name}}"
	fields_tpl     = "{{fields}}"
	fieldNames_tpl = "{{fieldNames}}"
	encodeTemplate = `
func ({{x}} *{{name}}) Encode() ([]byte,error){
	return jeans.Encode({{fields}})
}
`
	decodeTemplate = `
func ({{x}} *{{name}}) Decode(data []byte) error{
	return jeans.Decode(data,{{fields}})
}
`
	stringTemplate = `
func ({{x}} *{{name}}) String() string{
	return fmt.Sprintf("{{fieldNames}}",{{fields}})
}
`
	fieldsTemplate = `
func ({{x}} *{{name}}) Fields() []interface{} {
	return []interface{}{{{fields}}}
}
`
	fieldPointerTemplate = `
func ({{x}} *{{name}}) FieldPointers() []interface{}  {
	return []interface{}{{{fields}}}
}
`
)
