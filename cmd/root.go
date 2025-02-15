package cmd

import (
	"automata-compiler/pkg/compiler"
	"automata-compiler/pkg/lexer"
	"fmt"
	"os"

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

func init() {
	// TODO: add flag for timeout
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.automata-compiler.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
		return fmt.Errorf("error during lexing stage: %s", err.Error())
	}
	tmc := compiler.NewTuringMachineCompiler(tokens)
	tm, err := tmc.Compile()
	if err != nil {
		return fmt.Errorf("error during compiling stage: %s", err.Error())
	}
	result, err := tm.Run()
	if err != nil {
		// We don't return error here as it's not an actual program error, but rather
		// just not accepting state of machine
		fmt.Printf("Error during running stage: %s\n", err.Error())
		return nil
	}
	fmt.Printf("Calculations completed:\n\n%s\n", result.String())
	return nil
}
