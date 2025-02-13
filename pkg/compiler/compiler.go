package compiler

import (
	"automata-compiler/pkg/automata"
)

type Compiler interface {
	Compile() (automata.Automata, error)
}
