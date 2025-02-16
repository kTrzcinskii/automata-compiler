package automata

import (
	"context"
	"io"
)

type AutomataOptions struct {
	Output              io.Writer
	IncludeCalculations bool
}

type Automata interface {
	Run(ctx context.Context, opts AutomataOptions) (AutomataResult, error)
	CurrentCalculationsState() AutomataCurrentCalculationsState
}

type AutomataCurrentCalculationsState interface {
	SaveState(w io.Writer) error
}

type AutomataResult interface {
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
