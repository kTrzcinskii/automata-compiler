package automata

import (
	"errors"
	"fmt"
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

func (tmr TuringMachineResult) String() string {
	return fmt.Sprintf("state: %s, tape: %s", tmr.FinalState.Name, tmr.tapeString())
}

func (tmr TuringMachineResult) tapeString() string {
	tape := make([]string, 0, len(tmr.FinalTape))
	for _, v := range tmr.FinalTape {
		tape = append(tape, v.Name)
	}
	return strings.Join(tape, "|")
}

func (tm *TuringMachine) Run() (AutomataResult, error) {
	for !tm.isInAcceptingState() {
		if err := tm.makeTransition(); err != nil {
			var zero TuringMachineResult
			return zero, err
		}
	}
	finalState := tm.States[tm.CurrentState]
	finalTape := tm.getFinalTape()
	return TuringMachineResult{FinalState: finalState, FinalTape: finalTape}, nil
}

func (tm TuringMachine) isInAcceptingState() bool {
	return tm.States[tm.CurrentState].Accepting
}

func (tm TuringMachine) getFinalTape() []Symbol {
	finalTape := make([]Symbol, 0, len(tm.Tape))
	for _, v := range tm.Tape {
		finalTape = append(finalTape, tm.Symbols[v])
	}
	return finalTape
}

func (tm *TuringMachine) makeTransition() error {
	key := TMTransitionKey{
		StateName:  tm.CurrentState,
		SymbolName: tm.Tape[tm.TapeIt],
	}
	val, ok := tm.Transitions[key]
	if !ok {
		return fmt.Errorf("missing transition for state %s and symbol %s", key.StateName, key.SymbolName)
	}
	tm.Tape[tm.TapeIt] = val.SymbolName
	tm.CurrentState = val.StateName
	if val.Move == TapeMoveLeft {
		tm.TapeIt--
		if tm.TapeIt < 0 {
			return errors.New("turing machine went out of tape")
		}
	} else {
		tm.TapeIt++
		if tm.TapeIt >= len(tm.Tape) {
			tm.Tape = append(tm.Tape, BlankSymbol.Name)
		}
	}
	return nil
}
