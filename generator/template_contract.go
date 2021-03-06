package generator

import (
	"encoding/base64"
	"io/ioutil"
	"regexp"
)

type TemplateContract struct {
	Addr              string
	Compiled          []byte
	Source            string
	PlaceHolderSource string
	Variables         []*TemplateVariable
}

//Requires teal to be written with newlines, some short programs are semicolon delimited with spaces so...TODO?
var tmpl_re = regexp.MustCompile(".* TMPL_.*")

func NewTemplateContract(path string) (*TemplateContract, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tc := &TemplateContract{Source: string(b)}

	tc.initTemplatePositions()

	// If there are template variables
	if len(tc.Variables) > 0 {
		relpos := 0
		for _, v := range tc.Variables {
			ph := getPlaceholderType(v.Type)
			pos := v.SourceStart + relpos
			b = append(b[:pos], append([]byte(ph), b[pos+v.SourceLength:]...)...)
			relpos += len(ph) - v.SourceLength
		}
	}
	tc.PlaceHolderSource = string(b)

	result, err := compile(b)
	if err != nil {
		return nil, err
	}

	tc.Addr = result.Hash
	cbytes, err := base64.StdEncoding.DecodeString(result.Result)
	if err != nil {
		return nil, err
	}
	tc.Compiled = cbytes

	if len(tc.Variables) > 0 {
		tc.initCompiledPositions()
	}

	return tc, nil
}

func (tc *TemplateContract) initCompiledPositions() {
	const pushbytesOpcode = 128
	const pushintOpcode = 129

	pos := 1 // Version byte

	// int const block
	if tc.Compiled[pos] == 0x20 {
		size, _, _ := readIntConstBlock(tc.Compiled, pos)
		pos += size
	}

	// byte const block
	if tc.Compiled[pos] == 0x26 {
		size, _, _ := readByteConstBlock(tc.Compiled, pos)
		pos += size
	}

	var found int
	for pc := pos; pc < len(tc.Compiled); {
		opcode := tc.Compiled[pc]
		switch opcode {
		case pushintOpcode:
			size, _, _ := readPushIntOp(tc.Compiled, pc)

			tc.Variables[found].CompiledPosition = pc + 1 // Account for opcode
			tc.Variables[found].CompiledLength = size - 1

			found += 1
			pc += size
		case pushbytesOpcode:
			size, _, _ := readPushByteOp(tc.Compiled, pc)

			tc.Variables[found].CompiledPosition = pc + 2 // Account for opcode and bytesize
			tc.Variables[found].CompiledLength = size - 2

			found += 1
			pc += size
		default:
			// Done with this section, terminate
			pc = len(tc.Compiled)
		}
	}

}

func (tc *TemplateContract) initTemplatePositions() {
	matches := tmpl_re.FindAllStringIndex(tc.Source, -1)
	for _, m := range matches {
		tc.Variables = append(tc.Variables, NewTemplateVariable(tc.Source, m))
	}
}
