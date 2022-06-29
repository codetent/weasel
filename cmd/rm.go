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

	"github.com/codetent/weasel/pkg/weasel/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yuk7/wsllib-go"
)

type RmCmd struct {
	EnvName string
}

func NewRmCmd() *cobra.Command {
	cmd := &RmCmd{}

	rmCmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove distribution",
		Args:  cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			cmd.EnvName = args[0]
			return cmd.Run()
		},
	}

	return rmCmd
}

func (cmd *RmCmd) Run() error {
	configFile, err := config.LocateConfigFile()
	if err != nil {
		return err
	}
	log.Infof("configuration found at %s", configFile.Path)

	config, err := configFile.Content()
	if err != nil {
		return err
	}

	if _, ok := config.Environments[cmd.EnvName]; !ok {
		return fmt.Errorf("undefined environment '%s'", cmd.EnvName)
	}

	distName := config.Name + "-" + cmd.EnvName

	if !wsllib.WslIsDistributionRegistered(distName) {
		return fmt.Errorf("distribution '%s' not found", distName)
	}

	err = wsllib.WslUnregisterDistribution(distName)
	if err != nil {
		return err
	}

	return nil
}
