package automaton

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRunPA(t *testing.T) {
	var zero AutomatonResult
	data := []struct {
		name           string
		pa             *PushdownAutomaton
		expected       AutomatonResult
		expectedErrMsg string
	}{
		{
			"empty input with accepting state",
			&PushdownAutomaton{
				States: map[string]State{
					"qState": {Name: "qState", Accepting: true},
				},
				CurrentState: "qState",
				Symbols: map[string]Symbol{
					InputEndSymbol.Name:   InputEndSymbol,
					StackStartSymbol.Name: StackStartSymbol,
				},
				Input:   []string{InputEndSymbol.Name},
				InputIt: 0,
				Stack:   []string{StackStartSymbol.Name},
				Transitions: map[PATransitionKey]PATransitionValue{
					{
						StateName:       "qState",
						InputSymbolName: InputEndSymbol.Name,
						StackSymbolName: StackStartSymbol.Name,
					}: {
						StateName:        "qState",
						StackSymbolNames: []string{StackStartSymbol.Name},
					},
				},
			},
			PushdownAutomatonResult{
				FinalState: State{
					Name:      "qState",
					Accepting: true,
				},
				Stack: []Symbol{StackStartSymbol},
			},
			"",
		},
		{
			"few transitions without using stack",
			&PushdownAutomaton{
				States: map[string]State{
					"qA":   {Name: "qA"},
					"qB":   {Name: "qB"},
					"qC":   {Name: "qC"},
					"qEnd": {Name: "qEnd"},
					"qAcc": {Name: "qAcc", Accepting: true},
				},
				CurrentState: "qA",
				Symbols: map[string]Symbol{
					InputEndSymbol.Name:   InputEndSymbol,
					StackStartSymbol.Name: StackStartSymbol,
					"A":                   {Name: "A"},
					"B":                   {Name: "B"},
					"C":                   {Name: "C"},
				},
				Input:   []string{"A", "B", "C", InputEndSymbol.Name},
				InputIt: 0,
				Stack:   []string{StackStartSymbol.Name},
				Transitions: map[PATransitionKey]PATransitionValue{
					{
						StateName:       "qA",
						InputSymbolName: "A",
						StackSymbolName: StackStartSymbol.Name,
					}: {
						StateName:        "qB",
						StackSymbolNames: []string{StackStartSymbol.Name},
					},
					{
						StateName:       "qB",
						InputSymbolName: "B",
						StackSymbolName: StackStartSymbol.Name,
					}: {
						StateName:        "qC",
						StackSymbolNames: []string{StackStartSymbol.Name},
					},
					{
						StateName:       "qC",
						InputSymbolName: "C",
						StackSymbolName: StackStartSymbol.Name,
					}: {
						StateName:        "qEnd",
						StackSymbolNames: []string{StackStartSymbol.Name},
					},
					{
						StateName:       "qEnd",
						InputSymbolName: InputEndSymbol.Name,
						StackSymbolName: StackStartSymbol.Name,
					}: {
						StateName:        "qAcc",
						StackSymbolNames: []string{StackStartSymbol.Name},
					},
				},
			},
			PushdownAutomatonResult{
				FinalState: State{
					Name:      "qAcc",
					Accepting: true,
				},
				Stack: []Symbol{StackStartSymbol},
			},
			"",
		},
		{
			"few transitions with changing stack",
			&PushdownAutomaton{
				States: map[string]State{
					"qState": {Name: "qState", Accepting: true},
				},
				CurrentState: "qState",
				Symbols: map[string]Symbol{
					InputEndSymbol.Name:   InputEndSymbol,
					StackStartSymbol.Name: StackStartSymbol,
					"A":                   {Name: "A"},
				},
				Input:   []string{"A", InputEndSymbol.Name},
				InputIt: 0,
				Stack:   []string{StackStartSymbol.Name},
				Transitions: map[PATransitionKey]PATransitionValue{
					{
						StateName:       "qState",
						InputSymbolName: "A",
						StackSymbolName: StackStartSymbol.Name,
					}: {
						StateName:        "qState",
						StackSymbolNames: []string{StackStartSymbol.Name, "A", "A", "A", "A"},
					},
					{
						StateName:       "qState",
						InputSymbolName: InputEndSymbol.Name,
						StackSymbolName: "A",
					}: {
						StateName:        "qState",
						StackSymbolNames: []string{},
					},
				},
			},
			PushdownAutomatonResult{
				FinalState: State{
					Name:      "qState",
					Accepting: true,
				},
				Stack: []Symbol{StackStartSymbol, {Name: "A"}, {Name: "A"}, {Name: "A"}},
			},
			"",
		},
		{
			"missing transition",
			&PushdownAutomaton{
				States: map[string]State{
					"qA": {Name: "qA"},
					"qB": {Name: "qB"},
				},
				CurrentState: "qA",
				Symbols: map[string]Symbol{
					InputEndSymbol.Name:   InputEndSymbol,
					StackStartSymbol.Name: StackStartSymbol,
					"A":                   {Name: "A"},
					"B":                   {Name: "B"},
				},
				Input:   []string{"A", InputEndSymbol.Name},
				InputIt: 0,
				Stack:   []string{StackStartSymbol.Name},
				Transitions: map[PATransitionKey]PATransitionValue{
					{
						StateName:       "qA",
						InputSymbolName: "A",
						StackSymbolName: StackStartSymbol.Name,
					}: {
						StateName:        "qB",
						StackSymbolNames: []string{StackStartSymbol.Name, "B"},
					},
					{
						StateName:       "qB",
						InputSymbolName: InputEndSymbol.Name,
						StackSymbolName: StackStartSymbol.Name,
					}: {
						StateName:        "qB",
						StackSymbolNames: []string{StackStartSymbol.Name},
					},
				},
			},
			zero,
			"cannot continue calculations, missing transition for state qB, symbol { and stack symbol B",
		},
		{
			"transition that makes stack empty",
			&PushdownAutomaton{
				States: map[string]State{
					"qA": {Name: "qA"},
				},
				CurrentState: "qA",
				Symbols: map[string]Symbol{
					InputEndSymbol.Name:   InputEndSymbol,
					StackStartSymbol.Name: StackStartSymbol,
					"A":                   {Name: "A"},
				},
				Input:   []string{"A", InputEndSymbol.Name},
				InputIt: 0,
				Stack:   []string{StackStartSymbol.Name},
				Transitions: map[PATransitionKey]PATransitionValue{
					{
						StateName:       "qA",
						InputSymbolName: "A",
						StackSymbolName: StackStartSymbol.Name,
					}: {
						StateName:        "qA",
						StackSymbolNames: []string{},
					},
					{
						StateName:       "qA",
						InputSymbolName: InputEndSymbol.Name,
						StackSymbolName: StackStartSymbol.Name,
					}: {
						StateName:        "qA",
						StackSymbolNames: []string{StackStartSymbol.Name},
					},
				},
			},
			zero,
			"stack is empty",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			ctx := context.Background()
			opts := AutomatonOptions{
				Output:              io.Discard,
				IncludeCalculations: false,
			}
			result, err := Run(ctx, d.pa, opts)
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

func TestRunWithIncludedCalculationsPA(t *testing.T) {
	pa := &PushdownAutomaton{
		States: map[string]State{
			"qStart": {Name: "qStart"},
			"qA":     {Name: "qA"},
			"qB":     {Name: "qB"},
			"qEnd":   {Name: "qEnd"},
			"qAcc":   {Name: "qAcc", Accepting: true},
		},
		CurrentState: "qStart",
		Symbols: map[string]Symbol{
			InputEndSymbol.Name:   InputEndSymbol,
			StackStartSymbol.Name: StackStartSymbol,
			"A":                   {Name: "A"},
			"B":                   {Name: "B"},
			"X":                   {Name: "X"},
		},
		Input:   []string{"X", "A", "B", InputEndSymbol.Name},
		InputIt: 0,
		Stack:   []string{StackStartSymbol.Name},
		Transitions: map[PATransitionKey]PATransitionValue{
			{
				StateName:       "qStart",
				InputSymbolName: "X",
				StackSymbolName: StackStartSymbol.Name,
			}: {
				StateName:        "qA",
				StackSymbolNames: []string{StackStartSymbol.Name, "A"},
			},
			{
				StateName:       "qA",
				InputSymbolName: "A",
				StackSymbolName: "A",
			}: {
				StateName:        "qB",
				StackSymbolNames: []string{"B"},
			},
			{
				StateName:       "qB",
				InputSymbolName: "B",
				StackSymbolName: "B",
			}: {
				StateName:        "qEnd",
				StackSymbolNames: []string{},
			},
			{
				StateName:       "qEnd",
				InputSymbolName: InputEndSymbol.Name,
				StackSymbolName: StackStartSymbol.Name,
			}: {
				StateName:        "qAcc",
				StackSymbolNames: []string{StackStartSymbol.Name},
			},
		},
	}
	expected := PushdownAutomatonResult{
		FinalState: State{
			Name:      "qAcc",
			Accepting: true,
		},
		Stack: []Symbol{StackStartSymbol},
	}
	expectedCalculations := "current state: qStart, input left: X|A|B|{, stack: }\ncurrent state: qA, input left: A|B|{, stack: }|A\ncurrent state: qB, input left: B|{, stack: }|B\ncurrent state: qEnd, input left: {, stack: }\ncurrent state: qAcc, input left: , stack: }\n"
	sb := &strings.Builder{}
	opts := AutomatonOptions{
		Output:              sb,
		IncludeCalculations: true,
	}
	result, err := Run(context.Background(), pa, opts)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Error(diff)
	}
	if expectedCalculations != sb.String() {
		t.Errorf("invalid calculations, expected:\n%s, got:\n%s", expectedCalculations, sb.String())
	}

}

func TestSaveStatePA(t *testing.T) {
	aSym := Symbol{Name: "A"}
	bSym := Symbol{Name: "B"}
	pac := PushdownAutomatonCurrentCalculationsState{
		CurrentState: State{Name: "qCS", Accepting: false},
		Stack: []Symbol{
			StackStartSymbol,
			aSym,
		},
		InputLeft: []Symbol{bSym, bSym, InputEndSymbol},
	}
	var result strings.Builder
	pac.SaveState(&result)
	expected := "current state: qCS, input left: B|B|{, stack: }|A\n"
	if result.String() != expected {
		t.Errorf("invalid result string, expected:\n%s, got:\n%s", expected, result.String())
	}
}

func TestSaveResultPA(t *testing.T) {
	cSym := Symbol{Name: "C"}
	data := []struct {
		name     string
		par      PushdownAutomatonResult
		expected string
	}{
		{
			"with accepting state",
			PushdownAutomatonResult{
				FinalState: State{Name: "qAcc", Accepting: true},
				Stack: []Symbol{
					StackStartSymbol,
					cSym,
					cSym,
				},
			},
			"final state: qAcc, accepted: true, stack: }|C|C\n",
		},
		{
			"with rejecting state",
			PushdownAutomatonResult{
				FinalState: State{Name: "qAcc", Accepting: false},
				Stack: []Symbol{
					StackStartSymbol,
					cSym,
					cSym,
					cSym,
				},
			},
			"final state: qAcc, accepted: false, stack: }|C|C|C\n",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			var result strings.Builder
			d.par.SaveResult(&result)
			if result.String() != d.expected {
				t.Errorf("invalid result string, expected:\n%s, got:\n%s", d.expected, result.String())
			}
		})
	}
}
