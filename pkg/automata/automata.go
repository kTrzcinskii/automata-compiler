package automata

import "context"

type Automata interface {
	Run(ctx context.Context) (AutomataResult, error)
}

type AutomataResult interface {
	String() string
}

type State struct {
	Name      string
	Accepting bool
}

type Symbol struct {
	Name string
	// TODO: I think we need more fields here (may be in different kinds of automata?)
}

var BlankSymbol = Symbol{Name: "B"}

type TapeMoveType int

const (
	_ TapeMoveType = iota
	TapeMoveLeft
	TapeMoveRight
)
