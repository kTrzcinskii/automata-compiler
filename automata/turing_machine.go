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

type TuringMachine struct {
	States       map[string]State
	Symbols      map[string]Symbol
	InitialState string
	Tape         []Symbol
	Transitions  map[TMTransitionKey]TMTransitionValue
}
