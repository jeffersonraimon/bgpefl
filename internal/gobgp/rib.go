package gobgp


import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
	"encoding/json"
)

type RibEntry struct {
	Prefix string `json:"prefix"`
}

func AddPrefix(prefix string) error {

	if isIPv4(prefix) {
		return exec.Command("gobgp", "global", "rib", "add", "-a", "ipv4", prefix).Run()
	}

	return exec.Command("gobgp", "global", "rib", "add", "-a", "ipv6", prefix).Run()
}

func isIPv4(prefix string) bool {
	return !contains(prefix, ":")
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (len(sub) == 0 || (len(s) > 0 && (string(s[0]) == sub || contains(s[1:], sub))))
}

func GetRIB(family string) ([]string, error) {

	cmd := exec.Command("gobgp", "global", "rib", "-a", family)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var routes []string

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "Network") {
			continue
		}

		if strings.TrimSpace(line) != "" {
			routes = append(routes, line)
		}
	}

	return routes, nil
}

func ClearRIB(family string) error {
	return exec.Command("gobgp", "global", "rib", "del", "all").Run()
}

func GetRIBJSON(family string) ([]RibEntry, error) {

	cmd := exec.Command("gobgp", "global", "rib", "-a", family, "-j")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var entries []RibEntry
	err = json.Unmarshal(out, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func ClearRIBSoft(family string) error {

	entries, err := GetRIBJSON(family)
	if err != nil {
		return err
	}

	seen := make(map[string]bool)

	for _, e := range entries {

		if e.Prefix == "" {
			continue
		}

		if seen[e.Prefix] {
			continue
		}

		seen[e.Prefix] = true

		exec.Command("gobgp", "global", "rib", "del",
			"-a", family,
			e.Prefix).Run()
	}

	return nil
}