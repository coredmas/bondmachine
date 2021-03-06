package procbuilder

import ()

type Opcode interface {
	Op_get_name() string
	Op_get_desc() string
	Op_show_assembler(*Arch) string
	Op_get_instruction_len(*Arch) int
	Op_instruction_verilog_header(*Config, *Arch, string) string
	Op_instruction_verilog_reset(*Arch, string) string
	Op_instruction_verilog_internal_state(*Arch, string) string
	Op_instruction_verilog_default_state(*Arch, string) string
	Op_instruction_verilog_state_machine(*Arch, string) string
	Op_instruction_verilog_footer(*Arch, string) string
	Op_instruction_verilog_extra_modules(*Arch, string) ([]string, []string)
	Op_instruction_verilog_extra_block(*Arch, string, uint8, string, []string) string
	Abstract_Assembler(*Arch, []string) ([]UsageNotify, error)
	Assembler(*Arch, []string) (string, error)
	Disassembler(*Arch, string) (string, error)
	Simulate(*VM, string) error
	Generate(*Arch) string
	Required_shared() (bool, []string)
	Required_modes() (bool, []string)
	Forbidden_modes() (bool, []string)
}

type Sharedel interface {
	Shr_get_name() string
	Shortname() string
	Get_header(*Arch, string, int) string
	Get_params(*Arch, string, int) string
	Get_internal_params(*Arch, string, int) string
}

type Prerror struct {
	string
}

func (e Prerror) Error() string {
	return e.string
}

type ByName []Opcode

func (op ByName) Len() int           { return len(op) }
func (op ByName) Swap(i, j int)      { op[i], op[j] = op[j], op[i] }
func (op ByName) Less(i, j int) bool { return op[i].Op_get_name() < op[j].Op_get_name() }
