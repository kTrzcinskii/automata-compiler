package automaton

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestSaveResult(t *testing.T) {
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
	var result strings.Builder
	tmr.SaveResult(&result)
	expected := "final state: qState, tape: s1|s2|B\n"
	if result.String() != expected {
		t.Errorf("invalid result string, expected:\n%s, got:\n%s", expected, result.String())
	}
}

func TestSaveState(t *testing.T) {
	tmc := TuringMachineCurrentCalculationsState{
		State: State{
			Name: "qState",
		},
		Tape: []Symbol{
			{Name: "s1"},
			{Name: "s2"},
			{Name: BlankSymbol.Name},
		},
		It: 1,
	}
	var result strings.Builder
	tmc.SaveState(&result)
	expected := "current state: qState, tape: s1|s2|B\n"
	id := strings.Index(expected, "s2")
	spaces := strings.Repeat(" ", id)
	expected += spaces + "^\n"
	if result.String() != expected {
		t.Errorf("invalid result string, expected:\n%s, got:\n%s", expected, result.String())
	}
}

func TestRun(t *testing.T) {
	var zero AutomatonResult
	data := []struct {
		name           string
		tm             *TuringMachine
		timeout        int
		options        AutomatonOptions
		expected       AutomatonResult
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
			AutomatonOptions{Output: io.Discard},
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
			AutomatonOptions{Output: io.Discard},
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
			AutomatonOptions{Output: io.Discard},
			TuringMachineResult{
				FinalState: State{Name: "qState2", Accepting: true},
				FinalTape: []Symbol{
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
			AutomatonOptions{Output: io.Discard},
			zero,
			"cannot continue calculations, missing transition for state qState and symbol symbol1",
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
			AutomatonOptions{Output: io.Discard},
			zero,
			"cannot continue calculations, turing machine went out of tape",
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
			AutomatonOptions{Output: io.Discard},
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
			result, err := Run(ctx, d.tm, d.options)
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

func TestRunWithIncludedCalculations(t *testing.T) {
	tm := &TuringMachine{
		States: map[string]State{
			"qState":  {Name: "qState"},
			"qState2": {Name: "qState2", Accepting: true},
		},
		CurrentState: "qState",
		Symbols: map[string]Symbol{
			BlankSymbol.Name: BlankSymbol,
		},
		Transitions: map[TMTransitionKey]TMTransitionValue{
			{StateName: "qState", SymbolName: BlankSymbol.Name}: {StateName: "qState2", SymbolName: BlankSymbol.Name, Move: TapeMoveRight},
		},
		Tape: []string{
			BlankSymbol.Name,
		},
		TapeIt: 0,
	}
	sb := &strings.Builder{}
	opts := AutomatonOptions{
		Output:              sb,
		IncludeCalculations: true,
	}
	result, err := Run(context.Background(), tm, opts)
	expected := TuringMachineResult{
		FinalState: State{Name: "qState2", Accepting: true},
		FinalTape: []Symbol{
			{Name: BlankSymbol.Name},
		},
	}
	expectedCalculations := "current state: qState, tape: B\n"
	l := len(expectedCalculations)
	expectedCalculations += strings.Repeat(" ", l-2) + "^\n"
	secondLine := "current state: qState2, tape: B|B\n"
	l2 := len(secondLine)
	secondLine += strings.Repeat(" ", l2-2) + "^\n"
	expectedCalculations += secondLine
	expectedErrMsg := ""
	if diff := cmp.Diff(expected, result); diff != "" {
		t.Error(diff)
	}
	if expectedCalculations != sb.String() {
		t.Errorf("invalid calculations, expected:\n%s, got:\n%s", expectedCalculations, sb.String())
	}
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	if errMsg != expectedErrMsg {
		t.Errorf("invalid error message, expected: %s, got: %s", expectedErrMsg, errMsg)
	}
}

func TestRunWithCalculationsOutputToFile(t *testing.T) {
	output, err := os.CreateTemp(t.TempDir(), "temp-output")
	if err != nil {
		t.Fatal(err)
	}
	defer output.Close()
	tm := &TuringMachine{
		States: map[string]State{
			"qState":  {Name: "qState"},
			"qState2": {Name: "qState2", Accepting: true},
		},
		CurrentState: "qState",
		Symbols: map[string]Symbol{
			BlankSymbol.Name: BlankSymbol,
		},
		Transitions: map[TMTransitionKey]TMTransitionValue{
			{StateName: "qState", SymbolName: BlankSymbol.Name}: {StateName: "qState2", SymbolName: BlankSymbol.Name, Move: TapeMoveRight},
		},
		Tape: []string{
			BlankSymbol.Name,
		},
		TapeIt: 0,
	}
	opts := AutomatonOptions{
		Output:              output,
		IncludeCalculations: true,
	}
	_, err = Run(context.Background(), tm, opts)
	expectedFileContent := "current state: qState, tape: B\n"
	l := len(expectedFileContent)
	expectedFileContent += strings.Repeat(" ", l-2) + "^\n"
	secondLine := "current state: qState2, tape: B|B\n"
	l2 := len(secondLine)
	secondLine += strings.Repeat(" ", l2-2) + "^\n"
	expectedFileContent += secondLine
	expectedErrMsg := ""
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	if errMsg != expectedErrMsg {
		t.Errorf("invalid error message, expected: %s, got: %s", expectedErrMsg, errMsg)
	}
	if err := output.Sync(); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(output.Name())
	if err != nil {
		t.Fatal(err)
	}
	outputContent := string(b)
	if outputContent != expectedFileContent {
		t.Errorf("invalid file content, expected:\n%s, got:\n%s", expectedFileContent, outputContent)
	}
}
