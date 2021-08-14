package generator

import (
	"encoding/base64"
	"io/ioutil"
	"log"
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

	for idx, v := range tc.Variables {
		log.Printf("%+v %+v", idx, v)
		//        pos += 1 # pushbytes opcode byte
		//        var.length = int(self.assembled_bytes[pos]) # Get length byte
		//        pos += 1  # length opcode byte
		//        var.start = pos

		//        if vidx == 0:
		//            var.distance = pos
		//        else:
		//            pre = self.template_vars[vidx-1]
		//            var.distance = pos - (pre.start + pre.length)

		//        pos += var.length

		//        if var.is_integer:
		//            pos += 1 # btoi

		//        pos += 2 #store opcode + slot id byte
	}

}

func (tc *TemplateContract) initTemplatePositions() {
	matches := tmpl_re.FindAllStringIndex(tc.Source, -1)
	log.Printf("%+v", matches)

	for _, m := range matches {
		tc.Variables = append(tc.Variables, NewTemplateVariable(tc.Source, m))
	}
}
