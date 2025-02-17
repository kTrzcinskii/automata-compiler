package compiler

import (
	"automata-compiler/pkg/automaton"
	"automata-compiler/pkg/lexer"
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
	data := []struct {
		name           string
		tokens         []lexer.Token
		expected       *automaton.TuringMachine
		expectedErrMsg string
	}{
		{
			"simple program",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				// First
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 6},
				{Type: lexer.RightParenToken, Value: ")", Line: 6},
				{Type: lexer.ArrowToken, Value: ">", Line: 6},
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState2", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.MoveLeftToken, Value: "L", Line: 6},
				{Type: lexer.RightParenToken, Value: ")", Line: 6},
				// Second
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState3", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 6},
				{Type: lexer.RightParenToken, Value: ")", Line: 6},
				{Type: lexer.ArrowToken, Value: ">", Line: 6},
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.MoveRightToken, Value: "R", Line: 6},
				{Type: lexer.RightParenToken, Value: ")", Line: 6},
				{Type: lexer.SemicolonToken, Value: ";", Line: 6},
				// Initial tape
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 7},
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 7},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 7},
				{Type: lexer.BlankSymbolToken, Value: "B", Line: 7},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 7},
				{Type: lexer.SemicolonToken, Value: ";", Line: 7},
				{Type: lexer.EOFToken, Value: "", Line: 7},
			},
			&automaton.TuringMachine{
				States: map[string]automaton.State{
					"qState":  {Name: "qState"},
					"qState2": {Name: "qState2", Accepting: true},
					"qState3": {Name: "qState3", Accepting: true},
					"qState4": {Name: "qState4"},
				},
				CurrentState: "qState",
				Symbols: map[string]automaton.Symbol{
					automaton.BlankSymbol.Name: automaton.BlankSymbol,
					"symbol1":                  {Name: "symbol1"},
					"symbol2":                  {Name: "symbol2"},
					"symbol3":                  {Name: "symbol3"},
				},
				Transitions: map[automaton.TMTransitionKey]automaton.TMTransitionValue{
					{StateName: "qState", SymbolName: "symbol1"}:  {StateName: "qState2", SymbolName: "symbol2", Move: automaton.TapeMoveLeft},
					{StateName: "qState3", SymbolName: "symbol3"}: {StateName: "qState", SymbolName: "symbol3", Move: automaton.TapeMoveRight},
				},
				Tape: []string{
					"symbol2",
					"symbol1",
					"symbol2",
					"B",
					"symbol3",
				},
				TapeIt: 0,
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
			nil,
			"[Line 1] missing ';' at the end of states section",
		},
		{
			"duplicated state",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState", Line: 1},
			},
			nil,
			"[Line 1] state qState already declared, each state must have unique name",
		},
		{
			"unexpected token in state section",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.SymbolToken, Value: "Symbol", Line: 1},
			},
			nil,
			"[Line 1] invalid token type, expected: StateToken or SemicolonToken, got: SymbolToken",
		},
		{
			"no state defined",
			[]lexer.Token{
				{Type: lexer.SemicolonToken, Value: ";", Line: 1},
			},
			nil,
			"[Line 1] there must be at least one state defined",
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
			nil,
			"[Line 3] missing initial state section",
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
			nil,
			"[Line 3] invalid token type, expected: StateToken, got: SymbolToken",
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
			nil,
			"[Line 3] invalid initial state, state qState5 was not declared in states list",
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
			nil,
			"[Line 3] missing ';' after initial state",
		},
		{
			"missing semicolon after accepting states section",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
			},
			nil,
			"[Line 4] missing ';' at the end of accepting states section",
		},
		{
			"invalid token in accepting states section",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				{Type: lexer.ArrowToken, Value: ">", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
			},
			nil,
			"[Line 3] invalid token type, expected: StateToken or SemicolonToken, got: ArrowToken",
		},
		{
			"undefined state in accepting states section",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				{Type: lexer.StateToken, Value: "qState5", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
			},
			nil,
			"[Line 4] state qState5 not found, any accepting state must be defined in state list",
		},
		{
			"missing semicolon after symbols section",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
			},
			nil,
			"[Line 5] missing ';' at the end of symbols section",
		},
		{
			"invalid token in symbols section",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.LeftParenToken, Value: "(", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
			},
			nil,
			"[Line 5] invalid token type, expected: SymbolToken or SemicolonToken, got: LeftParenToken",
		},
		{
			"duplicated symbol in symbols section",
			[]lexer.Token{
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
			},
			nil,
			"[Line 5] symbol symbol1 already declared, each symbol must have unique name",
		},
		{
			"missing semicolon after transitions section",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				// First
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 6},
				{Type: lexer.RightParenToken, Value: ")", Line: 6},
				{Type: lexer.ArrowToken, Value: ">", Line: 6},
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState2", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.MoveLeftToken, Value: "L", Line: 6},
				{Type: lexer.RightParenToken, Value: ")", Line: 6},
				// Second
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState3", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 6},
				{Type: lexer.RightParenToken, Value: ")", Line: 6},
				{Type: lexer.ArrowToken, Value: ">", Line: 6},
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.MoveRightToken, Value: "R", Line: 6},
				{Type: lexer.RightParenToken, Value: ")", Line: 6},
			},
			nil,
			"[Line 6] missing ';' at the end of transitions section",
		},
		{
			"invalid token at the beginning of the transition",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
			},
			nil,
			"[Line 5] invalid token type, expected: LeftParenToken or SemicolonToken, got: SymbolToken",
		},
		{
			"unfinished transition",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 6},
			},
			nil,
			"[Line 6] unfinished transition",
		},
		{
			"undefined state in left side transition",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState5", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 6},
			},
			nil,
			"[Line 6] undefined state qState5 used in transition function left side",
		},
		{
			"undefined symbol in left side transition",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol4", Line: 6},
			},
			nil,
			"[Line 6] undefined symbol symbol4 used in transition function left side",
		},
		{
			"undefined state in right side transition",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 6},
				{Type: lexer.RightParenToken, Value: ")", Line: 6},
				{Type: lexer.ArrowToken, Value: ">", Line: 6},
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState10", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.MoveRightToken, Value: "R", Line: 6},
			},
			nil,
			"[Line 6] undefined state qState10 used in transition function right side",
		},
		{
			"undefined symbol in right side transition",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 6},
				{Type: lexer.RightParenToken, Value: ")", Line: 6},
				{Type: lexer.ArrowToken, Value: ">", Line: 6},
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol30", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.MoveRightToken, Value: "R", Line: 6},
			},
			nil,
			"[Line 6] undefined symbol symbol30 used in transition function right side",
		},
		{
			"invalid movement token in right side transition",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 6},
				{Type: lexer.RightParenToken, Value: ")", Line: 6},
				{Type: lexer.ArrowToken, Value: ">", Line: 6},
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.ArrowToken, Value: ">", Line: 6},
			},
			nil,
			"[Line 6] invalid token type, expected one of: MoveLeftToken, MoveRightToken, got: ArrowToken",
		},
		{
			"missing semicolon after tape section",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				{Type: lexer.SemicolonToken, Value: ";", Line: 6},
				// Initial tape
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 7},
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 7},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 7},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 7},
			},
			nil,
			"[Line 7] missing ';' at the end of tape section",
		},
		{
			"undefined symbol in initial tape",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				{Type: lexer.SemicolonToken, Value: ";", Line: 6},
				// Initial tape
				{Type: lexer.SymbolToken, Value: "symbol100", Line: 7},
			},
			nil,
			"[Line 7] invalid symbol symbol100 in initial tape, each symbol must be defined in symbols section",
		},
		{
			"invalid token type in tape section",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				{Type: lexer.SemicolonToken, Value: ";", Line: 6},
				// Initial tape
				{Type: lexer.ArrowToken, Value: ">", Line: 7},
			},
			nil,
			"[Line 7] invalid token type, expected: SemicolonToken, SymbolToken or BlankSymbolToken, got: ArrowToken",
		},
		{
			"missing EOF at the end of token list",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				{Type: lexer.SemicolonToken, Value: ";", Line: 6},
				// Initial tape
				{Type: lexer.SemicolonToken, Value: ";", Line: 7},
			},
			nil,
			"missing EOF token at the end of source",
		},
		{
			"unexpected token after EOF",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				{Type: lexer.SemicolonToken, Value: ";", Line: 6},
				// Initial tape
				{Type: lexer.SemicolonToken, Value: ";", Line: 7},
				// EOF
				{Type: lexer.EOFToken, Value: "", Line: 7},
				{Type: lexer.SemicolonToken, Value: ";", Line: 7},
			},
			nil,
			"unexpected token after EOF token",
		},
		{
			"invalid usage of blank symbol in user-defined symbols",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.StateToken, Value: "qState2", Line: 1},
				{Type: lexer.StateToken, Value: "qState3", Line: 1},
				{Type: lexer.StateToken, Value: "qState4", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.StateToken, Value: "qState2", Line: 4},
				{Type: lexer.StateToken, Value: "qState3", Line: 4},
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SymbolToken, Value: "symbol1", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol2", Line: 5},
				{Type: lexer.SymbolToken, Value: "symbol3", Line: 5},
				{Type: lexer.BlankSymbolToken, Value: "B", Line: 5},
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
			},
			nil,
			"[Line 5] invalid token type, expected: SymbolToken or SemicolonToken, got: BlankSymbolToken",
		},
		{
			"empty tape should contain one blank symbol",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				{Type: lexer.SemicolonToken, Value: ";", Line: 6},
				// Initial tape
				{Type: lexer.SemicolonToken, Value: ";", Line: 7},
				{Type: lexer.EOFToken, Value: "", Line: 7},
			},
			&automaton.TuringMachine{
				States: map[string]automaton.State{
					"qState": {Name: "qState"},
				},
				CurrentState: "qState",
				Symbols: map[string]automaton.Symbol{
					automaton.BlankSymbol.Name: automaton.BlankSymbol,
				},
				Transitions: map[automaton.TMTransitionKey]automaton.TMTransitionValue{},
				Tape: []string{
					"B",
				},
				TapeIt: 0,
			},
			"",
		},
		{
			"blank symbol can be used in transitions",
			[]lexer.Token{
				// States
				{Type: lexer.StateToken, Value: "qState", Line: 1},
				{Type: lexer.SemicolonToken, Value: ";", Line: 2},
				// Initial state
				{Type: lexer.StateToken, Value: "qState", Line: 3},
				{Type: lexer.SemicolonToken, Value: ";", Line: 3},
				// Accepting states
				{Type: lexer.SemicolonToken, Value: ";", Line: 4},
				// Symbols
				{Type: lexer.SemicolonToken, Value: ";", Line: 5},
				// Transitions
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.BlankSymbolToken, Value: "B", Line: 6},
				{Type: lexer.RightParenToken, Value: ")", Line: 6},
				{Type: lexer.ArrowToken, Value: ">", Line: 6},
				{Type: lexer.LeftParenToken, Value: "(", Line: 6},
				{Type: lexer.StateToken, Value: "qState", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.BlankSymbolToken, Value: "B", Line: 6},
				{Type: lexer.CommaToken, Value: ",", Line: 6},
				{Type: lexer.MoveLeftToken, Value: "R", Line: 6},
				{Type: lexer.RightParenToken, Value: ")", Line: 6},
				{Type: lexer.SemicolonToken, Value: ";", Line: 6},
				// Initial tape
				{Type: lexer.SemicolonToken, Value: ";", Line: 7},
				{Type: lexer.EOFToken, Value: "", Line: 7},
			},
			&automaton.TuringMachine{
				States: map[string]automaton.State{
					"qState": {Name: "qState"},
				},
				CurrentState: "qState",
				Symbols: map[string]automaton.Symbol{
					automaton.BlankSymbol.Name: automaton.BlankSymbol,
				},
				Transitions: map[automaton.TMTransitionKey]automaton.TMTransitionValue{
					{StateName: "qState", SymbolName: automaton.BlankSymbol.Name}: {StateName: "qState", SymbolName: automaton.BlankSymbol.Name, Move: automaton.TapeMoveLeft},
				},
				Tape: []string{
					"B",
				},
				TapeIt: 0,
			},
			"",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			tmc := NewTuringMachineCompiler(d.tokens)
			result, err := tmc.Compile()
			if !(d.expected == nil && result == nil) {
				if diff := cmp.Diff(d.expected, result); diff != "" {
					t.Error(diff)
				}
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
