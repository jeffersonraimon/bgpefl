package cmd

import (
	"fmt"
	"net"

	"github.com/jeffersonraimon/bgpefl/internal/gobgp"
	"github.com/jeffersonraimon/bgpefl/internal/netutil"


	"github.com/spf13/cobra"
)

var (
	iface     string
	ipAddr    string
	cidr      int
	neighbor  string
	remoteAS  uint32
	localAS   uint32
	routerID  string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Inicializa sessão BGP para lab",
	RunE: func(cmd *cobra.Command, args []string) error {

		// 🔎 Validação básica
		if localAS == 0 || remoteAS == 0 {
			return fmt.Errorf("AS não pode ser 0")
		}

		if net.ParseIP(neighbor) == nil {
			return fmt.Errorf("neighbor inválido")
		}

		// 🔥 Router-ID opcional com fallback
		if routerID == "" {
			ipParsed := net.ParseIP(ipAddr)
			if ipParsed == nil {
				return fmt.Errorf("IP inválido")
			}
		
			// Se for IPv6
			if ipParsed.To4() == nil {
				routerID = "10.99.99.99"
			} else {
				// Se for IPv4, usa o próprio IP como router-id
				routerID = ipAddr
			}
		}

        running := gobgp.IsRunning()

		currentAS, _ := gobgp.GetGlobalAS()

		if running && currentAS != localAS {
		    return fmt.Errorf("gobgpd já está rodando com AS %d", currentAS)
		}

		if !running {
		    if err := gobgp.StartDaemon(); err != nil {
		        return fmt.Errorf("erro ao iniciar gobgpd: %v", err)
		    }
		
		    if err := gobgp.ConfigureGlobal(localAS, routerID); err != nil {
		        return fmt.Errorf("erro ao configurar global: %v", err)
		    }
		} else {
		    fmt.Println("ℹ️ gobgpd já está rodando, pulando configuração global")
		}

		if err := netutil.AddIP(iface, ipAddr, cidr); err != nil {
			return fmt.Errorf("erro ao configurar IP: %v", err)
		}

		if err := gobgp.AddNeighbor(neighbor, remoteAS); err != nil {
			return fmt.Errorf("erro ao adicionar neighbor: %v", err)
		}

		fmt.Println("✅ Sessão BGP configurada.")
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&iface, "int", "", "Interface")
	initCmd.Flags().StringVar(&ipAddr, "ip", "", "IP address")
	initCmd.Flags().IntVar(&cidr, "cidr", 0, "CIDR mask")
	initCmd.Flags().StringVar(&neighbor, "neighbor", "", "Neighbor IP")
	initCmd.Flags().Uint32Var(&remoteAS, "remote-as", 0, "Remote AS")
	initCmd.Flags().Uint32Var(&localAS, "local-as", 0, "Local AS")
	initCmd.Flags().StringVar(&routerID, "router-id", "", "Router ID (optional)")

	initCmd.MarkFlagRequired("int")
	initCmd.MarkFlagRequired("ip")
	initCmd.MarkFlagRequired("cidr")
	initCmd.MarkFlagRequired("neighbor")
	initCmd.MarkFlagRequired("remote-as")
	initCmd.MarkFlagRequired("local-as")

	rootCmd.AddCommand(initCmd)
}
