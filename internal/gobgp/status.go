package gobgp

import (
	"bytes"
	"os/exec"
)

func GetGlobal() (string, error) {
	cmd := exec.Command("gobgp", "global")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}

func GetNeighbors() (string, error) {
	cmd := exec.Command("gobgp", "neighbor")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return out.String(), err
}
