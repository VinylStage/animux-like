package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/VinylStage/animux-like/daemon"
	"github.com/VinylStage/animux-like/ui"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play with your pet",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := daemon.SendRequest("/play", nil)
		if err != nil {
			fmt.Println(ui.RenderError(err))
			return
		}
		fmt.Println(ui.RenderResponse(state, "Yay! So fun!"))
	},
}
