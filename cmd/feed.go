package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/VinylStage/animux-like/daemon"
	"github.com/VinylStage/animux-like/ui"
)

var feedCmd = &cobra.Command{
	Use:   "feed",
	Short: "Feed your pet",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := daemon.SendRequest("/feed", nil)
		if err != nil {
			fmt.Println(ui.RenderError(err))
			return
		}
		fmt.Println(ui.RenderResponse(state, "Yum! That was delicious."))
	},
}
