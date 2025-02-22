package lexer

type TokenType int

const (
	_ TokenType = iota
	// Used in all automata
	CommaToken
	ArrowToken
	LeftParenToken
	RightParenToken
	SemicolonToken
	StateToken
	SymbolToken
	EOFToken

	// Used in TM
	BlankSymbolToken
	MoveLeftToken
	MoveRightToken

	// Used in PA
	InputEndToken
	StackStartToken
)

func (tt TokenType) String() string {
	switch tt {
	case CommaToken:
		return "CommaToken"
	case ArrowToken:
		return "ArrowToken"
	case LeftParenToken:
		return "LeftParenToken"
	case RightParenToken:
		return "RightParenToken"
	case SemicolonToken:
		return "SemicolonToken"
	case StateToken:
		return "StateToken"
	case SymbolToken:
		return "SymbolToken"
	case BlankSymbolToken:
		return "BlankSymbolToken"
	case MoveLeftToken:
		return "MoveLeftToken"
	case MoveRightToken:
		return "MoveRightToken"
	case EOFToken:
		return "EOFToken"
	case InputEndToken:
		return "InputEndToken"
	case StackStartToken:
		return "StackStartToken"
	default:
		return "Invalid Token Type"
	}
}

type Token struct {
	Type  TokenType
	Value string
	Line  int
}
