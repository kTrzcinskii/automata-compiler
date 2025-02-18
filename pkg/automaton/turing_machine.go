package automaton

import (
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

type TMTransitionFunction map[TMTransitionKey]TMTransitionValue

type TuringMachine struct {
	States       map[string]State
	Symbols      map[string]Symbol
	CurrentState string
	Tape         []string
	TapeIt       int
	Transitions  TMTransitionFunction
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

func (tm TuringMachine) currentCalculationsState() AutomatonCurrentCalculationsState {
	state := tm.States[tm.CurrentState]
	tape := tm.getTape()
	return TuringMachineCurrentCalculationsState{State: state, Tape: tape, It: tm.TapeIt}
}

func (tm TuringMachine) calculationsFinished() bool {
	return tm.isInAcceptingState()
}

func (tm TuringMachine) result() AutomatonResult {
	finalState := tm.States[tm.CurrentState]
	finalTape := removeUnnecessaryBlanks(tm.getTape())
	return TuringMachineResult{FinalState: finalState, FinalTape: finalTape}
}

func (tm *TuringMachine) makeMove() error {
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

// removeUnnecessaryBlanks removes blank symbols starting from the end of the tape until there is at most one
// blank symbol in the row at the end
func removeUnnecessaryBlanks(tape []Symbol) []Symbol {
	id := len(tape) - 1
	for id > 0 {
		if tape[id].Name == BlankSymbol.Name && tape[id-1].Name == BlankSymbol.Name {
			id--
		} else {
			break
		}
	}
	out := make([]Symbol, id+1)
	copy(out, tape)
	return out
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

func (tmc TuringMachineCurrentCalculationsState) SaveState(w io.Writer) error {
	firstLine := fmt.Sprintf("current state: %s, tape: ", tmc.State.Name)
	fpLen := len(firstLine)
	firstLine += tapeString(tmc.Tape) + "\n"
	offset := 0
	for i := 0; i < tmc.It; i++ {
		offset += len(tmc.Tape[i].Name) + 1
	}
	secondLine := fmt.Sprintf("%s^\n", strings.Repeat(" ", fpLen+offset))
	out := firstLine + secondLine
	_, err := w.Write([]byte(out))
	return err
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
