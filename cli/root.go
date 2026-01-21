package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zivlakmilos/perfin/cfg"
)

var rootCmd = &cobra.Command{
	Use:              "perfin",
	Short:            "perfin cli client application",
	PersistentPreRun: rootPreRun,
}

var config *cfg.CliConfig

var customCfgPath string

func init() {
	config = cfg.LoadCliConfig()

	loginCmd.PersistentFlags().StringVar(&customCfgPath, "config", "", "custom path for config.json")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "exit with error: %v\n", err)
		os.Exit(1)
	}
}

func rootPreRun(cmd *cobra.Command, argts []string) {
	if customCfgPath != "" {
		config = cfg.LoadCustomConfig(customCfgPath)
	}
}

func saveConfig() {
	if customCfgPath != "" {
		cfg.SaveCustomCliConfig(customCfgPath, config)
	} else {
		cfg.SaveCliConfig(config)
	}
}
