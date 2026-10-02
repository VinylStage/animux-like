package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/VinylStage/animux-like/daemon"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start the background daemon",
	Run: func(cmd *cobra.Command, args []string) {
		if err := daemon.Start(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}
