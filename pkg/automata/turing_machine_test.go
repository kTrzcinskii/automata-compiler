package automata

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRun(t *testing.T) {
	data := []struct {
		name           string
		tm             *TuringMachine
		expected       TuringMachineResult
		expectedErrMsg string
	}{
		{
			"initial state accepting",
			&TuringMachine{
				States: map[string]State{
					"qState": {Name: "qState", Accepting: true},
				},
				CurrentState: "qState",
				Symbols: map[string]Symbol{
					"symbol1": {Name: "symbol1"},
				},
				Transitions: map[TMTransitionKey]TMTransitionValue{},
				Tape: []string{
					"symbol1",
				},
				TapeIt: 0,
			},
			TuringMachineResult{
				FinalState: State{Name: "qState", Accepting: true},
				FinalTape: []Symbol{
					{Name: "symbol1"},
				},
			},
			"",
		},
		{
			"few simple transitions",
			&TuringMachine{
				States: map[string]State{
					"qState":  {Name: "qState"},
					"qState2": {Name: "qState2"},
					"qState3": {Name: "qState3", Accepting: true},
				},
				CurrentState: "qState",
				Symbols: map[string]Symbol{
					"symbol1": {Name: "symbol1"},
					"symbol2": {Name: "symbol2"},
				},
				Transitions: map[TMTransitionKey]TMTransitionValue{
					{StateName: "qState", SymbolName: "symbol1"}:  {StateName: "qState2", SymbolName: "symbol2", Move: TapeMoveRight},
					{StateName: "qState2", SymbolName: "symbol1"}: {StateName: "qState3", SymbolName: "symbol2", Move: TapeMoveRight},
				},
				Tape: []string{
					"symbol1",
					"symbol1",
					"symbol1",
				},
				TapeIt: 0,
			},
			TuringMachineResult{
				FinalState: State{Name: "qState3", Accepting: true},
				FinalTape: []Symbol{
					{Name: "symbol2"},
					{Name: "symbol2"},
					{Name: "symbol1"},
				},
			},
			"",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result, err := d.tm.Run()
			if diff := cmp.Diff(d.expected, result); diff != "" {
				t.Error(diff)
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
