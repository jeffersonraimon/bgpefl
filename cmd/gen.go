package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"strconv"

	"github.com/spf13/cobra"
	api "github.com/osrg/gobgp/v3/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Variáveis para as flags
var (
	targetAS   int
	limit      int
	limitV4    int
	limitV6    int
	onlyV4     bool
	onlyV6     bool
	minV4      int
	minV6      int
	dryRun     bool
	irrServer  string
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Gera e injeta rotas baseadas no ASN informado",
	Run: func(cmd *cobra.Command, args []string) {
		if targetAS == 0 {
			log.Fatal("❌ O parâmetro --as é obrigatório.")
		}

		fmt.Printf("🔎 Buscando prefixos do AS%d em %s...\n", targetAS, irrServer)
		
		prefixes, err := fetchIRRPrefixes(targetAS, irrServer)
		if err != nil {
			log.Fatalf("❌ Erro ao buscar no IRR: %v", err)
		}

		// Conectar ao GoBGP
		conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("❌ Erro ao conectar ao GoBGP gRPC: %v", err)
		}
		defer conn.Close()
		client := api.NewGobgpServiceClient(conn)

		v4Count, v6Count, total := 0, 0, 0

		for _, prefix := range prefixes {
			// Lógica de Filtros
			isV6 := strings.Contains(prefix, ":")
			parts := strings.Split(prefix, "/")
			mask, _ := strconv.Atoi(parts[1])

			if isV6 && onlyV4 { continue }
			if !isV6 && onlyV6 { continue }

			// Filtro de tamanho de máscara
			if !isV6 && minV4 > 0 && mask < minV4 { continue }
			if isV6 && minV6 > 0 && mask < minV6 { continue }

			// Filtros de limites
			if !isV6 && limitV4 > 0 && v4Count >= limitV4 { continue }
			if isV6 && limitV6 > 0 && v6Count >= limitV6 { continue }
			if limit > 0 && total >= limit {
				fmt.Printf("⚠️ Limite total atingido (%d)\n", limit)
				break
			}

			if dryRun {
				fmt.Printf("[DRY] %s\n", prefix)
			} else {
				addRoute(client, prefix, isV6)
			}

			if isV6 { v6Count++ } else { v4Count++ }
			total++
		}

		fmt.Printf("\n========== RESUMO ==========\n")
		fmt.Printf("IPv4: %d\nIPv6: %d\nTotal: %d\n", v4Count, v6Count, total)
		fmt.Println("============================")
	},
}

// fetchIRRPrefixes conecta via TCP ao servidor WHOIS e extrai os prefixos
func fetchIRRPrefixes(asn int, server string) ([]string, error) {
	conn, err := net.Dial("tcp", server+":43")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Comando padrão RADB para buscar rotas por origem
	query := fmt.Sprintf("-i origin AS%d\r\n", asn)
	fmt.Fprintf(conn, query)

	var prefixes []string
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "route:") || strings.HasPrefix(line, "route6:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				prefixes = append(prefixes, parts[1])
			}
		}
	}
	return prefixes, nil
}

// addRoute injeta a rota no GoBGP via gRPC
func addRoute(client api.GobgpServiceClient, prefix string, isV6 bool) {
	family := &api.Family{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_UNICAST}
	if isV6 {
		family.Afi = api.Family_AFI_IP6
	}

	parts := strings.Split(prefix, "/")
	addr := parts[0]
	mask, _ := strconv.Atoi(parts[1])

	nlri, _ := api.NewNlri(family, addr, uint32(mask))
	path := &api.Path{
		Family: family,
		Nlri:   nlri,
		Pattributes: []*api.PathAttribute{}, // Adicione atributos como NextHop se necessário
	}

	_, err := client.AddPath(context.Background(), &api.AddPathRequest{Path: path})
	if err != nil {
		fmt.Printf("❌ Erro ao adicionar %s: %v\n", prefix, err)
	}
}

func init() {
	genCmd.Flags().IntVar(&targetAS, "as", 0, "AS para buscar prefixos (Obrigatório)")
	genCmd.Flags().IntVar(&limit, "limit", 0, "Limite total de rotas")
	genCmd.Flags().IntVar(&limitV4, "limit-v4", 0, "Limite apenas IPv4")
	genCmd.Flags().IntVar(&limitV6, "limit-v6", 0, "Limite apenas IPv6")
	genCmd.Flags().BoolVar(&onlyV4, "only-v4", false, "Apenas IPv4")
	genCmd.Flags().BoolVar(&onlyV6, "only-v6", false, "Apenas IPv6")
	genCmd.Flags().IntVar(&minV4, "min-v4", 0, "Prefixo mínimo IPv4 (ex: 24)")
	genCmd.Flags().IntVar(&minV6, "min-v6", 0, "Prefixo mínimo IPv6 (ex: 48)")
	genCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Não adiciona no gobgp")
	genCmd.Flags().StringVar(&irrServer, "irr", "whois.radb.net", "IRR server")
}