package cmd

import (
	"fmt"

	"github.com/jeffersonraimon/bgpefl/internal/gobgp"
	"github.com/jeffersonraimon/bgpefl/internal/netutil"
	"github.com/jeffersonraimon/bgpefl/internal/system"

	"github.com/spf13/cobra"
)

var (
	clearRIB bool
	rmInt    string
	rmIP     string
	rmCIDR   int
)

var stopCmd = &cobra.Command{
	var	force bool
	
	Use:   "stop",
	Short: "Para o BGP Lab",
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println("🛑 Parando BGP Lab...")

		// Remove neighbors
		fmt.Println("Removendo neighbors...")
		neighbors, _ := gobgp.ListNeighbors()

		for _, n := range neighbors {
			fmt.Printf("Removendo neighbor %s\n", n)
			gobgp.RemoveNeighbor(n)
		}

		// Clear RIB
		if clearRIB {
			fmt.Println("Limpando RIB IPv4...")
			gobgp.ClearRIB("ipv4")

			fmt.Println("Limpando RIB IPv6...")
			gobgp.ClearRIB("ipv6")
		}

		// Remove IP
		if rmInt != "" && rmIP != "" && rmCIDR > 0 {
			fmt.Printf("Removendo IP %s/%d da interface %s\n", rmIP, rmCIDR, rmInt)
			netutil.RemoveIP(rmInt, rmIP, rmCIDR)
		}

		// Stop daemon
		fmt.Println("Finalizando gobgpd...")

		err := system.StopProcess("gobgpd", force)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("\n✅ BGP Lab parado.")
		return nil
	},
}

func init() {

	stopCmd.Flags().BoolVar(&clearRIB, "clear-rib", false, "Remove todas as rotas")
	stopCmd.Flags().BoolVar(&force, "force", false, "Força kill")
	stopCmd.Flags().StringVar(&rmInt, "remove-int", "", "Interface para remover IP")
	stopCmd.Flags().StringVar(&rmIP, "remove-ip", "", "IP a remover")
	stopCmd.Flags().IntVar(&rmCIDR, "remove-cidr", 0, "CIDR do IP")

	rootCmd.AddCommand(stopCmd)
}
