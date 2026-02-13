package gobgp

import (
	"fmt"
	"os/exec"
)

func ConfigureGlobal(as uint32, routerID string) error {
	cmd := exec.Command("gobgp", "global", "as", fmt.Sprint(as), "router-id", routerID)
	return cmd.Run()
}

func AddNeighbor(ip string, as uint32) error {
	cmd := exec.Command("gobgp", "neighbor", "add", ip, "as", fmt.Sprint(as))
	return cmd.Run()
}
