package cmd

import (
	"automata-compiler/pkg/automaton"
	"automata-compiler/pkg/compiler"
	"automata-compiler/pkg/lexer"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "automata-compiler AUTOMATON_TYPE PAHT_TO_INPUT_FILE",
	Short: "automata-compiler is a tool for simulating automata",
	Long: `The automata-compiler is a CLI application for compiling and running automata code.
It supports Deterministic Finite Automata, Pushdown Automaton and Turing Machines. 
AUTOMATON_TYPE is one of the following
- DFA (for Deterministic Finite Automaton)
- PA (for Pushdown Automaton)
- TM (for Turing Machine)`,
	RunE: runRootCmd,
	Args: cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

type flag struct {
	name  string
	short string
}

var (
	timeoutFlag         = flag{name: "timeout", short: "t"}
	output              = flag{name: "output", short: "o"}
	includeCalculations = flag{name: "include-calculations", short: "i"}
)

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.Flags().Uint32P(timeoutFlag.name, timeoutFlag.short, 3000, "Timeout in miliseconds after which program will stop any remaining calculations. It's useful as many automata can enter infinite loop for some input values. Set this value to 0 if you don't want any timeout.")
	rootCmd.Flags().StringP(output.name, output.short, "", "Use this flag to specify filepath where output should be placed. If you want to use `stdout` leave this option empty.")
	rootCmd.Flags().BoolP(includeCalculations.name, includeCalculations.short, false, "If set to true all calculations done by automaton will be written to output.")
}

func runRootCmd(cmd *cobra.Command, args []string) error {
	// automaton type
	aType := args[0]

	// source
	b, err := os.ReadFile(args[1])
	if err != nil {
		return err
	}
	source := string(b)

	// parse automaton options
	opts, cleanupFunc, err := automatonOptions(cmd)
	if err != nil {
		return err
	}
	defer cleanupFunc()

	// parse timeout
	timeout, err := cmd.Flags().GetUint32(timeoutFlag.name)
	if err != nil {
		return err
	}

	// start processing
	err = processAutomaton(aType, source, opts, timeout)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	return nil
}

func automatonOptions(cmd *cobra.Command) (automaton.AutomatonOptions, func(), error) {
	opts := automaton.AutomatonOptions{}
	cleanupFunc := func() {}
	// output
	output, err := cmd.Flags().GetString(output.name)
	if err != nil {
		return opts, nil, err
	}
	if output == "" {
		opts.Output = os.Stdout
	} else {
		err := os.MkdirAll(filepath.Dir(output), 0777)
		if err != nil {
			return opts, nil, err
		}
		f, err := os.Create(output)
		if err != nil {
			return opts, nil, err
		}
		opts.Output = f
		cleanupFunc = func() {
			f.Close()
		}
	}
	// include calculations
	ic, err := cmd.Flags().GetBool(includeCalculations.name)
	if err != nil {
		return opts, nil, err
	}
	opts.IncludeCalculations = ic
	return opts, cleanupFunc, nil
}

func getCompiler(tokens []lexer.Token, aType string) (compiler.Compiler, error) {
	switch strings.ToLower(aType) {
	case "dfa":
		return compiler.NewDeterministicFiniteAutomatonCompiler(tokens), nil
	case "tm":
		return compiler.NewTuringMachineCompiler(tokens), nil
	case "pa":
		return compiler.NewPushdownAutomatonCompiler(tokens), nil
	default:
		return nil, fmt.Errorf("unsupported automaton type: '%s'", aType)
	}
}

func createContextWithTimeout(timeout uint32) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		ctx, fun := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeout))
		return ctx, fun
	}
	emptyFun := func() {}
	return context.Background(), emptyFun
}

func processAutomaton(aType string, source string, opts automaton.AutomatonOptions, timeout uint32) error {
	l := lexer.NewLexer(source)
	tokens, err := l.ScanTokens()
	if err != nil {
		return fmt.Errorf("error during lexing stage: %s", err.Error())
	}
	c, err := getCompiler(tokens, aType)
	if err != nil {
		return err
	}
	a, err := c.Compile()
	if err != nil {
		return fmt.Errorf("error during compiling stage: %s", err.Error())
	}
	ctx, cancelFunc := createContextWithTimeout(timeout)
	defer cancelFunc()
	result, err := automaton.Run(ctx, a, opts)
	if err != nil {
		return fmt.Errorf("error during running stage: %s", err.Error())
	}
	result.SaveResult(opts.Output)
	return nil
}
