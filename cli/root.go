package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "perfin",
	Short: "perfin cli client application",
}

func init() {
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "exit with error: %v\n", err)
		os.Exit(1)
	}
}
