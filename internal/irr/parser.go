package irr

import (
	"bufio"
	"strings"
)

func ParsePrefixes(raw string) []string {

	seen := make(map[string]bool)
	var result []string

	scanner := bufio.NewScanner(strings.NewReader(raw))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "route:") ||
			strings.HasPrefix(line, "route6:") {

			parts := strings.Fields(line)
			if len(parts) >= 2 {
				prefix := parts[1]
				if !seen[prefix] {
					seen[prefix] = true
					result = append(result, prefix)
				}
			}
		}
	}

	return result
}
