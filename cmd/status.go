package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	api "github.com/osrg/gobgp/v3/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Mostra o status atual do GoBGP e rotas",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Println("gobgpd: STOPPED")
			return
		}
		defer conn.Close()
		client := api.NewGobgpServiceClient(conn)
		ctx := context.Background()

		fmt.Println("========== BGP LAB STATUS ==========")
		fmt.Println("gobgpd: RUNNING\n")

		// 1. Configuração Global
		g, _ := client.GetBgp(ctx, &api.GetBgpRequest{})
		if g != nil {
			fmt.Printf("Global config:\nAS: %d\nRouter-ID: %s\n\n", g.Global.As, g.Global.RouterId)
		}

		// 2. Neighbors
		fmt.Println("Neighbors:")
		client.ListPeer(ctx, &api.ListPeerRequest{}, func(p *api.Peer) {
			fmt.Printf("Peer: %s | AS: %d | State: %s\n", 
				p.Conf.NeighborAddress, p.Conf.PeerAs, p.State.SessionState)
		})

		// 3. Contagem de Rotas (v4 e v6)
		v4Count := countRoutes(client, api.Family_AFI_IP)
		v6Count := countRoutes(client, api.Family_AFI_IP6)

		fmt.Printf("\nRotas IPv4: %d\n", v4Count)
		fmt.Printf("Rotas IPv6: %d\n", v6Count)
		fmt.Println("====================================")
	},
}

func countRoutes(client api.GobgpServiceClient, afi api.Family_Afi) int {
	count := 0
	err := client.ListPath(context.Background(), &api.ListPathRequest{
		TableType: api.TableType_GLOBAL,
		Family:    &api.Family{Afi: afi, Safi: api.Family_SAFI_UNICAST},
	}, func(p *api.Destination) {
		count++
	})
	if err != nil { return 0 }
	return count
}