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
