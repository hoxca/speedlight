/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "speedlight",
	Short: "Is a tool to sumarize your lights directory",
	Long: `Speedlight implement 2 commands [report and flats]
which make different reports based on your local lights directory.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {},
}

var (
	lightsdir    string
	writeConsole bool
	writeReport  bool
	rotation     bool
	verbosity    string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	initConfig()
}

func initConfig() {
	rootCmd.PersistentFlags().BoolVar(&writeConsole, "console", true, "write report to the console")
	rootCmd.PersistentFlags().BoolVar(&writeReport, "report", true, "write report to the filesystem")
	rootCmd.PersistentFlags().StringVar(&lightsdir, "dir", "D:/Data/Voyager/Lights/", "lights directory")
	rootCmd.PersistentFlags().BoolVar(&rotation, "rotation", true, "manage rotation in lights report")
	rootCmd.PersistentFlags().StringVar(&verbosity, "level", "warn", "set log level")
}
