package lexer

import (
	"fmt"
	"unicode"
)

type Lexerer interface {
	ScanTokens() ([]Token, error)
}

type Lexer struct {
	// Source code of the program
	source string
	// List of tokens scanned from the source code
	tokens []Token
	// Current line
	line int
	// Start (rune id) of currenlty analyzed token
	start int
	// End (rune id) of currently analyzed token (exclusive)
	current int
}

func (l *Lexer) ScanTokens() ([]Token, error) {
	for {
		l.skipWhitespaces()
		l.skipComments()
		if l.isAtEnd() {
			break
		}
		l.start = l.current
		t, err := l.scanToken()
		if err != nil {
			var zero []Token
			return zero, err
		}
		l.tokens = append(l.tokens, t)
	}
	l.tokens = append(l.tokens, Token{Type: EOFToken, Value: "", Line: l.line})
	return l.tokens, nil
}

func NewLexer(source string) *Lexer {
	return &Lexer{
		source:  source,
		tokens:  make([]Token, 0),
		line:    1,
		start:   0,
		current: 0,
	}
}

func (l Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) scanToken() (Token, error) {
	r := l.advance()
	c := string(r)

	switch c {
	case "q":
		state := l.readAlphanumeric()
		return Token{Type: StateToken, Value: state, Line: l.line}, nil
	case "(":
		return Token{Type: LeftParenToken, Value: c, Line: l.line}, nil
	case ")":
		return Token{Type: RightParenToken, Value: c, Line: l.line}, nil
	case ",":
		return Token{Type: CommaToken, Value: c, Line: l.line}, nil
	case ";":
		return Token{Type: SemicolonToken, Value: c, Line: l.line}, nil
	case ">":
		return Token{Type: ArrowToken, Value: c, Line: l.line}, nil
	case "L":
		return Token{Type: MoveLeftToken, Value: c, Line: l.line}, nil
	case "R":
		return Token{Type: MoveRightToken, Value: c, Line: l.line}, nil
	case "B":
		return Token{Type: BlankSymbolToken, Value: c, Line: l.line}, nil
	case "}":
		return Token{Type: StackStartToken, Value: c, Line: l.line}, nil
	case "{":
		return Token{Type: InputEndToken, Value: c, Line: l.line}, nil
	default:
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			symbol := l.readAlphanumeric()
			return Token{Type: SymbolToken, Value: symbol, Line: l.line}, nil
		}
		var zero Token
		return zero, fmt.Errorf("[Line %d] unknown symbol %s", l.line, c)
	}
}

func (l *Lexer) advance() rune {
	c := l.peek()
	l.current++
	return c
}

func (l Lexer) peek() rune {
	if l.isAtEnd() {
		return 0
	}
	c := []rune(l.source)[l.current]
	return c
}

// readAlphanumeric calls advance untill the next rune is not letter nor digit and returns consumed string
func (l *Lexer) readAlphanumeric() string {
	for {
		c := l.peek()
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			break
		}
		l.advance()
	}
	return l.sourceFragment(l.start, l.current)
}

// sourceFragment returns fragment of the source code, starting at rune with id `from` and ending at rune
// with id `to` (exclusive)
func (l Lexer) sourceFragment(from, to int) string {
	runes := []rune(l.source)
	fragment := runes[from:to]
	return string(fragment)
}

// skipWhitespaces calls advance until the next rune is not unicode whitespace
//
// if rune is a `\n` then it increments new lines count
func (l *Lexer) skipWhitespaces() {
	for {
		c := l.peek()
		if !unicode.IsSpace(c) {
			break
		}
		if c == '\n' {
			l.line++
		}
		l.advance()
	}
}

func (l *Lexer) skipComments() {
	for {
		c := l.peek()
		if c == '#' {
			l.skipLine()
			l.skipWhitespaces()
		} else {
			break
		}
	}
}

// skipLine calls advance until new line is reached
func (l *Lexer) skipLine() {
	for {
		c := l.advance()
		if c == '\n' {
			l.line++
			break
		}
		if l.isAtEnd() {
			break
		}
	}
}
