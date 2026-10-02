package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/VinylStage/animux-like/daemon"
	"github.com/VinylStage/animux-like/ui"
)

var adoptCmd = &cobra.Command{
	Use:   "adopt [species] [name]",
	Short: "Adopt a new pet (e.g. cat, dog, penguin)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		query := map[string]string{
			"species": args[0],
			"name":    args[1],
		}
		state, err := daemon.SendRequest("/adopt", query)
		if err != nil {
			fmt.Println(ui.RenderError(err))
			return
		}
		fmt.Println(ui.RenderResponse(state, fmt.Sprintf("You just adopted %s!", state.Name)))
	},
}
