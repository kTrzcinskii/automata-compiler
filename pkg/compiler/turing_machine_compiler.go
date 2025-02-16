package compiler

import (
	"automata-compiler/pkg/automata"
	"automata-compiler/pkg/lexer"
	"errors"
	"fmt"
	"strings"
)

type TuringMachineCompiler struct {
	tokens []lexer.Token
	it     int
}

func (tm *TuringMachineCompiler) Compile() (automata.Automata, error) {
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
	return &automata.TuringMachine{States: states, CurrentState: initialState, Symbols: symbols, Transitions: tf, Tape: initialTape, TapeIt: 0}, nil
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

func checkTokenType(t lexer.Token, expected ...lexer.TokenType) error {
	if expected == nil {
		panic("no token type provided")
	}
	expectedStr := make([]string, 0, len(expected))
	for _, v := range expected {
		if t.Type == v {
			return nil
		}
		expectedStr = append(expectedStr, v.String())
	}
	if len(expected) == 1 {
		return fmt.Errorf("invalid token type, expected: %s, got: %s", expectedStr[0], t.Type.String())
	}
	all := strings.Join(expectedStr, ", ")
	return fmt.Errorf("invalid token type, expected one of: %s, got: %s", all, t.Type.String())
}

func (tm *TuringMachineCompiler) consumeTokenWithType(atEndErrMsg string, expected ...lexer.TokenType) (lexer.Token, error) {
	var zero lexer.Token
	if tm.isAtEnd() {
		return zero, errors.New(atEndErrMsg)
	}
	token := tm.advance()
	if err := checkTokenType(token, expected...); err != nil {
		return zero, err
	}
	return token, nil
}

func (tm TuringMachineCompiler) prevTokenLine() int {
	if tm.it == 0 {
		return 0
	}
	return tm.tokens[tm.it-1].Line
}

func addLinePrefixForErr(err error, line int) error {
	return fmt.Errorf("[Line %d] %s", line, err.Error())
}

func (tm TuringMachineCompiler) addLinePrefixForErrPrevToken(err error) error {
	return addLinePrefixForErr(err, tm.prevTokenLine())
}

func (tm *TuringMachineCompiler) processStates() (map[string]automata.State, error) {
	states := make(map[string]automata.State)
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
			states[name] = automata.State{Name: name, Accepting: false}
		default:
			return nil, fmt.Errorf("invalid token type, expected: %s or %s, got: %s", lexer.StateToken.String(), lexer.SemicolonToken.String(), t.Type.String())
		}
	}
	return nil, errors.New("missing ';' at the end of states section")
}

func (tm *TuringMachineCompiler) processInitialState(states map[string]automata.State) (string, error) {
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

func (tm *TuringMachineCompiler) processSymbols() (map[string]automata.Symbol, error) {
	symbols := make(map[string]automata.Symbol)
	symbols[automata.BlankSymbol.Name] = automata.BlankSymbol
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
			symbols[name] = automata.Symbol{Name: t.Value}
		default:
			return nil, fmt.Errorf("invalid token type, expected: %s or %s, got: %s", lexer.SymbolToken.String(), lexer.SemicolonToken.String(), t.Type.String())
		}
	}
	return nil, errors.New("missing ';' at the end of symbols section")
}

func (tm *TuringMachineCompiler) processTransitions(states map[string]automata.State, symbols map[string]automata.Symbol) (automata.TransitionFunction, error) {
	tf := make(automata.TransitionFunction)
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

func (tm *TuringMachineCompiler) processSingleTransition(states map[string]automata.State, symbols map[string]automata.Symbol, tf automata.TransitionFunction) error {
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

func (tm *TuringMachineCompiler) processTransitionLeftSide(states map[string]automata.State, symbols map[string]automata.Symbol, atEndErrMsg string) (automata.TMTransitionKey, error) {
	var zero automata.TMTransitionKey
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
	return automata.TMTransitionKey{StateName: state.Value, SymbolName: symbol.Value}, nil
}

func (tm *TuringMachineCompiler) processTransitionRightSide(states map[string]automata.State, symbols map[string]automata.Symbol, atEndErrMsg string) (automata.TMTransitionValue, error) {
	var zero automata.TMTransitionValue
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
	moveValue := automata.TapeMoveLeft
	if move.Type == lexer.MoveRightToken {
		moveValue = automata.TapeMoveRight
	}
	return automata.TMTransitionValue{StateName: state.Value, SymbolName: symbol.Value, Move: moveValue}, nil
}

func (tm *TuringMachineCompiler) processTape(symbols map[string]automata.Symbol) ([]string, error) {
	tape := make([]string, 0)
	for !tm.isAtEnd() {
		t := tm.advance()
		switch t.Type {
		case lexer.SemicolonToken:
			if len(tape) == 0 {
				tape = append(tape, automata.BlankSymbol.Name)
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

func (tm *TuringMachineCompiler) checkForCorrectEndingSequnce() error {
	if _, err := tm.consumeTokenWithType("missing EOF token at the end of source", lexer.EOFToken); err != nil {
		return err
	}
	if !tm.isAtEnd() {
		return errors.New("unexpected token after EOF token")
	}
	return nil
}
