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
	"strings"

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
	ImageRef     string
	EnvName      string
	Recreate     bool
	RegisterOnly bool
}

func NewEnterCmd() *cobra.Command {
	cmd := &EnterCmd{}

	enterCmd := &cobra.Command{
		Use:   "enter",
		Short: "Enter environment",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			cmd.ImageRef = args[0]

			if len(args) > 1 {
				cmd.EnvName = args[1]
			} else {
				cmd.EnvName = ""
			}

			return cmd.Run()
		},
	}

	enterCmd.Flags().BoolVarP(&cmd.Recreate, "recreate", "r", false, "Recreate environment")
	enterCmd.Flags().BoolVar(&cmd.RegisterOnly, "register", false, "Only register environment without entering it")
	return enterCmd
}

func (cmd *EnterCmd) Run() error {
	imageRef, err := name.ParseReference(cmd.ImageRef)
	if err != nil {
		return err
	}

	distName := cmd.EnvName
	if distName == "" {
		imageRefParts := strings.Split(imageRef.Context().Name(), "/")
		distName = imageRefParts[len(imageRefParts)-1]
	}

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
		image, err := remote.Image(imageRef, remote.WithAuthFromKeychain(authn.DefaultKeychain))
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				return fmt.Errorf("requested image %s not found", cmd.ImageRef)
			} else {
				return err
			}
		}
		imageDigest, err := image.Digest()
		if err != nil {
			return err
		}

		log.Infof("Setting up environment using %s:%s", imageRef.Context().Name(), imageDigest.Hex)

		imagePath, err := config.GetImageCachePath(imageRef, imageDigest)
		if err != nil {
			return err
		}

		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(imagePath), os.ModePerm)
			if err != nil {
				return err
			}

			log.Debug("Exporting image rootfs as tarball")

			err = oci.ExportRootFs(image, imageRef, imagePath)
			if err != nil {
				return err
			}
		} else {
			log.Debug("Image tarball already found in cache")
		}

		workspacePath, err := config.GetWorkspaceCachePath(imageRef, imageDigest)
		if err != nil {
			return err
		}
		workspaceVhdx := filepath.Join(workspacePath, "ext4.vhdx")

		if _, err := os.Stat(workspaceVhdx); os.IsNotExist(err) {
			err = os.MkdirAll(workspacePath, os.ModePerm)
			if err != nil {
				return err
			}

			log.Infof("Importing environment into WSL as %s", distName)

			err = wsl.Import(distName, workspacePath, imagePath)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("foreign virtual disk in cache at %s", workspaceVhdx)
		}
	}

	if cmd.RegisterOnly {
		return nil
	}

	_, err = wsllib.WslLaunchInteractive(distName, "", true)
	return err
}
