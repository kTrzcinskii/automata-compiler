package lexer

type TokenType int

const (
	_ TokenType = iota
	CommaToken
	ArrowToken
	LeftParenToken
	RightParenToken
	SemicolonToken
	StateToken
	SymbolToken
	BlankSymbolToken
	MoveLeftToken
	MoveRightToken
	EOFToken
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
}
