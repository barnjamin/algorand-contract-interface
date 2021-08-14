package generator

type ContractType string

const (
	LogicType ContractType = "logicsig"
	AppType   ContractType = "application"
)

// Manifest
type ContractDefinition struct {
	Type  ContractType `json:"type"`
	Logic `json:"logic,omitempty"`
	App   `json:"application,omitempty"`
}

type Logic struct {
	Path string `json:"path"`
}

type App struct {
	ApprovalPath Logic `json:"approval_path"`
	ClearPath    Logic `json:"clear_path"`
	GlobalState  State `json:"global_state_schema"`
	LocalState   State `json:"local_state_schema"`
}

// Interface
type ContractInterface struct {
	Repository string                    `json:"repo"`
	Ref        string                    `json:"ref"`
	Contracts  map[string]ContractSchema `json:"contracts"`
}

type ContractSchema struct {
	Type        ContractType `json:"type"`
	LogicSchema `json:"logic,omitempty"`
	AppSchema   `json:"application,omitempty"`
}

type LogicSchema struct {
	Program `json:"program"`
}

type AppSchema struct {
	Approval    Program `json:"approval_program"`
	Clear       Program `json:"clear_program"`
	GlobalState State   `json:"global_state_schema"`
	LocalState  State   `json:"local_state_schema"`
}

type Program struct {
	Bytecode  string              `json:"bytecode"`  // Base64 encoded compiled bytes of program
	Address   string              `json:"address"`   // Sha512_256 of bytecode
	Size      int                 `json:"size"`      // Number of bytes for compiled program
	Variables []*TemplateVariable `json:"variables"` // Array of Template Vars from contract source
	Source    string              `json:"source"`    // URL Path to the source
}

type State struct {
	Ints  int `json:"num_uints"`
	Bytes int `json:"num_byte_slices"`
}
