package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "animux",
	Short: "A terminal virtual pet",
	Long:  `Animux is a virtual pet that lives in your terminal. You can adopt, feed, play, and take care of it!`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(adoptCmd)
	rootCmd.AddCommand(feedCmd)
	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(daemonCmd)
}
