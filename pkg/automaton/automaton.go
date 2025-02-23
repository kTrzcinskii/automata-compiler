package automaton

import (
	"context"
	"errors"
	"io"
	"strings"
)

type State struct {
	Name      string
	Accepting bool
}

type Symbol struct {
	Name string
}

type Automaton interface {
	currentCalculationsState() AutomatonCurrentCalculationsState
	calculationsFinished() bool
	result() AutomatonResult
	makeMove() error
}

type AutomatonCurrentCalculationsState interface {
	SaveState(w io.Writer) error
}

type AutomatonResult interface {
	SaveResult(w io.Writer) error
}

func Run(ctx context.Context, a Automaton, opts AutomatonOptions) (AutomatonResult, error) {
	err := opts.validate()
	if err != nil {
		panic(err)
	}
	var zero AutomatonResult
	for {
		select {
		case <-ctx.Done():
			return zero, errors.New("timeout reached")
		default:
			if opts.IncludeCalculations {
				err := writeCurrentState(a, opts.Output)
				if err != nil {
					return zero, err
				}
			}
			if a.calculationsFinished() {
				return a.result(), nil
			}
			if err := a.makeMove(); err != nil {
				return zero, err
			}
		}
	}
}

func writeCurrentState(a Automaton, w io.Writer) error {
	cs := a.currentCalculationsState()
	err := cs.SaveState(w)
	return err
}

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

func symbolsToString(symbols []Symbol) string {
	symbolsStr := make([]string, 0, len(symbols))
	for _, v := range symbols {
		symbolsStr = append(symbolsStr, v.Name)
	}
	return strings.Join(symbolsStr, "|")
}
