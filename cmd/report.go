package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"speedlight/utils"
	"strings"

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

		lightsdir, _ = strings.CutSuffix(lightsdir, "/")
		lightsdir, _ = strings.CutSuffix(lightsdir, "\\")
		fmt.Println("Scanning lights directory: " + lightsdir + "\n")

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

func init() {

	rootCmd.AddCommand(reportCmd)

}
