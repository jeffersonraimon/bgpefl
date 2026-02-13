package gobgp

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"
)

func ListNeighbors() ([]string, error) {

	cmd := exec.Command("gobgp", "neighbor")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var neighbors []string

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "Neighbor") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) > 0 {
			neighbors = append(neighbors, fields[0])
		}
	}

	return neighbors, nil
}

func RemoveNeighbor(ip string) error {
	return exec.Command("gobgp", "neighbor", "del", ip).Run()
}
