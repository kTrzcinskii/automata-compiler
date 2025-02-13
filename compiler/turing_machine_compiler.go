package compiler

import (
	"automata-compiler/automata"
	"automata-compiler/lexer"
	"errors"
	"fmt"
)

type TuringMachineCompiler struct {
	tokens []lexer.Token
	it     int
}

func (tm *TuringMachineCompiler) Compile() (automata.Automata, error) {
	var zero automata.TuringMachine
	states, err := tm.processStates()
	if err != nil {
		return zero, err
	}
	initialState, err := tm.processInitialState(states)
	if err != nil {
		return zero, err
	}
	err = tm.processAcceptingStates(states)
	if err != nil {
		return zero, err
	}
	return automata.TuringMachine{States: states, InitialState: initialState}, nil
}

func NewTuringMachineCompiler(tokens []lexer.Token) *TuringMachineCompiler {
	return &TuringMachineCompiler{tokens: tokens, it: 0}
}

func (tm TuringMachineCompiler) isAtEnd() bool {
	return tm.it >= len(tm.tokens)
}

func (tm *TuringMachineCompiler) advance() lexer.Token {
	if tm.isAtEnd() {
		var t lexer.Token
		return t
	}
	t := tm.tokens[tm.it]
	tm.it++
	return t
}

func (tm *TuringMachineCompiler) processStates() (map[string]automata.State, error) {
	states := make(map[string]automata.State)
	for !tm.isAtEnd() {
		t := tm.advance()
		switch t.Type {
		case lexer.SemicolonToken:
			return states, nil
		case lexer.StateToken:
			name := t.Value
			if _, ok := states[name]; ok {
				return nil, fmt.Errorf("state %s already declared, each state must have unique name", name)
			}
			states[name] = automata.State{Name: name, Accepting: false}
		default:
			return nil, fmt.Errorf("invalid token type, expected: %s or %s, got: %s", lexer.StateToken.String(), lexer.SemicolonToken.String(), t.Type.String())
		}
	}
	return nil, errors.New("missing ';' at the end of states section")
}

func (tm *TuringMachineCompiler) processInitialState(states map[string]automata.State) (string, error) {
	if tm.isAtEnd() {
		return "", fmt.Errorf("missing initial state section")
	}

	initialState := tm.advance()
	if initialState.Type != lexer.StateToken {
		return "", fmt.Errorf("invalid initial state token, expected: %s, got: %s", lexer.StateToken.String(), initialState.Type.String())
	}

	if _, ok := states[initialState.Value]; !ok {
		return "", fmt.Errorf("invalid initial state, state %s was not declared in states list", initialState.Value)
	}

	if tm.isAtEnd() || tm.advance().Type != lexer.SemicolonToken {
		return "", fmt.Errorf("missing ';' after initial state")
	}

	return initialState.Value, nil
}

func (tm *TuringMachineCompiler) processAcceptingStates(states map[string]automata.State) error {
	for !tm.isAtEnd() {
		t := tm.advance()
		switch t.Type {
		case lexer.SemicolonToken:
			return nil
		case lexer.StateToken:
			name := t.Value
			if _, ok := states[name]; !ok {
				return fmt.Errorf("state %s not found, any accepting state must be defined in state list", name)
			}
			states[name] = automata.State{Name: name, Accepting: true}
		default:
			return fmt.Errorf("invalid token type, expected: %s or %s, got: %s", lexer.StateToken.String(), lexer.SemicolonToken.String(), t.Type.String())
		}
	}
	return errors.New("missing ';' at the end of accepting states section")
}
