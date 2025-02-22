package compiler

import (
	"automata-compiler/pkg/automaton"
	"automata-compiler/pkg/lexer"
	"errors"
	"fmt"
)

type DeterministicFiniteAutomatonCompiler struct {
	BaseCompiler
}

func NewDeterministicFiniteAutomatonCompiler(tokens []lexer.Token) *DeterministicFiniteAutomatonCompiler {
	return &DeterministicFiniteAutomatonCompiler{BaseCompiler: newBaseCompiler(tokens)}
}

func (dfa *DeterministicFiniteAutomatonCompiler) Compile() (automaton.Automaton, error) {
	states, err := dfa.processStates()
	if err != nil {
		return nil, dfa.addLinePrefixForErrPrevToken(err)
	}
	initialState, err := dfa.processInitialState(states)
	if err != nil {
		return nil, dfa.addLinePrefixForErrPrevToken(err)
	}
	err = dfa.processAcceptingStates(states)
	if err != nil {
		return nil, dfa.addLinePrefixForErrPrevToken(err)
	}
	// DFA doesn't have any special symbol so we pass an empty map
	symbols, err := dfa.processSymbols(make(map[string]automaton.Symbol))
	if err != nil {
		return nil, dfa.addLinePrefixForErrPrevToken(err)
	}
	tf, err := dfa.processTransitions(states, symbols)
	if err != nil {
		return nil, dfa.addLinePrefixForErrPrevToken(err)
	}
	input, err := dfa.processInput(symbols)
	if err != nil {
		return nil, dfa.addLinePrefixForErrPrevToken(err)
	}
	err = dfa.checkForCorrectEndingSequnce()
	if err != nil {
		// It's more lexer error than user provided source,
		// so we don't include line here
		return nil, err
	}
	return &automaton.DeterministicFiniteAutomaton{
		States:       states,
		Symbols:      symbols,
		CurrentState: initialState,
		Input:        input,
		InputIt:      0,
		Transitions:  tf,
	}, nil
}

func (dfa *DeterministicFiniteAutomatonCompiler) processTransitions(states map[string]automaton.State, symbols map[string]automaton.Symbol) (automaton.DFATransitionFunction, error) {
	tf := make(automaton.DFATransitionFunction)
	for !dfa.isAtEnd() {
		t := dfa.advance()
		switch t.Type {
		case lexer.SemicolonToken:
			return tf, nil
		case lexer.LeftParenToken:
			err := dfa.processSingleTransition(states, symbols, tf)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid token type, expected: %s or %s, got: %s", lexer.LeftParenToken.String(), lexer.SemicolonToken.String(), t.Type.String())

		}
	}
	return nil, errors.New("missing ';' at the end of transitions section")
}

func (dfa *DeterministicFiniteAutomatonCompiler) processSingleTransition(states map[string]automaton.State, symbols map[string]automaton.Symbol, tf automaton.DFATransitionFunction) error {
	// Each transition is as follows:
	// (state, symbol) > (state)
	// At this point '(' has already been processed
	const atEndErrMsg = "unfinished transition"
	leftSide, err := dfa.processTransitionLeftSide(states, symbols, atEndErrMsg)
	if err != nil {
		return err
	}
	if _, err := dfa.consumeTokenWithType(atEndErrMsg, lexer.ArrowToken); err != nil {
		return err
	}
	rightSide, err := dfa.processTransitionRightSide(states, atEndErrMsg)
	if err != nil {
		return err
	}
	tf[leftSide] = rightSide
	return nil
}

func (dfa *DeterministicFiniteAutomatonCompiler) processTransitionLeftSide(states map[string]automaton.State, symbols map[string]automaton.Symbol, atEndErrMsg string) (automaton.DFATransitionKey, error) {
	var zero automaton.DFATransitionKey
	state, err := dfa.consumeTokenWithType(atEndErrMsg, lexer.StateToken)
	if err != nil {
		return zero, err
	}
	if _, ok := states[state.Value]; !ok {
		return zero, fmt.Errorf("undefined state %s used in transition function left side", state.Value)
	}
	if _, err := dfa.consumeTokenWithType(atEndErrMsg, lexer.CommaToken); err != nil {
		return zero, err
	}
	symbol, err := dfa.consumeTokenWithType(atEndErrMsg, lexer.SymbolToken, lexer.BlankSymbolToken)
	if err != nil {
		return zero, err
	}
	if _, ok := symbols[symbol.Value]; !ok {
		return zero, fmt.Errorf("undefined symbol %s used in transition function left side", symbol.Value)
	}
	if _, err := dfa.consumeTokenWithType(atEndErrMsg, lexer.RightParenToken); err != nil {
		return zero, err
	}
	return automaton.DFATransitionKey{StateName: state.Value, SymbolName: symbol.Value}, nil
}

func (dfa *DeterministicFiniteAutomatonCompiler) processTransitionRightSide(states map[string]automaton.State, atEndErrMsg string) (automaton.DFATransitionValue, error) {
	var zero automaton.DFATransitionValue
	if _, err := dfa.consumeTokenWithType(atEndErrMsg, lexer.LeftParenToken); err != nil {
		return zero, err
	}
	state, err := dfa.consumeTokenWithType(atEndErrMsg, lexer.StateToken)
	if err != nil {
		return zero, err
	}
	if _, ok := states[state.Value]; !ok {
		return zero, fmt.Errorf("undefined state %s used in transition function right side", state.Value)
	}
	if _, err := dfa.consumeTokenWithType(atEndErrMsg, lexer.RightParenToken); err != nil {
		return zero, err
	}
	return automaton.DFATransitionValue{StateName: state.Value}, nil
}

func (dfa *DeterministicFiniteAutomatonCompiler) processInput(symbols map[string]automaton.Symbol) ([]string, error) {
	input := make([]string, 0)
	for !dfa.isAtEnd() {
		t := dfa.advance()
		switch t.Type {
		case lexer.SemicolonToken:
			if len(input) == 0 {
				input = append(input, automaton.BlankSymbol.Name)
			}
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
