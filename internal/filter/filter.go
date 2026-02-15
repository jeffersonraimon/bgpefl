package filter

import (
	"net"
	"strconv"
	"strings"
)

type Config struct {
	Limit   int
	LimitV4 int
	LimitV6 int
	OnlyV4  bool
	OnlyV6  bool
	MinV4   int
	MinV6   int
}

func IsIPv4(prefix string) bool {
	ip := strings.Split(prefix, "/")[0]
	return net.ParseIP(ip).To4() != nil
}

func Apply(prefixes []string, cfg Config) []string {

	var v4 []string
	var v6 []string

	// Primeiro: filtra e separa
	for _, p := range prefixes {

		familyV4 := IsIPv4(p)

		maskStr := strings.Split(p, "/")[1]
		mask, _ := strconv.Atoi(maskStr)

		// only flags
		if cfg.OnlyV4 && !familyV4 {
			continue
		}
		if cfg.OnlyV6 && familyV4 {
			continue
		}

		// min mask
		if familyV4 && cfg.MinV4 > 0 && mask < cfg.MinV4 {
			continue
		}
		if !familyV4 && cfg.MinV6 > 0 && mask < cfg.MinV6 {
			continue
		}

		if familyV4 {
			v4 = append(v4, p)
		} else {
			v6 = append(v6, p)
		}
	}

	// Aplicar limites individuais
	if cfg.LimitV4 > 0 && len(v4) > cfg.LimitV4 {
		v4 = v4[:cfg.LimitV4]
	}

	if cfg.LimitV6 > 0 && len(v6) > cfg.LimitV6 {
		v6 = v6[:cfg.LimitV6]
	}

	// Aplicar limite global proporcional
	if cfg.Limit > 0 {
	
	    half := cfg.Limit / 2
	
	    // Tentativa inicial 50/50
	    v4Limit := half
	    v6Limit := cfg.Limit - half
	
	    // Ajusta se não houver suficiente
	    if v4Limit > len(v4) {
	        v4Limit = len(v4)
	    }
	    if v6Limit > len(v6) {
	        v6Limit = len(v6)
	    }
	
	    totalSelected := v4Limit + v6Limit
	
	    // Se ainda não atingiu o limite, redistribui sobra
	    if totalSelected < cfg.Limit {
		
	        remaining := cfg.Limit - totalSelected
		
	        // Primeiro tenta completar com IPv4
	        if v4Limit < len(v4) {
	            extra := len(v4) - v4Limit
	            if extra > remaining {
	                extra = remaining
	            }
	            v4Limit += extra
	            remaining -= extra
	        }
		
	        // Depois tenta completar com IPv6
	        if remaining > 0 && v6Limit < len(v6) {
	            extra := len(v6) - v6Limit
	            if extra > remaining {
	                extra = remaining
	            }
	            v6Limit += extra
	        }
	    }
	
	    v4 = v4[:v4Limit]
	    v6 = v6[:v6Limit]
	}

	return append(v4, v6...)
}

