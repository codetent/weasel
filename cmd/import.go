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

	log "github.com/sirupsen/logrus"

	"github.com/codetent/weasel/pkg/weasel/cache"
	"github.com/codetent/weasel/pkg/weasel/docker"
	"github.com/codetent/weasel/pkg/weasel/wsl"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
	"github.com/spf13/cobra"
	"github.com/yuk7/wsllib-go"
)

type ImportCmd struct {
	ImageRef string
	DistName string
	Force    bool
}

func NewImportCmd() *cobra.Command {
	cmd := &ImportCmd{}

	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Import distribution",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			cmd.ImageRef = args[0]

			cmd.DistName = cmd.ImageRef
			if len(args) > 1 {
				cmd.DistName = args[1]
			}

			return cmd.Run()
		},
	}

	importCmd.Flags().BoolVarP(&cmd.Force, "force", "f", false, "Overwrite already existing distribution")
	return importCmd
}

func (cmd *ImportCmd) Run() error {
	if wsllib.WslIsDistributionRegistered(cmd.DistName) {
		if cmd.Force {
			log.Warnf("Distribution with name already '%s' found. Unregistering it.", cmd.DistName)
			fmt.Println()

			err := wsllib.WslUnregisterDistribution(cmd.DistName)
			if err != nil {
				return err
			}
		} else {
			log.Warnf("Distribution '%s' already installed", cmd.DistName)
			return nil
		}
	}

	log.Infof("Pulling docker image '%s'", cmd.ImageRef)

	stream, err := docker.ImagePull(cmd.ImageRef)
	if err != nil {
		return err
	}
	defer stream.Close()
	writer := log.StandardLogger().Out
	termFd, isTerm := term.GetFdInfo(writer)
	jsonmessage.DisplayJSONMessagesStream(stream, writer, termFd, isTerm, nil)

	fmt.Println()
	log.Debugln("Locking cache")

	lock, err := cache.GetLock()
	if err != nil {
		return err
	}
	defer lock.Unlock()

	archivePath, err := cache.GetDistPath(cmd.ImageRef)
	if err != nil {
		return err
	}

	log.Infoln("Creating distribution using image rootfs")
	log.Debugf("Storing distribution at '%s'", archivePath)

	err = docker.ImageExport(cmd.ImageRef, archivePath)
	if err != nil {
		return err
	}
	defer os.Remove(archivePath)

	distWorkspace, err := cache.GetWorkspacePath(cmd.DistName)
	if err != nil {
		return err
	}

	log.Infof("Importing distribution into WSL as '%s'", cmd.DistName)
	log.Debugf("Workspace of distribution at '%s'", distWorkspace)

	err = wsl.Import(cmd.DistName, distWorkspace, archivePath)
	if err != nil {
		return err
	}

	return nil
}
