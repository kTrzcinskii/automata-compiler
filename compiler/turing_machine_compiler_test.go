package compiler

import (
	"automata-compiler/automata"
	"automata-compiler/lexer"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewTuringMachineCompiler(t *testing.T) {
	tokens := []lexer.Token{
		{Type: lexer.StateToken, Value: "qState", Line: 1},
		{Type: lexer.SemicolonToken, Value: ";", Line: 2},
		{Type: lexer.SemicolonToken, Value: ";", Line: 2},
		{Type: lexer.CommaToken, Value: ",", Line: 2},
		{Type: lexer.ArrowToken, Value: ">", Line: 2},
		{Type: lexer.ArrowToken, Value: ">", Line: 2},
		{Type: lexer.SymbolToken, Value: "symbol1", Line: 2},
		{Type: lexer.SymbolToken, Value: "symbol2", Line: 2},
		{Type: lexer.BlankSymbolToken, Value: "B", Line: 2},
		{Type: lexer.MoveLeftToken, Value: "L", Line: 2},
		{Type: lexer.MoveRightToken, Value: "R", Line: 2},
		{Type: lexer.MoveRightToken, Value: "R", Line: 2},
		{Type: lexer.CommaToken, Value: ",", Line: 2},
		{Type: lexer.CommaToken, Value: ",", Line: 2},
		{Type: lexer.LeftParenToken, Value: "(", Line: 2},
		{Type: lexer.LeftParenToken, Value: "(", Line: 2},
		{Type: lexer.RightParenToken, Value: ")", Line: 3},
		{Type: lexer.EOFToken, Value: "", Line: 3},
	}
	result := NewTuringMachineCompiler(tokens)
	expected := TuringMachineCompiler{
		tokens: tokens,
		it:     0,
	}
	if !slices.Equal(expected.tokens, result.tokens) {
		t.Error("invalid tokens, expected:", expected.tokens, " got:", result.tokens)
	}
	if result.it != expected.it {
		t.Errorf("invalid it, expected: %d, got: %d", expected.it, result.it)
	}
}

func TestCompile(t *testing.T) {
	var zero automata.TuringMachine
	data := []struct {
		name           string
		tokens         []lexer.Token
		expected       automata.TuringMachine
		expectedErrMsg string
	}{
		{
			"simple program",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
			},
			automata.TuringMachine{
				States: map[string]automata.State{
					"qState":  {Name: "qState"},
					"qState2": {Name: "qState2"},
					"qState3": {Name: "qState3"},
					"qState4": {Name: "qState4"},
				},
				InitialState: "qState",
			},
			"",
		},
		{
			"missing semicolon after states",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
			},
			zero,
			"missing ';' at the end of states section",
		},
		{
			"duplicated state",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState", Line: 1},
			},
			zero,
			"state qState already declared, each state must have unique name",
		},
		{
			"unexpected token in state section",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.SymbolToken, Value: "Symbol", Line: 1},
			},
			zero,
			"invalid token type, expected: StateToken or SemicolonToken, got: SymbolToken",
		},
		{
			"missing initial state",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
			},
			zero,
			"missing initial state section",
		},
		{
			"invalid token in initial state section",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				{Type: lexer.SymbolToken, Value: "symbol", Line: 3},
			},
			zero,
			"invalid initial state token, expected: StateToken, got: SymbolToken",
		},
		{
			"unknown state in initial state section",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				{Type: lexer.StateToken, Value: "qState5", Line: 3},
			},
			zero,
			"invalid initial state, state qState5 was not declared in states list",
		},
		{
			"missing semicolon after initial state section",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				{Type: lexer.StateToken, Value: "qState", Line: 3},
			},
			zero,
			"missing ';' after initial state",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			tmc := NewTuringMachineCompiler(d.tokens)
			result, err := tmc.Compile()
			if diff := cmp.Diff(d.expected, result); diff != "" {
				t.Error(diff)
			}
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			if errMsg != d.expectedErrMsg {
				t.Errorf("invalid error message, expected: %s, got: %s", d.expectedErrMsg, errMsg)
			}

		})
	}
}
