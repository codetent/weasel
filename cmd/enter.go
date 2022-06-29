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
	"path/filepath"

	"github.com/codetent/weasel/pkg/weasel/config"
	"github.com/codetent/weasel/pkg/weasel/oci"
	"github.com/codetent/weasel/pkg/weasel/wsl"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yuk7/wsllib-go"
)

type EnterCmd struct {
	EnvName  string
	Recreate bool
}

func NewEnterCmd() *cobra.Command {
	cmd := &EnterCmd{}

	enterCmd := &cobra.Command{
		Use:   "enter",
		Short: "Enter existing environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			cmd.EnvName = args[0]
			return cmd.Run()
		},
	}

	enterCmd.Flags().BoolVarP(&cmd.Recreate, "recreate", "r", false, "Recreate environment")
	return enterCmd
}

func (cmd *EnterCmd) Run() error {
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
	envExists := wsllib.WslIsDistributionRegistered(distName)

	if envExists && cmd.Recreate {
		log.Warn("Recreating already existing environment")

		err := wsllib.WslUnregisterDistribution(distName)
		if err != nil {
			return err
		}

		envExists = false
	}

	if envExists {
		log.Debug("Loading already existing environment")
	} else {
		imageRawRef := configContent.Environments[cmd.EnvName].Image
		imageRef, err := name.ParseReference(imageRawRef)
		if err != nil {
			return err
		}
		image, err := remote.Image(imageRef, remote.WithAuthFromKeychain(authn.DefaultKeychain))
		if err != nil {
			return err
		}
		imageDigest, err := image.Digest()
		if err != nil {
			return err
		}

		log.Infof("Setting up environment using %s:%s", imageRef.Context().Name(), imageDigest.Hex)

		archivePath := config.GetArchiveCachePath(configFile, imageRef, imageDigest)

		if _, err := os.Stat(archivePath); os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(archivePath), os.ModePerm)
			if err != nil {
				return err
			}

			log.Debug("Exporting image rootfs as tarball")

			err = oci.ExportRootFs(image, imageRef, archivePath)
			if err != nil {
				return err
			}
		} else {
			log.Debug("Image tarball already found in cache")
		}

		workspacePath := config.GetWorkspaceCachePath(configFile, imageRef, imageDigest)
		workspaceVhdx := filepath.Join(workspacePath, "ext4.vhdx")

		if _, err := os.Stat(workspaceVhdx); os.IsNotExist(err) {
			err = os.MkdirAll(workspacePath, os.ModePerm)
			if err != nil {
				return err
			}

			log.Infof("Importing environment into WSL as %s", distName)

			err = wsl.Import(distName, workspacePath, archivePath)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("foreign virtual disk in cache at %s", workspaceVhdx)
		}
	}

	_, err = wsllib.WslLaunchInteractive(distName, "", true)
	return err
}
