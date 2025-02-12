package compiler

import "automata-compiler/automata"

type Compiler interface {
	Compile() (automata.Automata, error)
}
