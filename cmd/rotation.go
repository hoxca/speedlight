package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"speedlight/utils"

	Log "github.com/apatters/go-conlog"
	"github.com/spf13/cobra"
)

// rotationCmd get the target rotation
var rotationCmd = &cobra.Command{
	Use:   "rotation",
	Short: "Get the target rotation",
	Long:  `Get the rotation used by the target during last acquisition night`,
	Run: func(cmd *cobra.Command, args []string) {
		lightsdir, _ = strings.CutSuffix(lightsdir, "/")
		lightsdir, _ = strings.CutSuffix(lightsdir, "\\")
		utils.Wdest = utils.WriteDestination{writeConsole, writeReport}

		targetNumber--
		utils.SetUpLogs(verbosity)
		utils.Wdest = utils.WriteDestination{writeConsole, writeReport}

		utils.RotUsed, _ = cmd.Flags().GetBool("rotation")
		err := filepath.Walk(lightsdir, utils.Flatsversal)
		if err != nil {
			log.Fatal(err)
		}

		Log.Debugf("Number of targets detected: %d", len(utils.Rotations))
		if targetNumber < len(utils.Rotations) {
			if utils.Wdest.WriteToConsole {
				fmt.Println(utils.Rotations[targetNumber])
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(rotationCmd)
	rotationCmd.Flags().IntVar(&targetNumber, "target", 1, "night target number, between 1 and 3")

}
