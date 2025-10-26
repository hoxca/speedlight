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

var (
	targetNumber int
)

// flatsCmd get the target list of filters
var filtersCmd = &cobra.Command{
	Use:   "filters",
	Short: "Get the target list of filters",
	Long:  `Get the list of filters used for this target during the last acquisition night`,
	Run: func(cmd *cobra.Command, args []string) {

		lightsdir, _ = strings.CutSuffix(lightsdir, "/")
		lightsdir, _ = strings.CutSuffix(lightsdir, "\\")
		utils.Wdest = utils.WriteDestination{writeConsole, writeReport}

		if targetNumber == 0 {
			rotation = false
		} else {
			targetNumber--
		}
		utils.SetUpLogs(verbosity)
		utils.Wdest = utils.WriteDestination{writeConsole, writeReport}

		utils.RotUsed, _ = cmd.Flags().GetBool("rotation")
		err := filepath.Walk(lightsdir, utils.Flatsversal)
		if err != nil {
			log.Fatal(err)
		}

		Log.Debugf("Number of targets detected: %d", len(utils.Rotations))
		if targetNumber < len(utils.Rotations) {
			if !rotation {
				fmt.Println(utils.FlatList[666])

			} else {
				if utils.Wdest.WriteToConsole {
					fltrs := utils.FlatList[utils.Rotations[targetNumber]]
					fmt.Println(fltrs)
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(filtersCmd)
	filtersCmd.Flags().IntVar(&targetNumber, "target", 0, "night target number, between 1 and 3")
}
