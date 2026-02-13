package cmd

import (
	"fmt"

	"github.com/jeffersonraimon/bgpefl/internal/gobgp"
	"github.com/jeffersonraimon/bgpefl/internal/system"

	"github.com/spf13/cobra"
)

var (
	clearV4 bool
	clearV6 bool
	soft    bool
	force   bool
)

var clearCmd = &cobra.Command{
	Use:   "clearrib",
	Short: "Limpa RIB do BGP Lab",
	RunE: func(cmd *cobra.Command, args []string) error {

		if !system.IsProcessRunning("gobgpd") {
			return fmt.Errorf("gobgpd não está rodando")
		}

		fmt.Println("🧹 Limpando RIB...")

		// Default: limpar ambos
		if !cmd.Flags().Changed("ipv4") && !cmd.Flags().Changed("ipv6") {
			clearV4 = true
			clearV6 = true
		}

		if force {
			if clearV4 {
				fmt.Println("Removendo todas rotas IPv4...")
				gobgp.ClearRIB("ipv4")
			}
			if clearV6 {
				fmt.Println("Removendo todas rotas IPv6...")
				gobgp.ClearRIB("ipv6")
			}
		} else {
			if clearV4 {
				fmt.Println("Removendo rotas IPv4...")
				gobgp.ClearRIBSoft("ipv4")
			}
			if clearV6 {
				fmt.Println("Removendo rotas IPv6...")
				gobgp.ClearRIBSoft("ipv6")
			}
		}

		fmt.Println("✅ RIB limpa.")
		return nil
	},
}

func init() {

	clearCmd.Flags().BoolVar(&clearV4, "ipv4", false, "Limpa apenas IPv4")
	clearCmd.Flags().BoolVar(&clearV6, "ipv6", false, "Limpa apenas IPv6")
	clearCmd.Flags().BoolVar(&soft, "soft", false, "Remove rota por rota (default)")
	clearCmd.Flags().BoolVar(&force, "force", false, "Usa del all direto")

	rootCmd.AddCommand(clearCmd)
}
