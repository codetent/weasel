/*
Copyright © 2021 Christoph Swoboda

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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

func Run(id string, opts *RunOpts, arg ...string) (int, error) {
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
	code := cmd.ProcessState.ExitCode()
	if err != nil {
		return code, fmt.Errorf("Run: Run(): %v", err)
	}

	return code, nil
}

func Terminate(id string) error {
	cmd := CreateWSLCommand("--terminate", id)
	out, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Error(string(exitError.Stderr))
		}
		return fmt.Errorf("Terminate: Output(): %v", err)
	}

	log.Debug(string(out))
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
