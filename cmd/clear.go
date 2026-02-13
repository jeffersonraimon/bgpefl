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

var (
	clearIPv4 bool
	clearIPv6 bool
	force     bool
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Limpa prefixos da RIB (IPv4/IPv6)",
	Run: func(cmd *cobra.Command, args []string) {
		// Se o usuário não especificou nada, limpa ambos por padrão
		if !clearIPv4 && !clearIPv6 {
			clearIPv4 = true
			clearIPv6 = true
		}

		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("❌ Erro ao conectar ao GoBGP: %v", err)
		}
		defer conn.Close()
		client := api.NewGobgpServiceClient(conn)
		ctx := context.Background()

		fmt.Println("🧹 Limpando RIB...")

		if clearIPv4 {
			clearFamily(ctx, client, api.Family_AFI_IP, force)
		}
		if clearIPv6 {
			clearFamily(ctx, client, api.Family_AFI_IP6, force)
		}

		fmt.Println("✅ RIB limpa.")
	},
}

func clearFamily(ctx context.Context, client api.GobgpServiceClient, afi api.Family_Afi, forceMode bool) {
	family := &api.Family{Afi: afi, Safi: api.Family_SAFI_UNICAST}
	label := "IPv4"
	if afi == api.Family_AFI_IP6 {
		label = "IPv6"
	}

	if forceMode {
		fmt.Printf("Removendo todas as rotas %s (modo force)...", label)
		// No gRPC, DeletePath sem especificar um prefixo limpa a tabela daquela família
		_, err := client.DeletePath(ctx, &api.DeletePathRequest{
			TableType: api.TableType_GLOBAL,
			Family:    family,
		})
		if err != nil {
			fmt.Printf("❌ Erro: %v\n", err)
		}
	} else {
		fmt.Printf("Removendo rotas %s (modo soft)...", label)
		// Lista e remove uma por uma
		err := client.ListPath(ctx, &api.ListPathRequest{
			TableType: api.TableType_GLOBAL,
			Family:    family,
		}, func(destination *api.Destination) {
			for _, path := range destination.Paths {
				client.DeletePath(ctx, &api.DeletePathRequest{
					TableType: api.TableType_GLOBAL,
					Family:    family,
					Path:      path,
				})
			}
		})
		if err != nil {
			fmt.Printf("❌ Erro ao listar: %v\n", err)
		}
	}
	fmt.Println(" pronto.")
}

func init() {
	clearCmd.Flags().BoolVar(&clearIPv4, "ipv4", false, "Limpa apenas IPv4")
	clearCmd.Flags().BoolVar(&clearIPv6, "ipv6", false, "Limpa apenas IPv6")
	clearCmd.Flags().BoolVar(&force, "force", false, "Usa remoção em massa (mais rápido)")
	// O modo 'soft' é o padrão, então não precisamos de uma flag específica
	// a menos que queira manter por compatibilidade de comando.
}