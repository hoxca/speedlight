/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"speedlight/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	Log "github.com/apatters/go-conlog"
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
	cfgFile         string
	cfgFileNotFound = false
	lightsdir       string
	writeConsole    bool
	writeReport     bool
	rotation        bool
	verbosity       string
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is conf/speedlight.yaml)")
	rootCmd.PersistentFlags().BoolVar(&writeConsole, "console", true, "write report to the console")
	rootCmd.PersistentFlags().BoolVar(&writeReport, "report", true, "write report to the filesystem")
	rootCmd.PersistentFlags().StringVar(&lightsdir, "dir", "", "lights directory")
	rootCmd.PersistentFlags().BoolVar(&rotation, "rotation", true, "manage rotation in lights report")
	rootCmd.PersistentFlags().StringVar(&verbosity, "level", "", "set log level")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Switch to default program path
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}

		confdir := fmt.Sprintf("%s/conf", dir)
		// we came from bin directory
		confdir1 := fmt.Sprintf("%s/../conf", dir)
		confdir2 := "./conf"
		// Search yaml config file in program path with name "speedlight.yaml".
		viper.AddConfigPath(confdir)
		viper.AddConfigPath(confdir1)
		viper.AddConfigPath(confdir2)
		viper.AddConfigPath(dir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("speedlight")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			cfgFileNotFound = true
			fmt.Println("Config file not found")
		} else {
			Log.Debug("Something look strange")
			Log.Debugf("error: %v\n", err)
		}
	} else {
		Log.Debugf("Using config file: %s\n", viper.ConfigFileUsed())
	}

	manageDefault()
}

func manageDefault() {

	if len(lightsdir) == 0 {
		lightsdir = viper.GetString("lightsdir")
	}
	if len(verbosity) == 0 {
		verbosity = viper.GetString("level")
	}
	utils.TimeFrame = viper.GetInt("time_frame")
	utils.Regex = viper.GetString("regexp")
	Log.Debugf("regex: %s\n", utils.Regex)

}
