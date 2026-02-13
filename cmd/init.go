package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
	api "github.com/osrg/gobgp/v3/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var routerID string

var initCmd = &cobra.Command{
	Use:   "init [interface] [ip] [cidr] [neighbor_ip] [remote_as] [local_as]",
	Short: "Configura a interface de rede e a sessão BGP inicial",
	Args:  cobra.ExactArgs(6),
	Run: func(cmd *cobra.Command, args []string) {
		ifaceName := args[0]
		ipAddr := args[1]
		cidr := args[2]
		neighborIP := args[3]
		remoteAS := args[4]
		localAS := args[5]

		// 1. Configurar IP na Interface (Equivalente ao 'ip addr add')
		err := configureInterface(ifaceName, ipAddr, cidr)
		if err != nil {
			log.Fatalf("❌ Erro ao configurar interface: %v", err)
		}
		fmt.Printf("✅ IP %s/%s configurado na interface %s\n", ipAddr, cidr, ifaceName)

		// 2. Conectar ao GoBGP via gRPC (Assume que o gobgpd já está rodando ou foi iniciado)
		// Dica: Você pode usar o pacote 'os/exec' para rodar o 'gobgpd &' se desejar
		ctx := context.Background()
		conn, err := grpc.DialContext(ctx, "localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("❌ Erro ao conectar ao GoBGP: %v", err)
		}
		defer conn.Close()
		client := api.NewGobgpServiceClient(conn)

		// 3. Configuração Global (AS e Router-ID)
		_, err = client.StartBgp(ctx, &api.StartBgpRequest{
			Global: &api.Global{
				As:       parseUint32(localAS),
				RouterId: routerID,
			},
		})
		if err != nil {
			fmt.Printf("⚠️ Aviso: GoBGP já iniciado ou erro na config: %v\n", err)
		} else {
			fmt.Printf("✅ GoBGP Global configurado (AS: %s, Router-ID: %s)\n", localAS, routerID)
		}

		// 4. Adicionar Neighbor
		_, err = client.AddPeer(ctx, &api.AddPeerRequest{
			Peer: &api.Peer{
				Conf: &api.PeerConf{
					NeighborAddress: neighborIP,
					PeerAs:          parseUint32(remoteAS),
				},
			},
		})
		if err != nil {
			log.Fatalf("❌ Erro ao adicionar neighbor: %v", err)
		}

		fmt.Printf("🚀 Sessão BGP com neighbor %s (AS %s) configurada com sucesso!\n", neighborIP, remoteAS)
	},
}

// Função auxiliar para configurar IP usando Netlink (mais performático que chamar shell)
func configureInterface(name, ipStr, maskStr string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}

	addr, err := netlink.ParseAddr(fmt.Sprintf("%s/%s", ipStr, maskStr))
	if err != nil {
		return err
	}

	return netlink.AddrAdd(link, addr)
}

// Helper para converter string em uint32
func parseUint32(s string) uint32 {
	var res uint32
	fmt.Sscanf(s, "%d", &res)
	return res
}

func init() {
	// Define a flag opcional para router-id
	initCmd.Flags().StringVarP(&routerID, "router-id", "r", "10.99.99.99", "Router ID do BGP")
}