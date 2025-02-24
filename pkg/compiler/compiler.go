package compiler

import (
	"automata-compiler/pkg/automaton"
	"automata-compiler/pkg/lexer"
	"errors"
	"fmt"
	"strings"
)

type Compiler interface {
	Compile() (automaton.Automaton, error)
}

// BaseCompiler implements simple utility functions that every automaton compiler needs
type BaseCompiler struct {
	tokens []lexer.Token
	it     int
}

func newBaseCompiler(tokens []lexer.Token) BaseCompiler {
	return BaseCompiler{
		tokens: tokens,
		it:     0,
	}
}

func (c BaseCompiler) isAtEnd() bool {
	return c.it >= len(c.tokens)
}

func (c *BaseCompiler) advance() lexer.Token {
	if c.isAtEnd() {
		var t lexer.Token
		return t
	}
	t := c.tokens[c.it]
	c.it++
	return t
}

// peek is same as `advance` but doesn't move the `it`
func (c BaseCompiler) peek() lexer.Token {
	if c.isAtEnd() {
		var t lexer.Token
		return t
	}
	return c.tokens[c.it]
}

func (c *BaseCompiler) consumeTokenWithType(atEndErrMsg string, expected ...lexer.TokenType) (lexer.Token, error) {
	var zero lexer.Token
	if c.isAtEnd() {
		return zero, errors.New(atEndErrMsg)
	}
	token := c.advance()
	if err := checkTokenType(token, expected...); err != nil {
		return zero, err
	}
	return token, nil
}

func (c BaseCompiler) prevTokenLine() int {
	if c.it == 0 {
		return 0
	}
	return c.tokens[c.it-1].Line
}

func (c BaseCompiler) addLinePrefixForErrPrevToken(err error) error {
	return addLinePrefixForErr(err, c.prevTokenLine())
}

func (c *BaseCompiler) checkForCorrectEndingSequnce() error {
	if _, err := c.consumeTokenWithType("missing EOF token at the end of source", lexer.EOFToken); err != nil {
		return err
	}
	if !c.isAtEnd() {
		return errors.New("unexpected token after EOF token")
	}
	return nil
}

func (c *BaseCompiler) processStates() (map[string]automaton.State, error) {
	states := make(map[string]automaton.State)
	for !c.isAtEnd() {
		t := c.advance()
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

func (c *BaseCompiler) processInitialState(states map[string]automaton.State) (string, error) {
	initialState, err := c.consumeTokenWithType("missing initial state section", lexer.StateToken)
	if err != nil {
		return "", err
	}
	if _, ok := states[initialState.Value]; !ok {
		return "", fmt.Errorf("invalid initial state, state %s was not declared in states list", initialState.Value)
	}
	if c.isAtEnd() || c.advance().Type != lexer.SemicolonToken {
		return "", fmt.Errorf("missing ';' after initial state")
	}
	return initialState.Value, nil
}

func (c *BaseCompiler) processAcceptingStates(states map[string]automaton.State) error {
	for !c.isAtEnd() {
		t := c.advance()
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

// symbols are passed so if automaton uses some special symbols (e.g. blank symbol in tm, input end and stack start symbols in pa) it can add it to list of symbols provided by user
func (c *BaseCompiler) processSymbols(symbols map[string]automaton.Symbol) (map[string]automaton.Symbol, error) {
	for !c.isAtEnd() {
		t := c.advance()
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

// addLinePrefixForErr modifies provided error message, adding prefix which contains line number
func addLinePrefixForErr(err error, line int) error {
	return fmt.Errorf("[Line %d] %s", line, err.Error())
}

// checkTokenType returns an error if provided `t` is not one of the `expected` types
//
// it panics when no token in `expected` is provided
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
