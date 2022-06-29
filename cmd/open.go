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
	"syscall"

	"github.com/codetent/weasel/pkg/weasel/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yuk7/wsllib-go"
)

type OpenCmd struct {
	EnvName string
}

func NewOpenCmd() *cobra.Command {
	cmd := &OpenCmd{}

	openCmd := &cobra.Command{
		Use:   "open",
		Short: "Open environment folder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			cmd.EnvName = args[0]
			return cmd.Run()
		},
	}

	return openCmd
}

func (cmd *OpenCmd) Run() error {
	configFile, err := config.LocateConfigFile()
	if err == nil {
		log.Debugf("Configuration located at %s", configFile.Path)
	} else {
		return err
	}

	configContent, err := configFile.Content()
	if err == nil {
		log.Debug("Configuration loaded successfully")
	} else {
		return err
	}

	if _, ok := configContent.Environments[cmd.EnvName]; !ok {
		return fmt.Errorf("undefined environment %s", cmd.EnvName)
	}

	distName := configContent.Name + "-" + cmd.EnvName

	if !wsllib.WslIsDistributionRegistered(distName) {
		return fmt.Errorf("environment %s not available. Enter it first", cmd.EnvName)
	}

	devNull := syscall.Handle(0)
	_, err = wsllib.WslLaunch(distName, "explorer.exe .", false, devNull, devNull, devNull)
	return err
}
