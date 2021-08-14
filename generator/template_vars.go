package generator

import "strings"

type VarType string

const (
	Uint64 VarType = "uint64"
	Addr   VarType = "addr"
	Bytes  VarType = "bytes"
)

type TemplateVariable struct {
	Name string
	Type VarType

	SourceStart  int
	SourceLength int

	CompiledPosition int
	CompiledLength   int
}

func NewTemplateVariable(src string, pos []int) *TemplateVariable {
	//TODO: check that lengths of stuff wont panic on slice idx access

	chunks := strings.Split(src[pos[0]:pos[1]], " ")

	// Get just the TMPL_ section
	start := pos[0] + len(chunks[0]) + 1

	return &TemplateVariable{
		Name:         chunks[1],
		Type:         getVarType(chunks[0]),
		SourceStart:  start,
		SourceLength: pos[1] - start,
	}
}

func getVarType(op string) VarType {
	switch op {
	case "int", "pushint":
		return Uint64
	case "bytes", "pushbytes":
		return Bytes
	default:
		return Addr
	}
}

func getPlaceholderType(t VarType) string {
	switch t {
	case Uint64:
		return "1"
	case Bytes:
		return "b64(ZHVtbXk=)"
	default:
		return "7LQ7U4SEYEVQ7P4KJVCHPJA5NSIFJTGIEXJ4V6MFS4SL5FMDW6MYHL2JXM"
	}
}
