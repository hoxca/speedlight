package cmd

import (
	"log"
	"path/filepath"
	"speedlight/utils"

	"github.com/spf13/cobra"
)

var wdest utils.WriteDestination

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "report is flash view on your lights directortory",
	Long: `report will generate a report on all the lights
produced by your voyager astronomy orchestrator

it will sum time exposure by target and temperature.`,
	Run: func(cmd *cobra.Command, args []string) {

		utils.SetUpLogs(verbosity)
		utils.Wdest = utils.WriteDestination{writeConsole, writeReport}

		utils.RotUsed, _ = cmd.Flags().GetBool("rotation")
		err := filepath.Walk(lightsdir, utils.Traversal)
		if err != nil {
			log.Fatal(err)
		}
		utils.ObjectList.PrintObjects(lightsdir)
	},
}

var (
	lightsdir    string
	writeConsole bool
	writeReport  bool
	rotation     bool
	verbosity    string
)

func init() {

	rootCmd.AddCommand(reportCmd)

	reportCmd.PersistentFlags().StringVar(&lightsdir, "dir", "D:/Data/Voyager/Lights/", "lights directory")
	reportCmd.Flags().BoolVar(&writeConsole, "console", true, "write report to the console")
	reportCmd.Flags().BoolVar(&writeReport, "report", false, "write report to the filesystem")
	reportCmd.PersistentFlags().BoolVar(&rotation, "rotation", false, "manage rotation in lights report")
	reportCmd.PersistentFlags().StringVar(&verbosity, "level", "warn", "set log level")

}
