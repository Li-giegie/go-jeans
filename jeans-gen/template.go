package main

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
