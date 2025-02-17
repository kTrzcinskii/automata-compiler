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
	var zeroTokens []Token
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
			"state",
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
		{
			"blank symbol",
			"B",
			[]Token{
				{Type: BlankSymbolToken, Value: "B", Line: 1},
				{Type: EOFToken, Value: "", Line: 1},
			},
			"",
		},
		{
			"move tokens",
			"L R",
			[]Token{
				{Type: MoveLeftToken, Value: "L", Line: 1},
				{Type: MoveRightToken, Value: "R", Line: 1},
				{Type: EOFToken, Value: "", Line: 1},
			},
			"",
		},
		{
			"parens",
			"( )",
			[]Token{
				{Type: LeftParenToken, Value: "(", Line: 1},
				{Type: RightParenToken, Value: ")", Line: 1},
				{Type: EOFToken, Value: "", Line: 1},
			},
			"",
		},
		{
			"comma, semicolon & arrow",
			", ; >",
			[]Token{
				{Type: CommaToken, Value: ",", Line: 1},
				{Type: SemicolonToken, Value: ";", Line: 1},
				{Type: ArrowToken, Value: ">", Line: 1},
				{Type: EOFToken, Value: "", Line: 1},
			},
			"",
		},
		{
			"symbol",
			"someSymbol",
			[]Token{
				{Type: SymbolToken, Value: "someSymbol", Line: 1},
				{Type: EOFToken, Value: "", Line: 1},
			},
			"",
		},
		{
			"invalid token",
			"|321321",
			zeroTokens,
			"[Line 1] unknown symbol |",
		},
		{
			"all",
			"qState\n;;,>>symbol1 symbol2\tBLRR,,((\n)",
			[]Token{
				{Type: StateToken, Value: "qState", Line: 1},
				{Type: SemicolonToken, Value: ";", Line: 2},
				{Type: SemicolonToken, Value: ";", Line: 2},
				{Type: CommaToken, Value: ",", Line: 2},
				{Type: ArrowToken, Value: ">", Line: 2},
				{Type: ArrowToken, Value: ">", Line: 2},
				{Type: SymbolToken, Value: "symbol1", Line: 2},
				{Type: SymbolToken, Value: "symbol2", Line: 2},
				{Type: BlankSymbolToken, Value: "B", Line: 2},
				{Type: MoveLeftToken, Value: "L", Line: 2},
				{Type: MoveRightToken, Value: "R", Line: 2},
				{Type: MoveRightToken, Value: "R", Line: 2},
				{Type: CommaToken, Value: ",", Line: 2},
				{Type: CommaToken, Value: ",", Line: 2},
				{Type: LeftParenToken, Value: "(", Line: 2},
				{Type: LeftParenToken, Value: "(", Line: 2},
				{Type: RightParenToken, Value: ")", Line: 3},
				{Type: EOFToken, Value: "", Line: 3},
			},
			"",
		},
		{
			"with comments",
			"qState\n#this line should be skipped\n\n#comment number 1\n### comment number 2\nqState",
			[]Token{
				{Type: StateToken, Value: "qState", Line: 1},
				{Type: StateToken, Value: "qState", Line: 6},
				{Type: EOFToken, Value: "", Line: 6},
			},
			"",
		},
		{
			"with comments at the end with no new line after last comment",
			"qState\n#this line should be skipped\n#this should also be skipped",
			[]Token{
				{Type: StateToken, Value: "qState", Line: 1},
				{Type: EOFToken, Value: "", Line: 3},
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
