package generator

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
)

type Manifest struct {
	Repository string                        `json:"repo"`
	Ref        string                        `json:"ref"`
	Directory  string                        `json:"-"`
	Contracts  map[string]ContractDefinition `json:"contracts"`
}

const manifest_file = "asc.manifest.json"

func NewManifest(path string) (*Manifest, error) {

	b, err := ioutil.ReadFile(path + "/" + manifest_file)
	if err != nil {
		return nil, err
	}

	m := &Manifest{}
	if err := json.Unmarshal(b, m); err != nil {
		return nil, err
	}

	m.Directory = path

	return m, nil
}

func (m *Manifest) GenerateInterface() (*ContractInterface, error) {
	ci := &ContractInterface{
		Repository: m.Repository,
		Ref:        m.Ref,
		Contracts:  map[string]ContractSchema{},
	}

	//Iterate over contracts and generate their schema
	for name, def := range m.Contracts {
		if def.Type == LogicType {
			path := m.Directory + "/" + def.Logic.Path
			tc, err := NewTemplateContract(path)
			if err != nil {
				return ci, err
			}

			enc := base64.StdEncoding.EncodeToString([]byte(tc.Compiled))
			ci.Contracts[name] = ContractSchema{
				Type: def.Type,
				LogicSchema: &LogicSchema{
					Program{
						Bytecode:  enc,
						Address:   tc.Addr,
						Size:      len(tc.Compiled),
						Variables: tc.Variables,
					},
				},
			}
		} else {
			approval_path := m.Directory + "/" + def.App.ApprovalPath.Path
			atc, err := NewTemplateContract(approval_path)
			if err != nil {
				return ci, err
			}
			aenc := base64.StdEncoding.EncodeToString([]byte(atc.Compiled))

			clear_path := m.Directory + "/" + def.App.ClearPath.Path
			ctc, err := NewTemplateContract(clear_path)
			if err != nil {
				return ci, err
			}
			cenc := base64.StdEncoding.EncodeToString([]byte(ctc.Compiled))

			ci.Contracts[name] = ContractSchema{
				Type: def.Type,
				AppSchema: &AppSchema{
					Approval: Program{
						Bytecode:  aenc,
						Address:   atc.Addr,
						Size:      len(atc.Compiled),
						Variables: atc.Variables,
					},
					Clear: Program{
						Bytecode:  cenc,
						Address:   ctc.Addr,
						Size:      len(ctc.Compiled),
						Variables: ctc.Variables,
					},
					GlobalState: def.GlobalState,
					LocalState:  def.LocalState,
				},
			}
		}
	}

	return ci, nil
}
