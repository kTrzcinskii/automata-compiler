package automata

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

type TMTransitionKey struct {
	StateName  string
	SymbolName string
}

type TMTransitionValue struct {
	StateName  string
	SymbolName string
	Move       TapeMoveType
}

type TransitionFunction map[TMTransitionKey]TMTransitionValue

type TuringMachine struct {
	States       map[string]State
	Symbols      map[string]Symbol
	CurrentState string
	Tape         []string
	TapeIt       int
	Transitions  TransitionFunction
}

type TuringMachineResult struct {
	FinalState State
	FinalTape  []Symbol
}

type TuringMachineCurrentCalculationsState struct {
	State State
	Tape  []Symbol
	It    int
}

func (tmr TuringMachineResult) SaveResult(w io.Writer) error {
	_, err := w.Write([]byte(fmt.Sprintf("final state: %s, tape: %s\n", tmr.FinalState.Name, tapeString(tmr.FinalTape))))
	return err
}

func tapeString(tape []Symbol) string {
	tapeStr := make([]string, 0, len(tape))
	for _, v := range tape {
		tapeStr = append(tapeStr, v.Name)
	}
	return strings.Join(tapeStr, "|")
}

func (tmc TuringMachineCurrentCalculationsState) SaveState(w io.Writer) error {
	firstLine := fmt.Sprintf("current state: %s, tape: ", tmc.State.Name)
	fpLen := len(firstLine)
	firstLine += tapeString(tmc.Tape) + "\n"
	offset := tmc.It * 2
	secondLine := fmt.Sprintf("%s^\n", strings.Repeat(" ", fpLen+offset))
	out := firstLine + secondLine
	_, err := w.Write([]byte(out))
	return err
}

func (tm *TuringMachine) Run(ctx context.Context, opts AutomataOptions) (AutomataResult, error) {
	err := validateOpts(opts)
	if err != nil {
		panic(err)
	}
	var zero TuringMachineResult
	for {
		select {
		case <-ctx.Done():
			return zero, errors.New("timeout reached")
		default:
			if opts.IncludeCalculations {
				err := tm.writeCurrentState(opts.Output)
				if err != nil {
					return zero, err
				}
			}
			if tm.isInAcceptingState() {
				return tm.FinalState(), nil
			}
			if err := tm.makeTransition(); err != nil {
				return zero, err
			}
		}
	}
}

func validateOpts(opts AutomataOptions) error {
	if opts.Output == nil {
		return errors.New("field `Output` must be set")
	}
	return nil
}

func (tm TuringMachine) CurrentCalculationsState() AutomataCurrentCalculationsState {
	state := tm.States[tm.CurrentState]
	tape := tm.getTape()
	return TuringMachineCurrentCalculationsState{State: state, Tape: tape, It: tm.TapeIt}
}

func (tm TuringMachine) FinalState() AutomataResult {
	finalState := tm.States[tm.CurrentState]
	finalTape := tm.getTape()
	return TuringMachineResult{FinalState: finalState, FinalTape: finalTape}
}

func (tm TuringMachine) isInAcceptingState() bool {
	return tm.States[tm.CurrentState].Accepting
}

func (tm TuringMachine) getTape() []Symbol {
	tape := make([]Symbol, 0, len(tm.Tape))
	for _, v := range tm.Tape {
		tape = append(tape, tm.Symbols[v])
	}
	return tape
}

func (tm *TuringMachine) makeTransition() error {
	key := TMTransitionKey{
		StateName:  tm.CurrentState,
		SymbolName: tm.Tape[tm.TapeIt],
	}
	val, ok := tm.Transitions[key]
	if !ok {
		return fmt.Errorf("cannot continue calucations, missing transition for state %s and symbol %s", key.StateName, key.SymbolName)
	}
	tm.Tape[tm.TapeIt] = val.SymbolName
	tm.CurrentState = val.StateName
	if val.Move == TapeMoveLeft {
		tm.TapeIt--
		if tm.TapeIt < 0 {
			return errors.New("cannot continue calucations, turing machine went out of tape")
		}
	} else {
		tm.TapeIt++
		if tm.TapeIt >= len(tm.Tape) {
			tm.Tape = append(tm.Tape, BlankSymbol.Name)
		}
	}
	return nil
}

func (tm TuringMachine) writeCurrentState(w io.Writer) error {
	cs := tm.CurrentCalculationsState()
	err := cs.SaveState(w)
	return err
}
