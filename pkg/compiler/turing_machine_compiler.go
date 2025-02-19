package compiler

import (
	"automata-compiler/pkg/automaton"
	"automata-compiler/pkg/lexer"
	"errors"
	"fmt"
)

type TuringMachineCompiler struct {
	BaseCompiler
}

func NewTuringMachineCompiler(tokens []lexer.Token) *TuringMachineCompiler {
	return &TuringMachineCompiler{BaseCompiler: newBaseCompiler(tokens)}
}

func (tm *TuringMachineCompiler) Compile() (automaton.Automaton, error) {
	states, err := tm.processStates()
	if err != nil {
		return nil, tm.addLinePrefixForErrPrevToken(err)
	}
	initialState, err := tm.processInitialState(states)
	if err != nil {
		return nil, tm.addLinePrefixForErrPrevToken(err)
	}
	err = tm.processAcceptingStates(states)
	if err != nil {
		return nil, tm.addLinePrefixForErrPrevToken(err)
	}
	symbols, err := tm.processSymbols()
	if err != nil {
		return nil, tm.addLinePrefixForErrPrevToken(err)
	}
	tf, err := tm.processTransitions(states, symbols)
	if err != nil {
		return nil, tm.addLinePrefixForErrPrevToken(err)
	}
	initialTape, err := tm.processTape(symbols)
	if err != nil {
		return nil, tm.addLinePrefixForErrPrevToken(err)
	}
	err = tm.checkForCorrectEndingSequnce()
	if err != nil {
		// It's more lexer error than user provided source,
		// so we don't include line here
		return nil, err
	}
	return &automaton.TuringMachine{States: states, CurrentState: initialState, Symbols: symbols, Transitions: tf, Tape: initialTape, TapeIt: 0}, nil
}

func (tm *TuringMachineCompiler) processStates() (map[string]automaton.State, error) {
	states := make(map[string]automaton.State)
	for !tm.isAtEnd() {
		t := tm.advance()
		switch t.Type {
		case lexer.SemicolonToken:
			if len(states) == 0 {
				return nil, errors.New("there must be at least one state defined")
			}
			return states, nil
		case lexer.StateToken:
			name := t.Value
			if _, ok := states[name]; ok {
				return nil, fmt.Errorf("state %s already declared, each state must have unique name", name)
			}
			states[name] = automaton.State{Name: name, Accepting: false}
		default:
			return nil, fmt.Errorf("invalid token type, expected: %s or %s, got: %s", lexer.StateToken.String(), lexer.SemicolonToken.String(), t.Type.String())
		}
	}
	return nil, errors.New("missing ';' at the end of states section")
}

func (tm *TuringMachineCompiler) processInitialState(states map[string]automaton.State) (string, error) {
	initialState, err := tm.consumeTokenWithType("missing initial state section", lexer.StateToken)
	if err != nil {
		return "", err
	}
	if _, ok := states[initialState.Value]; !ok {
		return "", fmt.Errorf("invalid initial state, state %s was not declared in states list", initialState.Value)
	}
	if tm.isAtEnd() || tm.advance().Type != lexer.SemicolonToken {
		return "", fmt.Errorf("missing ';' after initial state")
	}
	return initialState.Value, nil
}

func (tm *TuringMachineCompiler) processAcceptingStates(states map[string]automaton.State) error {
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
			states[name] = automaton.State{Name: name, Accepting: true}
		default:
			return fmt.Errorf("invalid token type, expected: %s or %s, got: %s", lexer.StateToken.String(), lexer.SemicolonToken.String(), t.Type.String())
		}
	}
	return errors.New("missing ';' at the end of accepting states section")
}

func (tm *TuringMachineCompiler) processSymbols() (map[string]automaton.Symbol, error) {
	symbols := make(map[string]automaton.Symbol)
	symbols[automaton.BlankSymbol.Name] = automaton.BlankSymbol
	for !tm.isAtEnd() {
		t := tm.advance()
		switch t.Type {
		case lexer.SemicolonToken:
			return symbols, nil
		case lexer.SymbolToken:
			name := t.Value
			if _, ok := symbols[name]; ok {
				return nil, fmt.Errorf("symbol %s already declared, each symbol must have unique name", name)
			}
			symbols[name] = automaton.Symbol{Name: t.Value}
		default:
			return nil, fmt.Errorf("invalid token type, expected: %s or %s, got: %s", lexer.SymbolToken.String(), lexer.SemicolonToken.String(), t.Type.String())
		}
	}
	return nil, errors.New("missing ';' at the end of symbols section")
}

func (tm *TuringMachineCompiler) processTransitions(states map[string]automaton.State, symbols map[string]automaton.Symbol) (automaton.TMTransitionFunction, error) {
	tf := make(automaton.TMTransitionFunction)
	for !tm.isAtEnd() {
		t := tm.advance()
		switch t.Type {
		case lexer.SemicolonToken:
			return tf, nil
		case lexer.LeftParenToken:
			err := tm.processSingleTransition(states, symbols, tf)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid token type, expected: %s or %s, got: %s", lexer.LeftParenToken.String(), lexer.SemicolonToken.String(), t.Type.String())
		}
	}
	return nil, errors.New("missing ';' at the end of transitions section")
}

func (tm *TuringMachineCompiler) processSingleTransition(states map[string]automaton.State, symbols map[string]automaton.Symbol, tf automaton.TMTransitionFunction) error {
	// Each transition is as follows:
	// (state, symbol) > (state, symbol, movement)
	// At this point '(' has already been processed
	const atEndErrMsg = "unfinished transition"
	leftSide, err := tm.processTransitionLeftSide(states, symbols, atEndErrMsg)
	if err != nil {
		return err
	}
	if _, err := tm.consumeTokenWithType(atEndErrMsg, lexer.ArrowToken); err != nil {
		return err
	}
	rightSide, err := tm.processTransitionRightSide(states, symbols, atEndErrMsg)
	if err != nil {
		return err
	}
	tf[leftSide] = rightSide
	return nil
}

func (tm *TuringMachineCompiler) processTransitionLeftSide(states map[string]automaton.State, symbols map[string]automaton.Symbol, atEndErrMsg string) (automaton.TMTransitionKey, error) {
	var zero automaton.TMTransitionKey
	state, err := tm.consumeTokenWithType(atEndErrMsg, lexer.StateToken)
	if err != nil {
		return zero, err
	}
	if _, ok := states[state.Value]; !ok {
		return zero, fmt.Errorf("undefined state %s used in transition function left side", state.Value)
	}
	if _, err := tm.consumeTokenWithType(atEndErrMsg, lexer.CommaToken); err != nil {
		return zero, err
	}
	symbol, err := tm.consumeTokenWithType(atEndErrMsg, lexer.SymbolToken, lexer.BlankSymbolToken)
	if err != nil {
		return zero, err
	}
	if _, ok := symbols[symbol.Value]; !ok {
		return zero, fmt.Errorf("undefined symbol %s used in transition function left side", symbol.Value)
	}
	if _, err := tm.consumeTokenWithType(atEndErrMsg, lexer.RightParenToken); err != nil {
		return zero, err
	}
	return automaton.TMTransitionKey{StateName: state.Value, SymbolName: symbol.Value}, nil
}

func (tm *TuringMachineCompiler) processTransitionRightSide(states map[string]automaton.State, symbols map[string]automaton.Symbol, atEndErrMsg string) (automaton.TMTransitionValue, error) {
	var zero automaton.TMTransitionValue
	if _, err := tm.consumeTokenWithType(atEndErrMsg, lexer.LeftParenToken); err != nil {
		return zero, err
	}
	state, err := tm.consumeTokenWithType(atEndErrMsg, lexer.StateToken)
	if err != nil {
		return zero, err
	}
	if _, ok := states[state.Value]; !ok {
		return zero, fmt.Errorf("undefined state %s used in transition function right side", state.Value)
	}
	if _, err := tm.consumeTokenWithType(atEndErrMsg, lexer.CommaToken); err != nil {
		return zero, err
	}
	symbol, err := tm.consumeTokenWithType(atEndErrMsg, lexer.SymbolToken, lexer.BlankSymbolToken)
	if err != nil {
		return zero, err
	}
	if _, ok := symbols[symbol.Value]; !ok {
		return zero, fmt.Errorf("undefined symbol %s used in transition function right side", symbol.Value)
	}
	if _, err := tm.consumeTokenWithType(atEndErrMsg, lexer.CommaToken); err != nil {
		return zero, err
	}
	move, err := tm.consumeTokenWithType(atEndErrMsg, lexer.MoveLeftToken, lexer.MoveRightToken)
	if err != nil {
		return zero, err
	}
	if _, err := tm.consumeTokenWithType(atEndErrMsg, lexer.RightParenToken); err != nil {
		return zero, err
	}
	moveValue := automaton.TapeMoveLeft
	if move.Type == lexer.MoveRightToken {
		moveValue = automaton.TapeMoveRight
	}
	return automaton.TMTransitionValue{StateName: state.Value, SymbolName: symbol.Value, Move: moveValue}, nil
}

func (tm *TuringMachineCompiler) processTape(symbols map[string]automaton.Symbol) ([]string, error) {
	tape := make([]string, 0)
	for !tm.isAtEnd() {
		t := tm.advance()
		switch t.Type {
		case lexer.SemicolonToken:
			if len(tape) == 0 {
				tape = append(tape, automaton.BlankSymbol.Name)
			}
			return tape, nil
		case lexer.SymbolToken:
			if _, ok := symbols[t.Value]; !ok {
				return nil, fmt.Errorf("invalid symbol %s in initial tape, each symbol must be defined in symbols section", t.Value)
			}
			tape = append(tape, t.Value)
		case lexer.BlankSymbolToken:
			tape = append(tape, t.Value)
		default:
			return nil, fmt.Errorf("invalid token type, expected: %s, %s or %s, got: %s", lexer.SemicolonToken.String(), lexer.SymbolToken.String(), lexer.BlankSymbolToken.String(), t.Type.String())
		}
	}
	return nil, errors.New("missing ';' at the end of tape section")
}
