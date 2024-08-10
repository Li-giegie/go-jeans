package main

var SupportList = map[string]struct{}{
	"int":        {},
	"int8":       {},
	"int16":      {},
	"int32":      {},
	"int64":      {},
	"uint":       {},
	"uint8":      {},
	"uint16":     {},
	"uint32":     {},
	"uint64":     {},
	"byte":       {},
	"float32":    {},
	"float64":    {},
	"string":     {},
	"bool":       {},
	"*int":       {},
	"*int8":      {},
	"*int16":     {},
	"*int32":     {},
	"*int64":     {},
	"*uint":      {},
	"*uint8":     {},
	"*uint16":    {},
	"*uint32":    {},
	"*uint64":    {},
	"*byte":      {},
	"*float32":   {},
	"*float64":   {},
	"*string":    {},
	"*bool":      {},
	"[]int":      {},
	"[]int8":     {},
	"[]int16":    {},
	"[]int32":    {},
	"[]int64":    {},
	"[]uint":     {},
	"[]uint8":    {},
	"[]uint16":   {},
	"[]uint32":   {},
	"[]uint64":   {},
	"[]byte":     {},
	"[]float32":  {},
	"[]float64":  {},
	"[]string":   {},
	"[]bool":     {},
	"*[]int":     {},
	"*[]int8":    {},
	"*[]int16":   {},
	"*[]int32":   {},
	"*[]int64":   {},
	"*[]uint":    {},
	"*[]uint8":   {},
	"*[]uint16":  {},
	"*[]uint32":  {},
	"*[]uint64":  {},
	"*[]byte":    {},
	"*[]float32": {},
	"*[]float64": {},
	"*[]string":  {},
	"*[]bool":    {},
}

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

const StringTemplate = `func ({{x}} *{{structName}}) String() string {
	return fmt.Sprintf({{args}})
}`
