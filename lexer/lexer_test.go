package lexer

import (
	"slices"
	"testing"
)

func TestNewLexer(t *testing.T) {
	sourceCode := "source code"
	result := NewLexer(sourceCode)
	expected := &Lexer{
		source:  sourceCode,
		tokens:  make([]Token, 0),
		line:    1,
		start:   0,
		current: 0,
	}
	if result.source != expected.source {
		t.Errorf("invalid source, expected: %s, got: %s", expected.source, result.source)
	}
	if !slices.Equal(expected.tokens, result.tokens) {
		t.Error("invalid tokens, expected:", expected.tokens, " got:", result.tokens)
	}
	if result.line != expected.line {
		t.Errorf("invalid line, expected: %d, got: %d", expected.line, result.line)
	}
	if result.start != expected.start {
		t.Errorf("invalid start, expected: %d, got: %d", expected.start, result.start)
	}
	if result.current != expected.current {
		t.Errorf("invalid current, expected: %d, got: %d", expected.current, result.current)
	}
}

func TestScanTokens(t *testing.T) {
	data := []struct {
		name           string
		sourceCode     string
		expectedTokens []Token
		expectedErrMsg string
	}{
		{
			"empty",
			"",
			[]Token{
				{Type: EOFToken, Value: "", Line: 1},
			},
			"",
		},
		{
			"single state",
			"qSingle",
			[]Token{
				{Type: StateToken, Value: "qSingle", Line: 1},
				{Type: EOFToken, Value: "", Line: 1},
			},
			"",
		},
		{
			"few states",
			"qSingle qNext qNext   qAfterCouple\n\n\nqNext\t\r\n    qLast\n\n",
			[]Token{
				{Type: StateToken, Value: "qSingle", Line: 1},
				{Type: StateToken, Value: "qNext", Line: 1},
				{Type: StateToken, Value: "qNext", Line: 1},
				{Type: StateToken, Value: "qAfterCouple", Line: 1},
				{Type: StateToken, Value: "qNext", Line: 4},
				{Type: StateToken, Value: "qLast", Line: 5},
				{Type: EOFToken, Value: "", Line: 7},
			},
			"",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			lexer := NewLexer(d.sourceCode)
			result, err := lexer.ScanTokens()
			if !slices.Equal(d.expectedTokens, result) {
				t.Error("invalid tokens, expected:", d.expectedTokens, ", got:", result)
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
