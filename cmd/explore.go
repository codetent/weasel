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
package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/codetent/weasel/pkg/weasel/wsl"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yuk7/wsllib-go"
)

type ExploreCmd struct {
	EnvName  string
	ShowOnly bool
}

func NewExploreCmd() *cobra.Command {
	cmd := &ExploreCmd{}

	exploreCmd := &cobra.Command{
		Use:   "explore",
		Short: "Open environment folder in explorer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			cmd.EnvName = args[0]
			return cmd.Run()
		},
	}

	exploreCmd.Flags().BoolVar(&cmd.ShowOnly, "show-only", false, "Show directory path only")
	return exploreCmd
}

func (cmd *ExploreCmd) Run() error {
	distName := cmd.EnvName

	if !wsllib.WslIsDistributionRegistered(distName) {
		return fmt.Errorf("environment %s not available. Enter it first", cmd.EnvName)
	}

	distPath, err := wsl.ExecuteSilently(distName, "wslpath -w .")
	if err != nil {
		return err
	}
	distPath = strings.TrimSpace(distPath)

	if cmd.ShowOnly {
		log.Info(distPath)
		return nil
	}

	proc := exec.Command("explorer.exe " + distPath)
	return proc.Run()
}
