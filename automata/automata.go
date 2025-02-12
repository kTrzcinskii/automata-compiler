package automata

type Automata interface {
	// TODO:
}

type State struct {
	Name      string
	Accepting bool
}

type Symbol struct {
	Name         string
	FromAlphabet bool
}

type TapeMoveType int

const (
	_ TapeMoveType = iota
	TapeMoveLeft
	TapeMoveRight
)
