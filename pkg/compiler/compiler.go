package compiler

import "automata-compiler/pkg/automaton"

type Compiler interface {
	Compile() (automaton.Automaton, error)
}
