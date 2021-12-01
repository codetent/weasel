package wsl

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type RunOpts struct {
	User        string
	AttachStdin bool
}

func CreateWSLCommand(arg ...string) *exec.Cmd {
	return exec.Command("wsl", arg...)
}

func DecodeUTF16(bytes []byte) string {
	buffer := make([]uint16, len(bytes)/2)
	for i := 0; i < len(buffer); i++ {
		buffer[i] = binary.LittleEndian.Uint16(bytes[(i * 2):])
	}

	return syscall.UTF16ToString(buffer)
}

func Import(id string, workspace string, archive string) error {
	cmd := CreateWSLCommand("--import", id, workspace, archive, "--version", "2")
	out, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Error(string(exitError.Stderr))
		}
		return fmt.Errorf("Import: Output(): %v", err)
	}

	log.Debug(out)
	return nil
}

func Unregister(id string) error {
	cmd := CreateWSLCommand("--unregister", id)
	out, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Error(string(exitError.Stderr))
		}
		return fmt.Errorf("Unregister: Output(): %v", err)
	}

	log.Debug(string(out))
	return nil
}

func Run(id string, opts *RunOpts, arg ...string) error {
	optArgs := []string{"--distribution", id}

	if opts.User != "" {
		optArgs = append(optArgs, "--user", opts.User)
	}

	cmd := CreateWSLCommand(append(optArgs, arg...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if opts.AttachStdin {
		cmd.Stdin = os.Stdin
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Run: Run(): %v", err)
	}

	return nil
}

func ListRunning() ([]string, error) {
	cmd := CreateWSLCommand("--list", "--quiet", "--running")
	out, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Error(string(exitError.Stderr))
		}
		return nil, fmt.Errorf("ListRunning: cmd.Output(): %v", err)
	}

	dists := strings.Fields(DecodeUTF16(out))
	return dists, nil
}
