package bondmachine

import (
	"fmt"
	"procbuilder"
	"simbox"
	"strconv"
)

type VM struct {
	Bmach                 *Bondmachine
	Processors            []*procbuilder.VM
	Inputs_regs           []interface{}
	Outputs_regs          []interface{}
	Internal_inputs_regs  []interface{}
	Internal_outputs_regs []interface{}

	send_chans   []chan int
	result_chans []chan string
	recv_chan    chan int

	wait_proc int

	abs_tick uint64
}

func (vm *VM) CopyState(vmsource *VM) {
	for i, pstate := range vmsource.Processors {
		vm.Processors[i].CopyState(pstate)
	}
	// TODO Finish
}

type Sim_config struct {
	Show_ticks   bool
	Show_io_pre  bool
	Show_io_post bool
}

// Simbox rules are converted in a sim drive when the simulation starts and applied during the simulation
type Sim_tick_set map[int]interface{}
type Sim_drive struct {
	Injectables []*interface{}
	AbsSet      map[uint64]Sim_tick_set
	PerSet      map[uint64]Sim_tick_set
}

// This is initializated when the simulation starts and filled on the way
type Sim_tick_get map[int]interface{}
type Sim_tick_show map[int]bool
type Sim_report struct {
	Reportables []*interface{}
	Showables   []*interface{}
	AbsGet      map[uint64]Sim_tick_get
	PerGet      map[uint64]Sim_tick_get
	AbsShow     map[uint64]Sim_tick_show
	PerShow     map[uint64]Sim_tick_show
}

func (vm *VM) Processor_execute(psc *procbuilder.Sim_config, instruct <-chan int, resp chan<- int, result_chan chan<- string, proc_id int) {
	for {
		switch <-instruct {
		case 0:
			resp <- proc_id
			break
		case 1:
			result, err := vm.Processors[proc_id].Step(psc)
			resp <- proc_id
			if err == nil {
				result_chan <- result
			} else {
				result_chan <- ""
			}
		}
	}
}

func (vm *VM) Init() error {
	vm.Processors = make([]*procbuilder.VM, len(vm.Bmach.Processors))
	vm.Inputs_regs = make([]interface{}, vm.Bmach.Inputs)
	vm.Outputs_regs = make([]interface{}, vm.Bmach.Outputs)
	vm.Internal_inputs_regs = make([]interface{}, len(vm.Bmach.Internal_inputs))
	vm.Internal_outputs_regs = make([]interface{}, len(vm.Bmach.Internal_outputs))
	vm.abs_tick = uint64(0)

	for i, proc_dom_id := range vm.Bmach.Processors {
		pvm := new(procbuilder.VM)
		pvm.Mach = vm.Bmach.Domains[proc_dom_id]
		pvm.Init()

		vm.Processors[i] = pvm
	}

	vm.send_chans = make([]chan int, len(vm.Bmach.Processors))
	vm.result_chans = make([]chan string, len(vm.Bmach.Processors))
	vm.recv_chan = make(chan int)

	vm.wait_proc = 0

	for i := 0; i < len(vm.Processors); i++ {
		vm.wait_proc = vm.wait_proc + 1
		vm.send_chans[i] = make(chan int)
		vm.result_chans[i] = make(chan string)
	}

	switch vm.Bmach.Rsize {
	case 8:
		for i := 0; i < vm.Bmach.Inputs; i++ {
			vm.Inputs_regs[i] = uint8(0)
		}
		for i := 0; i < vm.Bmach.Outputs; i++ {
			vm.Outputs_regs[i] = uint8(0)
		}
		for i := 0; i < len(vm.Bmach.Internal_inputs); i++ {
			vm.Internal_inputs_regs[i] = uint8(0)
		}
		for i := 0; i < len(vm.Bmach.Internal_outputs); i++ {
			vm.Internal_outputs_regs[i] = uint8(0)
		}
	case 16:
		for i := 0; i < vm.Bmach.Inputs; i++ {
			vm.Inputs_regs[i] = uint16(0)
		}
		for i := 0; i < vm.Bmach.Outputs; i++ {
			vm.Outputs_regs[i] = uint16(0)
		}
		for i := 0; i < len(vm.Bmach.Internal_inputs); i++ {
			vm.Internal_inputs_regs[i] = uint16(0)
		}
		for i := 0; i < len(vm.Bmach.Internal_outputs); i++ {
			vm.Internal_outputs_regs[i] = uint16(0)
		}
	default:
		// TODO Fix
	}
	//	// Set the initial state of the internal outputs registers
	//	for i, bond := range vm.Bmach.Internal_outputs {
	//		switch bond.Map_to {
	//		case 0:
	//			vm.Internal_outputs_regs[i] = vm.Inputs_regs[bond.Res_id]
	//		case 3:
	//			vm.Internal_outputs_regs[i] = vm.Processors[bond.Res_id].Outputs[bond.Ext_id]
	//		}
	//	}

	return nil
}

func (vm *VM) Launch_processors(s *simbox.Simbox) error {
	for i := 0; i < len(vm.Processors); i++ {

		psc := new(procbuilder.Sim_config)
		pscerr := psc.Init(s, vm.Processors[i])
		check(pscerr)

		go vm.Processor_execute(psc, vm.send_chans[i], vm.recv_chan, vm.result_chans[i], i)
	}
	return nil
}

func (vm *VM) Step(sc *Sim_config) (string, error) {

	result := ""

	if sc != nil {
		if sc.Show_ticks {
			result += "Absolute tick:" + strconv.Itoa(int(vm.abs_tick)) + "\n"
		}
	}

	// Set the internal outputs registers
	for i, bond := range vm.Bmach.Internal_outputs {
		switch bond.Map_to {
		case 0:
			vm.Internal_outputs_regs[i] = vm.Inputs_regs[bond.Res_id]
			//		case 3:
			//			vm.Internal_outputs_regs[i] = vm.Processors[bond.Res_id].Outputs[bond.Ext_id]
		}
	}

	// Transfer to the internal inputs registers the previous outputs according the links
	for i, j := range vm.Bmach.Links {
		if j != -1 {
			vm.Internal_inputs_regs[i] = vm.Internal_outputs_regs[j]
		}
	}

	// Transfer internal inputs registers to their destination
	for i, bond := range vm.Bmach.Internal_inputs {
		switch bond.Map_to {
		//		case 1:
		//			vm.Outputs_regs[bond.Res_id] = vm.Internal_inputs_regs[i]
		case 2:
			vm.Processors[bond.Res_id].Inputs[bond.Ext_id] = vm.Internal_inputs_regs[i]
		}
	}

	if sc != nil {
		if sc.Show_io_pre {
			result += "\tPre-compute IO: " + vm.Dump_io() + "\n"
		}
	}

	// Order the step to processors
	for i := 0; i < len(vm.Processors); i++ {
		vm.send_chans[i] <- 1
		vm.wait_proc = vm.wait_proc - 1
	}

	for {
		i := <-vm.recv_chan
		proc_result := <-vm.result_chans[i]
		if proc_result != "" {
			result += "\tProc: " + strconv.Itoa(i) + "\n"
			result += proc_result
		}
		vm.wait_proc = vm.wait_proc + 1
		if vm.wait_proc == len(vm.Processors) {
			break
		}
	}

	if sc != nil {
		if sc.Show_io_post {
			result += "\tPost-compute IO: " + vm.Dump_io() + "\n"
		}
	}

	// Set the internal outputs registers
	for i, bond := range vm.Bmach.Internal_outputs {
		switch bond.Map_to {
		//		case 0:
		//			vm.Internal_outputs_regs[i] = vm.Inputs_regs[bond.Res_id]
		case 3:
			vm.Internal_outputs_regs[i] = vm.Processors[bond.Res_id].Outputs[bond.Ext_id]
		}
	}

	// Transfer to the internal inputs registers the previous outputs according the links
	for i, j := range vm.Bmach.Links {
		if j != -1 {
			vm.Internal_inputs_regs[i] = vm.Internal_outputs_regs[j]
		}
	}

	// Transfer internal inputs registers to their destination
	for i, bond := range vm.Bmach.Internal_inputs {
		switch bond.Map_to {
		case 1:
			vm.Outputs_regs[bond.Res_id] = vm.Internal_inputs_regs[i]
			//		case 2:
			//			vm.Processors[bond.Res_id].Inputs[bond.Ext_id] = vm.Internal_inputs_regs[i]
		}
	}

	vm.abs_tick++

	return result, nil
}

func (vm *VM) Dump_io() string {
	result := ""
	for i, reg := range vm.Inputs_regs {
		switch vm.Bmach.Rsize {
		case 8:
			result = result + Get_input_name(i) + ": " + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint8)))) + " "
		case 16:
			result = result + Get_input_name(i) + ": " + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint16)))) + " "
		default:
			// TODO Fix
		}
	}
	for i, reg := range vm.Outputs_regs {
		switch vm.Bmach.Rsize {
		case 8:
			result = result + Get_output_name(i) + ": " + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint8)))) + " "
		case 16:
			result = result + Get_output_name(i) + ": " + zeros_prefix(int(vm.Bmach.Rsize), get_binary(int(reg.(uint16)))) + " "
		default:
			// TODO Fix
		}
	}
	return result
}

func (vm *VM) Get_element_location(mnemonic string) (*interface{}, error) {
	// TODO include others
	if len(mnemonic) > 1 && mnemonic[0] == 'i' {
		if i, err := strconv.Atoi(mnemonic[1:]); err == nil {
			if i < len(vm.Inputs_regs) {
				return &vm.Inputs_regs[i], nil
			}
		}
	}
	if len(mnemonic) > 1 && mnemonic[0] == 'o' {
		if i, err := strconv.Atoi(mnemonic[1:]); err == nil {
			if i < len(vm.Outputs_regs) {
				return &vm.Outputs_regs[i], nil
			}
		}
	}
	return nil, Prerror{mnemonic + " unknown"}
}

func (sc *Sim_config) Init(s *simbox.Simbox, vm *VM, conf *Config) error {

	if s != nil {

		for _, rule := range s.Rules {
			if conf.Debug {
				fmt.Println("Loading simbox rule:", rule)
			}
			// Intercept the set rules
			if rule.Timec == simbox.TIMEC_NONE && rule.Action == simbox.ACTION_CONFIG {
				switch rule.Object {
				case "show_ticks":
					sc.Show_ticks = true
				case "show_io_pre":
					sc.Show_io_pre = true
				case "show_io_post":
					sc.Show_io_post = true
				}
			}
		}
	}
	return nil
}

func (sd *Sim_drive) Init(s *simbox.Simbox, vm *VM) error {

	inj := make([]*interface{}, 0)
	absset := make(map[uint64]Sim_tick_set)
	perset := make(map[uint64]Sim_tick_set)

	for _, rule := range s.Rules {
		// Intercept the set rules
		if rule.Timec == simbox.TIMEC_ABS && rule.Action == simbox.ACTION_SET {
			if loc, err := vm.Get_element_location(rule.Object); err == nil {
				if val, err := strconv.Atoi(rule.Extra); err == nil {
					ipos := -1
					for i, iloc := range inj {
						if iloc == loc {
							ipos = i
							break
						}
					}
					if ipos == -1 {
						ipos = len(inj)
						inj = append(inj, loc)
					}

					if act_on_tick, ok := absset[rule.Tick]; ok {
						switch vm.Bmach.Rsize {
						case 8:
							act_on_tick[ipos] = uint8(val)
						case 16:
							act_on_tick[ipos] = uint16(val)
						default:
							// TODO Fix
						}
					} else {
						act_on_tick := make(map[int]interface{})
						switch vm.Bmach.Rsize {
						case 8:
							act_on_tick[ipos] = uint8(val)
						case 16:
							act_on_tick[ipos] = uint16(val)
						default:
							// TODO Fix
						}
						absset[rule.Tick] = act_on_tick
					}
				} else {
					return err
				}
			} else {
				return err
			}
		}
		// Intercept the periodic set rules
		if rule.Timec == simbox.TIMEC_REL && rule.Action == simbox.ACTION_SET {
			if loc, err := vm.Get_element_location(rule.Object); err == nil {
				if val, err := strconv.Atoi(rule.Extra); err == nil {
					ipos := -1
					for i, iloc := range inj {
						if iloc == loc {
							ipos = i
							break
						}
					}
					if ipos == -1 {
						ipos = len(inj)
						inj = append(inj, loc)
					}

					if act_on_tick, ok := perset[rule.Tick]; ok {
						switch vm.Bmach.Rsize {
						case 8:
							act_on_tick[ipos] = uint8(val)
						case 16:
							act_on_tick[ipos] = uint16(val)
						default:
							// TODO Fix
						}
					} else {
						act_on_tick := make(map[int]interface{})
						switch vm.Bmach.Rsize {
						case 8:
							act_on_tick[ipos] = uint8(val)
						case 16:
							act_on_tick[ipos] = uint16(val)
						default:
							// TODO Fix
						}
						perset[rule.Tick] = act_on_tick
					}
				} else {
					return err
				}
			} else {
				return err
			}
		}
	}

	sd.Injectables = inj
	sd.AbsSet = absset
	sd.PerSet = perset
	return nil
}

func (sd *Sim_report) Init(s *simbox.Simbox, vm *VM) error {

	rep := make([]*interface{}, 0)
	sho := make([]*interface{}, 0)
	absget := make(map[uint64]Sim_tick_get)
	perget := make(map[uint64]Sim_tick_get)
	absshow := make(map[uint64]Sim_tick_show)
	pershow := make(map[uint64]Sim_tick_show)

	for _, rule := range s.Rules {
		// Intercept the get rules in absolute time
		if rule.Timec == simbox.TIMEC_ABS && rule.Action == simbox.ACTION_GET {
			if loc, err := vm.Get_element_location(rule.Object); err == nil {
				ipos := -1
				for i, iloc := range rep {
					if iloc == loc {
						ipos = i
						break
					}
				}
				if ipos == -1 {
					ipos = len(rep)
					rep = append(rep, loc)
				}

				if str_on_tick, ok := absget[rule.Tick]; ok {
					switch vm.Bmach.Rsize {
					case 8:
						str_on_tick[ipos] = uint8(0)
					case 16:
						str_on_tick[ipos] = uint16(0)
					default:
						// TODO Fix
					}
				} else {
					str_on_tick := make(map[int]interface{})
					switch vm.Bmach.Rsize {
					case 8:
						str_on_tick[ipos] = uint8(0)
					case 16:
						str_on_tick[ipos] = uint16(0)
					default:
						// TODO Fix
					}
					absget[rule.Tick] = str_on_tick
				}
			} else {
				return err
			}
		}
		// Intercept the get rules in relative time
		if rule.Timec == simbox.TIMEC_REL && rule.Action == simbox.ACTION_GET {
			if loc, err := vm.Get_element_location(rule.Object); err == nil {
				ipos := -1
				for i, iloc := range rep {
					if iloc == loc {
						ipos = i
						break
					}
				}
				if ipos == -1 {
					ipos = len(rep)
					rep = append(rep, loc)
				}

				if str_on_tick, ok := perget[rule.Tick]; ok {
					switch vm.Bmach.Rsize {
					case 8:
						str_on_tick[ipos] = uint8(0)
					case 16:
						str_on_tick[ipos] = uint16(0)
					default:
						// TODO Fix
					}
				} else {
					str_on_tick := make(map[int]interface{})
					switch vm.Bmach.Rsize {
					case 8:
						str_on_tick[ipos] = uint8(0)
					case 16:
						str_on_tick[ipos] = uint16(0)
					default:
						// TODO Fix
					}
					perget[rule.Tick] = str_on_tick
				}
			} else {
				return err
			}
		}
		// Intercept the show rules in absolute time
		if rule.Timec == simbox.TIMEC_ABS && rule.Action == simbox.ACTION_SHOW {
			if loc, err := vm.Get_element_location(rule.Object); err == nil {
				ipos := -1
				for i, iloc := range sho {
					if iloc == loc {
						ipos = i
						break
					}
				}
				if ipos == -1 {
					ipos = len(sho)
					sho = append(sho, loc)
				}

				if str_on_tick, ok := absshow[rule.Tick]; ok {
					str_on_tick[ipos] = true
				} else {
					str_on_tick := make(map[int]bool)
					str_on_tick[ipos] = true
					absshow[rule.Tick] = str_on_tick
				}
			} else {
				return err
			}
		}
		// Intercept the show rules in relative time
		if rule.Timec == simbox.TIMEC_REL && rule.Action == simbox.ACTION_SHOW {
			if loc, err := vm.Get_element_location(rule.Object); err == nil {
				ipos := -1
				for i, iloc := range sho {
					if iloc == loc {
						ipos = i
						break
					}
				}
				if ipos == -1 {
					ipos = len(sho)
					sho = append(sho, loc)
				}

				if str_on_tick, ok := pershow[rule.Tick]; ok {
					str_on_tick[ipos] = true
				} else {
					str_on_tick := make(map[int]bool)
					str_on_tick[ipos] = true
					pershow[rule.Tick] = str_on_tick
				}
			} else {
				return err
			}
		}
	}

	sd.Reportables = rep
	sd.Showables = sho
	sd.AbsGet = absget
	sd.PerGet = perget
	sd.AbsShow = absshow
	sd.PerShow = pershow

	return nil
}
