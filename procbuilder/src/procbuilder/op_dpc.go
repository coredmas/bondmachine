package procbuilder

import (
	"strconv"
)

// The Dpc opcode is both a basic instruction and a template for other instructions.
type Dpc struct{}

func (op Dpc) Op_get_name() string {
	// TODO
	return "dpc"
}

func (op Dpc) Op_get_desc() string {
	// TODO
	return "No operation"
}

func (op Dpc) Op_show_assembler(arch *Arch) string {
	// TODO
	opbits := arch.Opcodes_bits()
	result := "dpc [" + strconv.Itoa(opbits) + "]	// No operation [" + strconv.Itoa(opbits) + "]\n"
	return result
}

func (op Dpc) Op_get_instruction_len(arch *Arch) int {
	// TODO
	opbits := arch.Opcodes_bits()
	return opbits
}

func (op Dpc) Op_instruction_verilog_header(conf *Config, arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Dpc) Op_instruction_verilog_state_machine(arch *Arch, flavor string) string {
	// TODO
	result := ""
	result += "				DPC: begin\n"
	result += "					$display(\"DPC\");\n"
	result += "					_pc <= _pc + 1'b1 ;\n"
	result += "				end\n"

	return result
}

func (op Dpc) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Dpc) Assembler(arch *Arch, words []string) (string, error) {
	// TODO
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()
	result := ""
	for i := opbits; i < rom_word; i++ {
		result += "0"
	}
	return result, nil
}

func (op Dpc) Disassembler(arch *Arch, instr string) (string, error) {
	// TODO
	return "", nil
}

// The simulation does nothing
func (op Dpc) Simulate(vm *VM, instr string) error {
	// TODO
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op Dpc) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Dpc) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op Dpc) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Dpc) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Dpc) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Dpc) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Dpc) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Dpc) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Dpc) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Dpc) Abstract_Assembler(arch *Arch, words []string) ([]UsageNotify, error) {
	result := make([]UsageNotify, 0)
	return result, nil
}

func (Op Dpc) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
