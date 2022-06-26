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
	"fmt"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func _CreateWSLCommand(arg ...string) *exec.Cmd {
	return exec.Command("wsl", arg...)
}

func Import(id string, workspace string, archive string) error {
	cmd := _CreateWSLCommand("--import", id, workspace, archive, "--version", "2")
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
