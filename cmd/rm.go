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
	"os"

	"github.com/codetent/weasel/pkg/weasel/cache"
	"github.com/spf13/cobra"
	"github.com/yuk7/wsllib-go"
)

type RmCmd struct {
	DistName string
}

func NewRmCmd() *cobra.Command {
	cmd := &RmCmd{}

	rmCmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove distribution",
		Args:  cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			cmd.DistName = args[0]
			return cmd.Run()
		},
	}

	return rmCmd
}

func (cmd *RmCmd) Run() error {
	if !wsllib.WslIsDistributionRegistered(cmd.DistName) {
		return fmt.Errorf("distribution '%s' not found", cmd.DistName)
	}

	err := wsllib.WslUnregisterDistribution(cmd.DistName)
	if err != nil {
		return err
	}

	distWorkspace, err := cache.GetWorkspacePath(cmd.DistName)
	if err != nil {
		return err
	}

	return os.RemoveAll(distWorkspace)
}
