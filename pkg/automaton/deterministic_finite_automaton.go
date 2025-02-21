package automaton

import (
	"fmt"
	"io"
	"strings"
)

type DFATransitionKey struct {
	StateName  string
	SymbolName string
}

type DFATransitionValue struct {
	StateName string
}

type DFATransitionFunction map[DFATransitionKey]DFATransitionValue

type DeterministicFiniteAutomaton struct {
	States       map[string]State
	Symbols      map[string]Symbol
	CurrentState string
	Input        []string
	InputIt      int
	Transitions  DFATransitionFunction
}

type DeterministicFiniteAutomatonResult struct {
	FinalState State
}

type DeterministicFiniteAutomatonCurrentCalculationsState struct {
	State     State
	InputLeft []Symbol
}

func (dfa DeterministicFiniteAutomaton) currentCalculationsState() AutomatonCurrentCalculationsState {
	inputLeft := make([]Symbol, 0, len(dfa.Input)-dfa.InputIt)
	for i := dfa.InputIt; i < len(dfa.Input); i++ {
		name := dfa.Input[i]
		inputLeft = append(inputLeft, dfa.Symbols[name])
	}
	return DeterministicFiniteAutomatonCurrentCalculationsState{
		State:     dfa.States[dfa.CurrentState],
		InputLeft: inputLeft,
	}
}

func (dfa DeterministicFiniteAutomaton) calculationsFinished() bool {
	return dfa.InputIt == len(dfa.Input)
}

func (dfa DeterministicFiniteAutomaton) result() AutomatonResult {
	return DeterministicFiniteAutomatonResult{
		FinalState: dfa.States[dfa.CurrentState],
	}
}

func (dfa *DeterministicFiniteAutomaton) makeMove() error {
	key := DFATransitionKey{
		StateName:  dfa.CurrentState,
		SymbolName: dfa.Input[dfa.InputIt],
	}
	val, ok := dfa.Transitions[key]
	if !ok {
		return fmt.Errorf("cannot continue calculations, missing transition for state %s and symbol %s", key.StateName, key.SymbolName)
	}
	dfa.CurrentState = val.StateName
	dfa.InputIt++
	return nil
}

func (dfa DeterministicFiniteAutomatonCurrentCalculationsState) SaveState(w io.Writer) error {
	inputStr := make([]string, 0, len(dfa.InputLeft))
	for _, v := range dfa.InputLeft {
		inputStr = append(inputStr, v.Name)
	}
	_, err := w.Write([]byte(fmt.Sprintf("current state: %s, input left: %s\n", dfa.State.Name, strings.Join(inputStr, "|"))))
	return err
}

func (dfa DeterministicFiniteAutomatonResult) SaveResult(w io.Writer) error {
	_, err := w.Write([]byte(fmt.Sprintf("final state: %s, accepted: %t\n", dfa.FinalState.Name, dfa.FinalState.Accepting)))
	return err
}
