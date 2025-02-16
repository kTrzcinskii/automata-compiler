package automata

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestString(t *testing.T) {
	tmr := TuringMachineResult{
		FinalState: State{
			Name:      "qState",
			Accepting: true,
		},
		FinalTape: []Symbol{
			{Name: "s1"},
			{Name: "s2"},
			{Name: BlankSymbol.Name},
		},
	}
	result := tmr.String()
	expected := "state: qState, tape: s1|s2|B"
	if result != expected {
		t.Errorf("invalid result string, expected: %s, got: %s", expected, result)
	}
}

func TestRun(t *testing.T) {
	var zero TuringMachineResult
	data := []struct {
		name           string
		tm             *TuringMachine
		timeout        int
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
					BlankSymbol.Name: BlankSymbol,
					"symbol1":        {Name: "symbol1"},
				},
				Transitions: map[TMTransitionKey]TMTransitionValue{},
				Tape: []string{
					"symbol1",
				},
				TapeIt: 0,
			},
			0,
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
					BlankSymbol.Name: BlankSymbol,
					"symbol1":        {Name: "symbol1"},
					"symbol2":        {Name: "symbol2"},
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
			0,
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
		{
			"add new tape symbols when going right",
			&TuringMachine{
				States: map[string]State{
					"qState":  {Name: "qState"},
					"qStat1":  {Name: "qState1"},
					"qState2": {Name: "qState2", Accepting: true},
				},
				CurrentState: "qState",
				Symbols: map[string]Symbol{
					BlankSymbol.Name: BlankSymbol,
				},
				Transitions: map[TMTransitionKey]TMTransitionValue{
					{StateName: "qState", SymbolName: BlankSymbol.Name}:  {StateName: "qState1", SymbolName: BlankSymbol.Name, Move: TapeMoveRight},
					{StateName: "qState1", SymbolName: BlankSymbol.Name}: {StateName: "qState2", SymbolName: BlankSymbol.Name, Move: TapeMoveRight},
				},
				Tape: []string{
					"B",
				},
				TapeIt: 0,
			},
			0,
			TuringMachineResult{
				FinalState: State{Name: "qState2", Accepting: true},
				FinalTape: []Symbol{
					BlankSymbol,
					BlankSymbol,
					BlankSymbol,
				},
			},
			"",
		},
		{
			"missing transition",
			&TuringMachine{
				States: map[string]State{
					"qState": {Name: "qState"},
				},
				CurrentState: "qState",
				Symbols: map[string]Symbol{
					BlankSymbol.Name: BlankSymbol,
					"symbol1":        {Name: "symbol1"},
				},
				Transitions: map[TMTransitionKey]TMTransitionValue{},
				Tape: []string{
					"symbol1",
				},
				TapeIt: 0,
			},
			0,
			zero,
			"cannot continue calucations, missing transition for state qState and symbol symbol1",
		},
		{
			"go out of tape",
			&TuringMachine{
				States: map[string]State{
					"qState": {Name: "qState"},
				},
				CurrentState: "qState",
				Symbols: map[string]Symbol{
					BlankSymbol.Name: BlankSymbol,
					"symbol1":        {Name: "symbol1"},
				},
				Transitions: map[TMTransitionKey]TMTransitionValue{
					{StateName: "qState", SymbolName: "symbol1"}: {StateName: "qState", SymbolName: BlankSymbol.Name, Move: TapeMoveLeft},
				},
				Tape: []string{
					"symbol1",
				},
				TapeIt: 0,
			},
			0,
			zero,
			"cannot continue calucations, turing machine went out of tape",
		},
		{
			"infinite loop with timeout",
			&TuringMachine{
				States: map[string]State{
					"q0": {Name: "q0"},
					"q1": {Name: "q1"},
				},
				CurrentState: "q0",
				Symbols: map[string]Symbol{
					BlankSymbol.Name: BlankSymbol,
				},
				Transitions: map[TMTransitionKey]TMTransitionValue{
					{StateName: "q0", SymbolName: BlankSymbol.Name}: {StateName: "q1", SymbolName: BlankSymbol.Name, Move: TapeMoveRight},
					{StateName: "q1", SymbolName: BlankSymbol.Name}: {StateName: "q0", SymbolName: BlankSymbol.Name, Move: TapeMoveLeft},
				},
				Tape: []string{
					BlankSymbol.Name,
				},
				TapeIt: 0,
			},
			500,
			zero,
			"timeout reached",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			var ctx context.Context
			if d.timeout == 0 {
				ctx = context.Background()
			} else {
				ctxWithT, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(d.timeout))
				ctx = ctxWithT
				defer cancelFunc()
			}
			result, err := d.tm.Run(ctx)
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
