package gobgp

import (
	"os/exec"
	"time"
)

func StartDaemon() error {
	cmd := exec.Command("gobgpd")
	if err := cmd.Start(); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)
	return nil
}
