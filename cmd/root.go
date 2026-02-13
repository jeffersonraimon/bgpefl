package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "bgpefl",
	Short: "BGP Easy for Labs",
}

func Execute() error {
	return rootCmd.Execute()
}
