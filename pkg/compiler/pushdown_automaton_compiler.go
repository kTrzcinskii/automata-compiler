package compiler

import (
	"automata-compiler/pkg/automaton"
	"automata-compiler/pkg/lexer"
	"errors"
	"fmt"
)

type PushdownAutomatonCompiler struct {
	BaseCompiler
}

func NewPushdownAutomatonCompiler(tokens []lexer.Token) *PushdownAutomatonCompiler {
	return &PushdownAutomatonCompiler{BaseCompiler: newBaseCompiler(tokens)}
}

func (pa *PushdownAutomatonCompiler) Compile() (automaton.Automaton, error) {
	states, err := pa.processStates()
	if err != nil {
		return nil, pa.addLinePrefixForErrPrevToken(err)
	}
	initialState, err := pa.processInitialState(states)
	if err != nil {
		return nil, pa.addLinePrefixForErrPrevToken(err)
	}
	err = pa.processAcceptingStates(states)
	if err != nil {
		return nil, pa.addLinePrefixForErrPrevToken(err)
	}
	specialSymbols := pa.getSpecialSymbols()
	symbols, err := pa.processSymbols(specialSymbols)
	if err != nil {
		return nil, pa.addLinePrefixForErrPrevToken(err)
	}
	tf, err := pa.processTransitions(states, symbols)
	if err != nil {
		return nil, pa.addLinePrefixForErrPrevToken(err)
	}
	initialInput, err := pa.processInput(symbols)
	if err != nil {
		return nil, pa.addLinePrefixForErrPrevToken(err)
	}
	err = pa.checkForCorrectEndingSequnce()
	if err != nil {
		// It's more lexer error than user provided source,
		// so we don't include line here
		return nil, err
	}
	return &automaton.PushdownAutomaton{
		States:       states,
		Symbols:      symbols,
		CurrentState: initialState,
		Input:        initialInput,
		InputIt:      0,
		Stack:        []string{automaton.StackStartSymbol.Name},
		Transitions:  tf,
	}, nil
}

func (pa PushdownAutomatonCompiler) getSpecialSymbols() map[string]automaton.Symbol {
	symbols := make(map[string]automaton.Symbol)
	symbols[automaton.InputEndSymbol.Name] = automaton.InputEndSymbol
	symbols[automaton.StackStartSymbol.Name] = automaton.StackStartSymbol
	return symbols
}

func (pa *PushdownAutomatonCompiler) processTransitions(states map[string]automaton.State, symbols map[string]automaton.Symbol) (automaton.PATransitionFunction, error) {
	tf := make(automaton.PATransitionFunction)
	for !pa.isAtEnd() {
		t := pa.advance()
		switch t.Type {
		case lexer.SemicolonToken:
			return tf, nil
		case lexer.LeftParenToken:
			err := pa.processSingleTransition(states, symbols, tf)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid token type, expected: %s or %s, got: %s", lexer.LeftParenToken.String(), lexer.SemicolonToken.String(), t.Type.String())
		}
	}
	return nil, errors.New("missing ';' at the end of transitions section")
}

func (pa *PushdownAutomatonCompiler) processSingleTransition(states map[string]automaton.State, symbols map[string]automaton.Symbol, tf automaton.PATransitionFunction) error {
	// Each transition is as follows:
	// (state, input_symbol, stack_symbol) > (state, stack_symbol1, stack_symbol2, ...)
	// At this point '(' has already been processed
	const atEndErrMsg = "unfinished transition"
	leftSide, err := pa.processTransitionLeftSide(states, symbols, atEndErrMsg)
	if err != nil {
		return err
	}
	if _, err := pa.consumeTokenWithType(atEndErrMsg, lexer.ArrowToken); err != nil {
		return err
	}
	rightSide, err := pa.processTransitionRightSide(states, symbols, atEndErrMsg)
	if err != nil {
		return err
	}
	tf[leftSide] = rightSide
	return nil
}

func (pa *PushdownAutomatonCompiler) processTransitionLeftSide(states map[string]automaton.State, symbols map[string]automaton.Symbol, atEndErrMsg string) (automaton.PATransitionKey, error) {
	var zero automaton.PATransitionKey
	state, err := pa.consumeTokenWithType(atEndErrMsg, lexer.StateToken)
	if err != nil {
		return zero, err
	}
	if _, ok := states[state.Value]; !ok {
		return zero, fmt.Errorf("undefined state %s used in transition function left side", state.Value)
	}
	if _, err := pa.consumeTokenWithType(atEndErrMsg, lexer.CommaToken); err != nil {
		return zero, err
	}
	inputSymbol, err := pa.consumeTokenWithType(atEndErrMsg, lexer.SymbolToken, lexer.InputEndToken)
	if err != nil {
		return zero, err
	}
	if _, ok := symbols[inputSymbol.Value]; !ok {
		return zero, fmt.Errorf("undefined input symbol %s used in transition function left side", inputSymbol.Value)
	}
	if _, err := pa.consumeTokenWithType(atEndErrMsg, lexer.CommaToken); err != nil {
		return zero, err
	}
	stackSymbol, err := pa.consumeTokenWithType(atEndErrMsg, lexer.SymbolToken, lexer.StackStartToken)
	if err != nil {
		return zero, err
	}
	if _, ok := symbols[stackSymbol.Value]; !ok {
		return zero, fmt.Errorf("undefined stack symbol %s used in transition function left side", stackSymbol.Value)
	}
	if _, err := pa.consumeTokenWithType(atEndErrMsg, lexer.RightParenToken); err != nil {
		return zero, err
	}
	return automaton.PATransitionKey{
		StateName:       state.Value,
		InputSymbolName: inputSymbol.Value,
		StackSymbolName: stackSymbol.Value,
	}, nil
}

func (pa *PushdownAutomatonCompiler) processTransitionRightSide(states map[string]automaton.State, symbols map[string]automaton.Symbol, atEndErrMsg string) (automaton.PATransitionValue, error) {
	var zero automaton.PATransitionValue
	if _, err := pa.consumeTokenWithType(atEndErrMsg, lexer.LeftParenToken); err != nil {
		return zero, err
	}
	state, err := pa.consumeTokenWithType(atEndErrMsg, lexer.StateToken)
	if err != nil {
		return zero, err
	}
	if _, ok := states[state.Value]; !ok {
		return zero, fmt.Errorf("undefined state %s used in transition function right side", state.Value)
	}
	stackSymbols := make([]string, 0)
	for pa.peek().Type == lexer.CommaToken {
		// Consume comma
		pa.advance()
		stackSymbol, err := pa.consumeTokenWithType(atEndErrMsg, lexer.SymbolToken, lexer.StackStartToken)
		if err != nil {
			return zero, err
		}
		if _, ok := symbols[stackSymbol.Value]; !ok {
			return zero, fmt.Errorf("undefined stack symbol %s used in transition function right side", stackSymbol.Value)
		}
		stackSymbols = append(stackSymbols, stackSymbol.Value)
	}
	// We pass CommaToken here only for better error message, at this point we know it can
	// only be RightParenToken
	if _, err := pa.consumeTokenWithType(atEndErrMsg, lexer.RightParenToken, lexer.CommaToken); err != nil {
		return zero, err
	}
	return automaton.PATransitionValue{
		StateName:        state.Value,
		StackSymbolNames: stackSymbols,
	}, nil
}

func (pa *PushdownAutomatonCompiler) processInput(symbols map[string]automaton.Symbol) ([]string, error) {
	input := make([]string, 0)
	for !pa.isAtEnd() {
		t := pa.advance()
		switch t.Type {
		case lexer.SemicolonToken:
			input = append(input, automaton.InputEndSymbol.Name)
			return input, nil
		case lexer.SymbolToken:
			if _, ok := symbols[t.Value]; !ok {
				return nil, fmt.Errorf("invalid symbol %s in input, each symbol must be defined in symbols section", t.Value)
			}
			input = append(input, t.Value)
		default:
			return nil, fmt.Errorf("invalid token type, expected: %s or %s, got: %s", lexer.SemicolonToken.String(), lexer.SymbolToken.String(), t.Type.String())
		}
	}
	return nil, errors.New("missing ';' at the end of input section")
}
