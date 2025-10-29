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

		if targetNumber != 0 {
			targetNumber--

			utils.SetUpLogs(verbosity)

			utils.RotUsed = true
			err := filepath.Walk(lightsdir, utils.Flatsversal)
			if err != nil {
				log.Fatal(err)
			}

			Log.Debugf("Number of targets detected: %d", len(utils.Rotations))
			if targetNumber < len(utils.Rotations) {
				fmt.Println(utils.Rotations[targetNumber])
			} else {
				fmt.Println(-1)
			}
		} else {
			fmt.Println(-1)
		}
	},
}

func init() {
	rootCmd.AddCommand(rotationCmd)
	rotationCmd.Flags().IntVar(&targetNumber, "target", 0, "night target number, between 1 and 3")

}
