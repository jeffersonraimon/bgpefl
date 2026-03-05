package netutil

import (
        "fmt"
        "net"
        "os/exec"
		"strings"
		"bytes"
		"bufio"
)

type IPVersion int

const (
        IPv4 IPVersion = 4
        IPv6 IPVersion = 6
        All  IPVersion = 0
)


func AddRoute(prefix string) error {
        ip, _, err := net.ParseCIDR(prefix)
        if err != nil {
                return fmt.Errorf("invalid prefix %s: %w", prefix, err)
        }

        var cmd *exec.Cmd

        if ip.To4() == nil {
                // IPv6
                cmd = exec.Command(
                        "ip", "-6", "route", "replace",
                        "local", prefix,
                        "dev", "lo", "metric", "888",
                )
        } else {
                // IPv4
                cmd = exec.Command(
                        "ip", "route", "replace",
                        "local", prefix,
                        "dev", "lo", "metric", "888",
                )
        }

        out, err := cmd.CombinedOutput()
        if err != nil {
                return fmt.Errorf("ip route failed: %v (%s)", err, string(out))
        }

        return nil
}


// ClearProgramLocalRoutesByVersion remove rotas locais criadas pelo programa.
// version = IPv4, IPv6 ou All
func ClearProgramLocalRoutesByVersion(version IPVersion) error {
        if version == IPv4 || version == All {
                if err := clearMarkedLocalRoutes("-4"); err != nil {
                        return fmt.Errorf("failed to clear IPv4 routes: %w", err)
                }
        }

        if version == IPv6 || version == All {
                if err := clearMarkedLocalRoutes("-6"); err != nil {
                        return fmt.Errorf("failed to clear IPv6 routes: %w", err)
                }
        }

        return nil
}

func clearMarkedLocalRoutes(ipVersion string) error {
        args := []string{ipVersion, "route", "show"}

        // Apenas IPv6 usa table local
        if ipVersion == "-6" {
                args = append(args, "table", "local")
        }

        cmd := exec.Command("ip", args...)
        out, err := cmd.Output()
        if err != nil {
                return fmt.Errorf("failed to list routes: %v", err)
        }

        scanner := bufio.NewScanner(bytes.NewReader(out))
        for scanner.Scan() {
                line := scanner.Text()

                if strings.HasPrefix(line, "local ") &&
                        strings.Contains(line, "dev lo") &&
                        strings.Contains(line, "metric 888") {

                        parts := strings.Fields(line)
                        if len(parts) < 2 {
                                continue
                        }

                        prefix := parts[1]

                        delArgs := []string{
                                ipVersion,
                                "route", "del",
                                "local", prefix,
                                "dev", "lo",
                                "metric", "888",
                        }

                        delCmd := exec.Command("ip", delArgs...)
                        if out, err := delCmd.CombinedOutput(); err != nil {
                                fmt.Printf("failed to delete route %s: %v (%s)\n", prefix, err, string(out))
                        }
                }
        }

        return nil
}