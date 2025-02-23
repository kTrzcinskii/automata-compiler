package automaton

import (
	"errors"
	"fmt"
	"io"
)

// Special symbols that must be included in every PA by compiler
var (
	// Input end symbol should be placed after the last element in `Input`
	InputEndSymbol = Symbol{
		Name: "{",
	}
	// Stack start should be placed at the beginning of the `Stack`
	StackStartSymbol = Symbol{
		Name: "}",
	}
)

type PATransitionKey struct {
	StateName       string
	InputSymbolName string
	StackSymbolName string
}

type PATransitionValue struct {
	StateName        string
	StackSymbolNames []string
}

type PATransitionFunction map[PATransitionKey]PATransitionValue

type PushdownAutomaton struct {
	States       map[string]State
	Symbols      map[string]Symbol
	CurrentState string
	Input        []string
	InputIt      int
	Stack        []string
	Transitions  PATransitionFunction
}

type PushdownAutomatonCurrentCalculationsState struct {
	CurrentState State
	Stack        []Symbol
	InputLeft    []Symbol
}

type PushdownAutomatonResult struct {
	FinalState State
	Stack      []Symbol
}

func (pa PushdownAutomaton) currentCalculationsState() AutomatonCurrentCalculationsState {
	stack := pa.getStack()
	input := pa.getInputLeft()
	return PushdownAutomatonCurrentCalculationsState{
		CurrentState: pa.States[pa.CurrentState],
		Stack:        stack,
		InputLeft:    input,
	}
}

func (pa PushdownAutomaton) calculationsFinished() bool {
	return pa.InputIt == len(pa.Input)
}

func (pa PushdownAutomaton) result() AutomatonResult {
	stack := pa.getStack()
	finalState := pa.States[pa.CurrentState]
	return PushdownAutomatonResult{
		FinalState: finalState,
		Stack:      stack,
	}
}

func (pa *PushdownAutomaton) makeMove() error {
	if len(pa.Stack) == 0 {
		return errors.New("stack is empty")
	}
	// Remove last element from stack
	// It's user's responsibility to always have at least one element ('}') on the stack
	stackSybmol := pa.Stack[len(pa.Stack)-1]
	pa.Stack = pa.Stack[:len(pa.Stack)-1]

	input := pa.Input[pa.InputIt]

	key := PATransitionKey{
		StateName:       pa.CurrentState,
		InputSymbolName: input,
		StackSymbolName: stackSybmol,
	}

	value, ok := pa.Transitions[key]
	if !ok {
		return fmt.Errorf("cannot continue calculations, missing transition for state %s, symbol %s and stack symbol %s", pa.CurrentState, input, stackSybmol)
	}

	pa.CurrentState = value.StateName
	pa.InputIt++
	pa.Stack = append(pa.Stack, value.StackSymbolNames...)
	return nil
}

func (pa PushdownAutomaton) getStack() []Symbol {
	stack := make([]Symbol, 0, len(pa.Stack))
	for _, v := range pa.Stack {
		stack = append(stack, pa.Symbols[v])
	}
	return stack
}

func (pa PushdownAutomaton) getInputLeft() []Symbol {
	inputLeft := make([]Symbol, 0, len(pa.Input)-pa.InputIt)
	for i := pa.InputIt; i < len(pa.Input); i++ {
		name := pa.Input[i]
		inputLeft = append(inputLeft, pa.Symbols[name])
	}
	return inputLeft
}

func (pa PushdownAutomatonResult) SaveResult(w io.Writer) error {
	stack := symbolsToString(pa.Stack)
	_, err := w.Write([]byte(fmt.Sprintf("final state: %s, accepted: %t, stack: %s\n", pa.FinalState.Name, pa.FinalState.Accepting, stack)))
	return err
}

func (pa PushdownAutomatonCurrentCalculationsState) SaveState(w io.Writer) error {
	stack := symbolsToString(pa.Stack)
	input := symbolsToString(pa.InputLeft)
	_, err := w.Write([]byte(fmt.Sprintf("current state: %s, input left: %s, stack: %s\n", pa.CurrentState.Name, input, stack)))
	return err
}
