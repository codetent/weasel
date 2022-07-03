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

	"github.com/spf13/cobra"
	"github.com/yuk7/wsllib-go"
)

type RemoveCmd struct {
	EnvName string
}

func NewRemoveCmd() *cobra.Command {
	cmd := &RemoveCmd{}

	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			cmd.EnvName = args[0]
			return cmd.Run()
		},
	}

	return removeCmd
}

func (cmd *RemoveCmd) Run() error {
	distName := cmd.EnvName

	if !wsllib.WslIsDistributionRegistered(distName) {
		return fmt.Errorf("environment %s not available. Enter it first", cmd.EnvName)
	}

	return wsllib.WslUnregisterDistribution(distName)
}
