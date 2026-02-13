package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
	api "github.com/osrg/gobgp/v3/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	stopClearRib bool
	stopForce    bool
	removeIP     []string // [interface, ip, cidr]
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Para a sessão BGP, limpa as configurações e finaliza o daemon",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🛑 Parando BGP Lab...")

		// 1. Conectar via gRPC para limpeza elegante
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithTimeout(2*time.Second))
		if err == nil {
			defer conn.Close()
			client := api.NewGobgpServiceClient(conn)
			ctx := context.Background()

			// Remover Neighbors
			fmt.Println("Removendo neighbors...")
			client.ListPeer(ctx, &api.ListPeerRequest{}, func(p *api.Peer) {
				fmt.Printf("Removendo neighbor %s\n", p.Conf.NeighborAddress)
				client.DeletePeer(ctx, &api.DeletePeerRequest{Address: p.Conf.NeighborAddress})
			})

			// Limpar RIB se solicitado
			if stopClearRib {
				fmt.Println("Limpando RIB IPv4 e IPv6...")
				families := []*api.Family{
					{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_UNICAST},
					{Afi: api.Family_AFI_IP6, Safi: api.Family_SAFI_UNICAST},
				}
				for _, f := range families {
					client.DeletePath(ctx, &api.DeletePathRequest{TableType: api.TableType_GLOBAL, Family: f})
				}
			}
		}

		// 2. Remover IP da Interface (se flags fornecidas)
		if len(removeIP) == 3 {
			iface, ip, cidr := removeIP[0], removeIP[1], removeIP[2]
			fmt.Printf("Removendo IP %s/%s da interface %s\n", ip, cidr, iface)
			err := removeInterfaceIP(iface, ip, cidr)
			if err != nil {
				fmt.Printf("⚠️ Erro ao remover IP: %v\n", err)
			}
		}

		// 3. Finalizar o processo gobgpd
		fmt.Println("Finalizando gobgpd...")
		terminateProcess("gobgpd", stopForce)

		fmt.Println("\n✅ BGP Lab parado.")
	},
}

func removeInterfaceIP(name, ipStr, maskStr string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}
	addr, err := netlink.ParseAddr(fmt.Sprintf("%s/%s", ipStr, maskStr))
	if err != nil {
		return err
	}
	return netlink.AddrDel(link, addr)
}

func terminateProcess(name string, force bool) {
	// Pega o PID usando pgrep
	out, _ := exec.Command("pgrep", name).Output()
	pids := strings.Fields(string(out))

	if len(pids) == 0 {
		fmt.Printf("%s não estava rodando.\n", name)
		return
	}

	for _, pidStr := range pids {
		cmd := exec.Command("kill", pidStr)
		if force {
			cmd = exec.Command("kill", "-9", pidStr)
			fmt.Println("Forçando kill (-9)...")
		}
		cmd.Run()
	}
}

func init() {
	stopCmd.Flags().BoolVar(&stopClearRib, "clear-rib", false, "Remove todas as rotas antes de parar")
	stopCmd.Flags().BoolVar(&stopForce, "force", false, "Força a parada do processo (kill -9)")
	stopCmd.Flags().StringSliceVar(&removeIP, "remove-ip", []string{}, "Remove IP (Ex: --remove-ip eth0,192.168.1.1,24)")
}