package automaton

import (
	"context"
	"io"
)

type AutomatonOptions struct {
	Output              io.Writer
	IncludeCalculations bool
}

type Automaton interface {
	Run(ctx context.Context, opts AutomatonOptions) (AutomatonResult, error)
	CurrentCalculationsState() AutomatonCurrentCalculationsState
}

type AutomatonCurrentCalculationsState interface {
	SaveState(w io.Writer) error
}

type AutomatonResult interface {
	SaveResult(w io.Writer) error
}

type State struct {
	Name      string
	Accepting bool
}

type Symbol struct {
	Name string
	// TODO: I think we need more fields here (may be in different kinds of automata?)
}

var BlankSymbol = Symbol{Name: "B"}

type TapeMoveType int

const (
	_ TapeMoveType = iota
	TapeMoveLeft
	TapeMoveRight
)
