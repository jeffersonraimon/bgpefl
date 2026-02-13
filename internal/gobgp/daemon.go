package gobgp

import (
	"bytes"
	"os/exec"
	"regexp"
	"time"
	"strconv"
	"github.com/jeffersonraimon/bgpefl/internal/system"
)

func StartDaemon() error {
	cmd := exec.Command("gobgpd")
	if err := cmd.Start(); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)
	return nil
}

func IsRunning() bool {
    return system.IsProcessRunning("gobgpd")
}

func GetGlobalAS() (uint32, error) {
	cmd := exec.Command("gobgp", "global")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0, err
	}

	// Procura algo como: "AS: 52872"
	re := regexp.MustCompile(`AS:\s+(\d+)`)
	matches := re.FindStringSubmatch(out.String())
	if len(matches) < 2 {
		return 0, nil
	}

	as, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	return uint32(as), nil
}