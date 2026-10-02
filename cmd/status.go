package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/VinylStage/animux-like/daemon"
	"github.com/VinylStage/animux-like/ui"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check your pet's status",
	Run: func(cmd *cobra.Command, args []string) {
		state, err := daemon.SendRequest("/status", nil)
		if err != nil {
			fmt.Println(ui.RenderError(err))
			return
		}
		msg := "Just hanging out."
		if state.IsSick {
			msg = "I don't feel so good..."
		} else if state.Hunger < 30 {
			msg = "I'm so hungry..."
		} else if state.Happiness < 30 {
			msg = "I'm bored..."
		} else if state.Cleanliness < 30 {
			msg = "It smells in here..."
		}
		fmt.Println(ui.RenderResponse(state, msg))
	},
}
