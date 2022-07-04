package utils

import (
	"os/exec"
	"strings"
)

func Which(name string) (string, error) {
	cmd := exec.Command("where", name)

	cmdBytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	path := string(cmdBytes)
	return strings.TrimSpace(path), nil
}
