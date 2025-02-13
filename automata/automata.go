package automata

type Automata interface {
	// TODO:
}

type State struct {
	Name      string
	Accepting bool
}

type Symbol struct {
	Name string
	// TODO: I think we need more fields here (may be in different kinds of automata?)
}

type TapeMoveType int

const (
	_ TapeMoveType = iota
	TapeMoveLeft
	TapeMoveRight
)
