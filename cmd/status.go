package cmd

import (
	"fmt"

	"bgpefl/internal/gobgp"
	"bgpefl/internal/system"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Mostra status do BGP Lab",
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println("========== BGP LAB STATUS ==========")

		running := system.IsProcessRunning("gobgpd")

		if !running {
			fmt.Println("gobgpd: STOPPED")
			fmt.Println("====================================")
			return nil
		}

		fmt.Println("gobgpd: RUNNING\n")

		global, _ := gobgp.GetGlobal()
		fmt.Println("Global config:")
		fmt.Println(global)

		neighbors, _ := gobgp.GetNeighbors()
		fmt.Println("\nNeighbors:")
		fmt.Println(neighbors)

		v4, _ := gobgp.GetRIB("ipv4")
		v6, _ := gobgp.GetRIB("ipv6")

		fmt.Printf("\nRotas IPv4: %d\n", len(v4))
		fmt.Printf("Rotas IPv6: %d\n", len(v6))

		fmt.Println("\nPreview IPv4:")
		printPreview(v4)

		fmt.Println("\nPreview IPv6:")
		printPreview(v6)

		fmt.Println("====================================")

		return nil
	},
}

func printPreview(routes []string) {
	limit := 7
	if len(routes) < limit {
		limit = len(routes)
	}

	for i := 0; i < limit; i++ {
		fmt.Println(routes[i])
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
