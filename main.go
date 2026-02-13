package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/jeffersonraimon/bgpefl/cmd"
)

var rootCmd = &cobra.Command{
	Use:   "bgpefl",
	Short: "BGPEFL - BGP Easy for Labs",
	Long:  `Ferramenta para facilitar a configuração de Upstream BGP em ambientes de Lab (EVE-NG/PNET)`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)   
	rootCmd.AddCommand(genCmd)   
	rootCmd.AddCommand(statusCmd) 
	rootCmd.AddCommand(clearCmd) 
	rootCmd.AddCommand(stopCmd)  
}