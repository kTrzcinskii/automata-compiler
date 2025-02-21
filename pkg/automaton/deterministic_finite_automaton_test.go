package automaton

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRunDFA(t *testing.T) {
	var zero AutomatonResult
	data := []struct {
		name           string
		dfa            *DeterministicFiniteAutomaton
		expected       AutomatonResult
		expectedErrMsg string
	}{
		{
			"empty input with default accepting state",
			&DeterministicFiniteAutomaton{
				States: map[string]State{
					"qOK": {
						Name:      "qOk",
						Accepting: true,
					},
				},
				Symbols:      map[string]Symbol{},
				CurrentState: "qOK",
				Input:        []string{},
				InputIt:      0,
				Transitions:  map[DFATransitionKey]DFATransitionValue{},
			},
			DeterministicFiniteAutomatonResult{
				FinalState: State{
					Name:      "qOk",
					Accepting: true,
				},
			},
			"",
		},
		{
			"single accepting state",
			&DeterministicFiniteAutomaton{
				States: map[string]State{
					"qOK": {
						Name:      "qOk",
						Accepting: true,
					},
				},
				Symbols: map[string]Symbol{
					"S": {
						Name: "S",
					},
				},
				CurrentState: "qOK",
				Input: []string{
					"S",
					"S",
					"S",
					"S",
					"S",
				},
				InputIt: 0,
				Transitions: map[DFATransitionKey]DFATransitionValue{
					{
						StateName:  "qOK",
						SymbolName: "S",
					}: {
						StateName: "qOK",
					},
				},
			},
			DeterministicFiniteAutomatonResult{
				FinalState: State{
					Name:      "qOk",
					Accepting: true,
				},
			},
			"",
		},
		{
			"single rejecting state",
			&DeterministicFiniteAutomaton{
				States: map[string]State{
					"qNotOK": {
						Name:      "qNotOk",
						Accepting: false,
					},
				},
				Symbols: map[string]Symbol{
					"S": {
						Name: "S",
					},
				},
				CurrentState: "qNotOK",
				Input: []string{
					"S",
					"S",
					"S",
					"S",
					"S",
				},
				InputIt: 0,
				Transitions: map[DFATransitionKey]DFATransitionValue{
					{
						StateName:  "qNotOK",
						SymbolName: "S",
					}: {
						StateName: "qNotOK",
					},
				},
			},
			DeterministicFiniteAutomatonResult{
				FinalState: State{
					Name:      "qNotOk",
					Accepting: false,
				},
			},
			"",
		},
		{
			"with few transitions",
			&DeterministicFiniteAutomaton{
				States: map[string]State{
					"qA": {
						Name:      "qA",
						Accepting: false,
					},
					"qB": {
						Name:      "qB",
						Accepting: false,
					},
					"qC": {
						Name:      "qC",
						Accepting: false,
					},
					"qAcc": {
						Name:      "qAcc",
						Accepting: true,
					},
				},
				Symbols: map[string]Symbol{
					"A": {
						Name: "A",
					},
					"B": {
						Name: "B",
					},
					"C": {
						Name: "C",
					},
				},
				CurrentState: "qA",
				Input: []string{
					"A",
					"B",
					"C",
				},
				InputIt: 0,
				Transitions: map[DFATransitionKey]DFATransitionValue{
					{
						StateName:  "qA",
						SymbolName: "A",
					}: {
						StateName: "qB",
					},
					{
						StateName:  "qB",
						SymbolName: "B",
					}: {
						StateName: "qC",
					},
					{
						StateName:  "qC",
						SymbolName: "C",
					}: {
						StateName: "qAcc",
					},
				},
			},
			DeterministicFiniteAutomatonResult{
				FinalState: State{
					Name:      "qAcc",
					Accepting: true,
				},
			},
			"",
		},
		{
			"missing transition",
			&DeterministicFiniteAutomaton{
				States: map[string]State{
					"qOK": {
						Name:      "qOK",
						Accepting: true,
					},
				},
				Symbols: map[string]Symbol{
					"P": {
						Name: "P",
					},
				},
				CurrentState: "qOK",
				Input: []string{
					"P",
				},
				InputIt:     0,
				Transitions: map[DFATransitionKey]DFATransitionValue{},
			},
			zero,
			"cannot continue calculations, missing transition for state qOK and symbol P",
		},
		{
			"long input",
			&DeterministicFiniteAutomaton{
				States: map[string]State{
					"qA": {
						Name:      "qA",
						Accepting: false,
					},
					"qAcc": {
						Name:      "qAcc",
						Accepting: true,
					},
				},
				Symbols: map[string]Symbol{
					"A": {
						Name: "A",
					},
					"B": {
						Name: "B",
					},
				},
				CurrentState: "qA",
				Input: []string{
					"A",
					"A",
					"A",
					"A",
					"A",
					"A",
					"A",
					"A",
					"A",
					"A",
					"A",
					"A",
					"A",
					"A",
					"A",
					"B",
				},
				InputIt: 0,
				Transitions: map[DFATransitionKey]DFATransitionValue{
					{
						StateName:  "qA",
						SymbolName: "A",
					}: {
						StateName: "qA",
					},
					{
						StateName:  "qA",
						SymbolName: "B",
					}: {
						StateName: "qAcc",
					},
				},
			},
			DeterministicFiniteAutomatonResult{
				FinalState: State{
					Name:      "qAcc",
					Accepting: true,
				},
			},
			"",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result, err := Run(context.Background(), d.dfa, AutomatonOptions{Output: io.Discard})
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

func TestRunWithIncludedCalculationsDFA(t *testing.T) {
	dfa := &DeterministicFiniteAutomaton{
		States: map[string]State{
			"qA": {
				Name:      "qA",
				Accepting: false,
			},
			"qB": {
				Name:      "qB",
				Accepting: false,
			},
			"qAcc": {
				Name:      "qAcc",
				Accepting: true,
			},
		},
		Symbols: map[string]Symbol{
			"A": {
				Name: "A",
			},
			"B": {
				Name: "B",
			},
		},
		CurrentState: "qA",
		Input: []string{
			"A",
			"B",
		},
		InputIt: 0,
		Transitions: map[DFATransitionKey]DFATransitionValue{
			{
				StateName:  "qA",
				SymbolName: "A",
			}: {
				StateName: "qB",
			},
			{
				StateName:  "qB",
				SymbolName: "B",
			}: {
				StateName: "qAcc",
			},
		},
	}
	expected := DeterministicFiniteAutomatonResult{
		FinalState: State{
			Name:      "qAcc",
			Accepting: true,
		},
	}
	expectedCalculations := "current state: qA, input left: A|B\ncurrent state: qB, input left: B\ncurrent state: qAcc, input left: \n"
	sb := &strings.Builder{}
	opts := AutomatonOptions{
		Output:              sb,
		IncludeCalculations: true,
	}
	result, err := Run(context.Background(), dfa, opts)
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

func TestSaveStateDFA(t *testing.T) {
	dfac := DeterministicFiniteAutomatonCurrentCalculationsState{
		State: State{
			Name: "qCurrent",
		},
		InputLeft: []Symbol{
			{Name: "A"},
			{Name: "B"},
		},
	}
	var result strings.Builder
	dfac.SaveState(&result)
	expected := "current state: qCurrent, input left: A|B\n"
	if result.String() != expected {
		t.Errorf("invalid result string, expected:\n%s, got:\n%s", expected, result.String())
	}
}

func TestSaveResultDFA(t *testing.T) {
	data := []struct {
		name     string
		dfar     DeterministicFiniteAutomatonResult
		expected string
	}{
		{
			"with rejecting state",
			DeterministicFiniteAutomatonResult{
				FinalState: State{
					Name: "qAcc",
				},
			},
			"final state: qAcc, accepted: false\n",
		},
		{
			"with accepting state",
			DeterministicFiniteAutomatonResult{
				FinalState: State{
					Name:      "qAcc",
					Accepting: true,
				},
			},
			"final state: qAcc, accepted: true\n",
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			var result strings.Builder
			d.dfar.SaveResult(&result)
			if result.String() != d.expected {
				t.Errorf("invalid result string, expected:\n%s, got:\n%s", d.expected, result.String())
			}
		})
	}
}
