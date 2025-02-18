package cmd

import (
	"automata-compiler/pkg/automaton"
	"automata-compiler/pkg/compiler"
	"automata-compiler/pkg/lexer"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "automata-compiler PAHT_TO_INPUT_FILE",
	Short: "automata-compiler is a tool for simulating automata",
	Long: `The automata-compiler is a CLI application for compiling and running automata code.
It currenlty only supports Turing Machines, but there should be more type of
automata available in the future.`,
	RunE: runRootCmd,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
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
	b, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	source := string(b)
	l := lexer.NewLexer(source)
	tokens, err := l.ScanTokens()
	if err != nil {
		fmt.Printf("Error during lexing stage: %s\n", err.Error())
		return nil
	}
	tmc := compiler.NewTuringMachineCompiler(tokens)
	tm, err := tmc.Compile()
	if err != nil {
		fmt.Printf("Error during compiling stage: %s\n", err.Error())
		return nil
	}
	opts, cleanupFunc, err := automatonOptions(cmd)
	if err != nil {
		return err
	}
	defer cleanupFunc()
	ctx, cancelFunc, err := createCmdContext(cmd)
	if err != nil {
		return err
	}
	defer cancelFunc()
	result, err := automaton.Run(ctx, tm, opts)
	if err != nil {
		fmt.Printf("Error during running stage: %s\n", err.Error())
		return nil
	}
	result.SaveResult(opts.Output)
	return nil
}

func createCmdContext(cmd *cobra.Command) (context.Context, context.CancelFunc, error) {
	timeout, err := cmd.Flags().GetUint32(timeoutFlag.name)
	if err != nil {
		return nil, nil, err
	}
	if timeout > 0 {
		ctx, fun := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(timeout))
		return ctx, fun, nil
	}
	emptyFun := func() {}
	return context.Background(), emptyFun, nil
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
