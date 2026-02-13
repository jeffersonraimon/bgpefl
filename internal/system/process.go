package system

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func IsRunning() bool {
    return system.IsProcessRunning("gobgpd")
}

func IsProcessRunning(name string) bool {
	err := exec.Command("pgrep", name).Run()
	return err == nil
}


func StopProcess(name string, force bool) error {

	out, err := exec.Command("pgrep", name).Output()
	if err != nil {
		return fmt.Errorf("%s não estava rodando", name)
	}

	pids := strings.Fields(string(out))

	for _, pidStr := range pids {
		pid, _ := strconv.Atoi(pidStr)
		process, _ := os.FindProcess(pid)
		process.Signal(os.Interrupt)
	}

	time.Sleep(1 * time.Second)

	// Check again
	if exec.Command("pgrep", name).Run() == nil {

		if force {
			fmt.Println("Forçando kill...")

			for _, pidStr := range pids {
				pid, _ := strconv.Atoi(pidStr)
				process, _ := os.FindProcess(pid)
				process.Kill()
			}
		} else {
			return fmt.Errorf("%s ainda está rodando. Use --force se necessário", name)
		}
	}

	return nil
}