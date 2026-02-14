package netutil

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

func AddIP(iface, ip string, cidr int) error {

	parsed := net.ParseIP(ip)
	if parsed == nil {
		return fmt.Errorf("IP inválido")
	}

	var cmd *exec.Cmd

	if parsed.To4() != nil {
		cmd = exec.Command("ip", "addr", "replace",
			fmt.Sprintf("%s/%d", ip, cidr),
			"dev", iface)
	} else {
		cmd = exec.Command("ip", "-6", "addr", "replace",
			fmt.Sprintf("%s/%d", ip, cidr),
			"dev", iface)
	}

	// Adiciona / substitui o IP
	if err := cmd.Run(); err != nil {
		return err
	}

	// Sobe a interface
	upCmd := exec.Command("ip", "link", "set", iface, "up")
	if err := upCmd.Run(); err != nil {
		return err
	}

	return nil
}

func RemoveIP(iface, ip string, cidr int) error {

	if isIPv6(ip) {
		return exec.Command("ip", "-6", "addr", "del",
			fmt.Sprintf("%s/%d", ip, cidr),
			"dev", iface).Run()
	}

	return exec.Command("ip", "addr", "del",
		fmt.Sprintf("%s/%d", ip, cidr),
		"dev", iface).Run()
}

func isIPv6(ip string) bool {
	return strings.Contains(ip, ":")
}
