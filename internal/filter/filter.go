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

	var result []string
	v4Count := 0
	v6Count := 0

	for _, p := range prefixes {

		familyV4 := IsIPv4(p)
		maskStr := strings.Split(p, "/")[1]
		mask, _ := strconv.Atoi(maskStr)

		// only-v4 / only-v6
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

		// limits
		if familyV4 && cfg.LimitV4 > 0 && v4Count >= cfg.LimitV4 {
			continue
		}
		if !familyV4 && cfg.LimitV6 > 0 && v6Count >= cfg.LimitV6 {
			continue
		}

		if cfg.Limit > 0 && len(result) >= cfg.Limit {
			break
		}

		result = append(result, p)

		if familyV4 {
			v4Count++
		} else {
			v6Count++
		}
	}

	return result
}
