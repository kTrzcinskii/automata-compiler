package automaton

import (
	"context"
	"errors"
	"io"
)

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

type AutomatonOptions struct {
	Output              io.Writer
	IncludeCalculations bool
}

func (opts AutomatonOptions) validate() error {
	if opts.IncludeCalculations && opts.Output == nil {
		return errors.New("field `Output` must be set when `IncludeCalculations` is enabled")
	}
	return nil
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

func writeCurrentState(a Automaton, w io.Writer) error {
	cs := a.CurrentCalculationsState()
	err := cs.SaveState(w)
	return err
}
