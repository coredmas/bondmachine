package procbuilder

import (
	"strconv"
)

// The Je opcode is both a basic instruction and a template for other instructions.
type Je struct{}

func (op Je) Op_get_name() string {
	// TODO
	return "je"
}

func (op Je) Op_get_desc() string {
	// TODO
	return "No operation"
}

func (op Je) Op_show_assembler(arch *Arch) string {
	// TODO
	opbits := arch.Opcodes_bits()
	result := "je [" + strconv.Itoa(opbits) + "]	// No operation [" + strconv.Itoa(opbits) + "]\n"
	return result
}

func (op Je) Op_get_instruction_len(arch *Arch) int {
	// TODO
	opbits := arch.Opcodes_bits()
	return opbits
}

func (op Je) Op_instruction_verilog_header(conf *Config, arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Je) Op_instruction_verilog_state_machine(arch *Arch, flavor string) string {
	// TODO
	result := ""
	result += "				JE: begin\n"
	result += "					$display(\"JE\");\n"
	result += "					_pc <= _pc + 1'b1 ;\n"
	result += "				end\n"

	return result
}

func (op Je) Op_instruction_verilog_footer(arch *Arch, flavor string) string {
	// TODO
	return ""
}

func (op Je) Assembler(arch *Arch, words []string) (string, error) {
	// TODO
	opbits := arch.Opcodes_bits()
	rom_word := arch.Max_word()
	result := ""
	for i := opbits; i < rom_word; i++ {
		result += "0"
	}
	return result, nil
}

func (op Je) Disassembler(arch *Arch, instr string) (string, error) {
	// TODO
	return "", nil
}

// The simulation does nothing
func (op Je) Simulate(vm *VM, instr string) error {
	// TODO
	vm.Pc = vm.Pc + 1
	return nil
}

// The random genaration does nothing
func (op Je) Generate(arch *Arch) string {
	// TODO
	return ""
}

func (op Je) Required_shared() (bool, []string) {
	// TODO
	return false, []string{}
}

func (op Je) Required_modes() (bool, []string) {
	return false, []string{}
}

func (op Je) Forbidden_modes() (bool, []string) {
	return false, []string{}
}

func (op Je) Op_instruction_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Je) Op_instruction_verilog_reset(arch *Arch, flavor string) string {
	return ""
}

func (Op Je) Op_instruction_verilog_default_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Je) Op_instruction_verilog_internal_state(arch *Arch, flavor string) string {
	return ""
}

func (Op Je) Op_instruction_verilog_extra_modules(arch *Arch, flavor string) ([]string, []string) {
	return []string{}, []string{}
}

func (Op Je) Abstract_Assembler(arch *Arch, words []string) ([]UsageNotify, error) {
	// TODO Partial
	result := make([]UsageNotify, 1)
	newnot := UsageNotify{C_OPCODE, "je", I_NIL}
	result[0] = newnot
	return result, nil
}

func (Op Je) Op_instruction_verilog_extra_block(arch *Arch, flavor string, level uint8, blockname string, objects []string) string {
	result := ""
	switch blockname {
	default:
		result = ""
	}
	return result
}
