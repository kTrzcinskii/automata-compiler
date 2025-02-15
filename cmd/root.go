package cmd

import (
	"automata-compiler/pkg/compiler"
	"automata-compiler/pkg/lexer"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "automata-compiler PAHT_TO_INPUT_FILE",
	Short: "automata-compiler is a tool for simulating automatas",
	Long: `The automata-compiler is a CLI application for compiling and running automatas code.
It currenlty only supports Turing Machines, but there should be more type of
automatas available in the future.`,
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

const (
	TimeoutFlag      = "timeout"
	TimeoutFlagShort = "t"
)

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.Flags().Uint32P(TimeoutFlag, TimeoutFlagShort, 3000, "Timeout in miliseconds after which program will stop any remaining calculations. It's useful as many automatas can enter infinite loop for some input values. Set this value to 0 if you don't want any timeout.")
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
	ctx, cancelFunc, err := createCmdContext(cmd)
	if err != nil {
		return err
	}
	defer cancelFunc()
	result, err := tm.Run(ctx)
	if err != nil {
		fmt.Printf("Error during running stage: %s\n", err.Error())
		return nil
	}
	fmt.Printf("Calculations completed:\n\n%s\n", result.String())
	return nil
}

func createCmdContext(cmd *cobra.Command) (context.Context, context.CancelFunc, error) {
	timeout, err := cmd.Flags().GetUint32(TimeoutFlag)
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
