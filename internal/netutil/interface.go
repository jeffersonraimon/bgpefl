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
		cmd = exec.Command("ip", "addr", "add", fmt.Sprintf("%s/%d", ip, cidr), "dev", iface)
	} else {
		cmd = exec.Command("ip", "-6", "addr", "add", fmt.Sprintf("%s/%d", ip, cidr), "dev", iface)
	}

	return cmd.Run()
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
