package cmd

import (
	"fmt"

	"bgpefl/internal/filter"
	"bgpefl/internal/gobgp"
	"bgpefl/internal/irr"

	"github.com/spf13/cobra"
)

var (
	targetAS uint32
	limit    int
	limitV4  int
	limitV6  int
	onlyV4   bool
	onlyV6   bool
	minV4    int
	minV6    int
	dryRun   bool
	irrHost  string
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Gera rotas BGP baseadas em IRR",
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Printf("🔎 Buscando prefixos do AS%d em %s...\n\n", targetAS, irrHost)

		raw, err := irr.FetchPrefixes(irrHost, targetAS)
		if err != nil {
			return err
		}

		prefixes := irr.ParsePrefixes(raw)

		if len(prefixes) == 0 {
			return fmt.Errorf("nenhum prefixo encontrado")
		}

		cfg := filter.Config{
			Limit:   limit,
			LimitV4: limitV4,
			LimitV6: limitV6,
			OnlyV4:  onlyV4,
			OnlyV6:  onlyV6,
			MinV4:   minV4,
			MinV6:   minV6,
		}

		final := filter.Apply(prefixes, cfg)

		v4, v6 := 0, 0

		for _, p := range final {

			if dryRun {
				fmt.Println("[DRY]", p)
			} else {
				if err := gobgp.AddPrefix(p); err != nil {
					return err
				}
			}

			if filter.IsIPv4(p) {
				v4++
			} else {
				v6++
			}
		}

		fmt.Println("\n========== RESUMO ==========")
		fmt.Printf("IPv4: %d\n", v4)
		fmt.Printf("IPv6: %d\n", v6)
		fmt.Printf("Total: %d\n", len(final))
		fmt.Println("============================")

		return nil
	},
}

func init() {

	genCmd.Flags().Uint32Var(&targetAS, "as", 0, "AS para buscar prefixos")
	genCmd.Flags().IntVar(&limit, "limit", 0, "Limite total")
	genCmd.Flags().IntVar(&limitV4, "limit-v4", 0, "Limite IPv4")
	genCmd.Flags().IntVar(&limitV6, "limit-v6", 0, "Limite IPv6")
	genCmd.Flags().BoolVar(&onlyV4, "only-v4", false, "Apenas IPv4")
	genCmd.Flags().BoolVar(&onlyV6, "only-v6", false, "Apenas IPv6")
	genCmd.Flags().IntVar(&minV4, "min-v4", 0, "Prefixo mínimo IPv4")
	genCmd.Flags().IntVar(&minV6, "min-v6", 0, "Prefixo mínimo IPv6")
	genCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Não adiciona no gobgp")
	genCmd.Flags().StringVar(&irrHost, "irr", "whois.radb.net", "Servidor IRR")

	genCmd.MarkFlagRequired("as")

	rootCmd.AddCommand(genCmd)
}
