package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/VinylStage/animux-like/daemon"
	"github.com/VinylStage/animux-like/ui"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean your pet's room",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := daemon.SendRequest("/clean", nil)
		if err != nil {
			fmt.Println(ui.RenderError(err))
			return
		}
		fmt.Println(ui.RenderResponse(state, "All sparkling clean!"))
	},
}
