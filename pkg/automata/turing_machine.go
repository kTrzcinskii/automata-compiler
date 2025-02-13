package automata

type TMTransitionKey struct {
	StateName  string
	SymbolName string
}

type TMTransitionValue struct {
	StateName  string
	SymbolName string
	Move       TapeMoveType
}

type TransitionFunction map[TMTransitionKey]TMTransitionValue

type TuringMachine struct {
	States       map[string]State
	Symbols      map[string]Symbol
	InitialState string
	Tape         []string
	Transitions  TransitionFunction
}
