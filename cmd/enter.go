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
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codetent/weasel/pkg/weasel/config"
	"github.com/codetent/weasel/pkg/weasel/utils"
	"github.com/codetent/weasel/pkg/weasel/wsl"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
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
		Short: "Enter existing distribution",
		Args:  cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			cmd.EnvName = args[0]
			return cmd.Run()
		},
	}

	enterCmd.Flags().BoolVarP(&cmd.Recreate, "recreate", "r", false, "Recreate distribution")
	return enterCmd
}

func (cmd *EnterCmd) Run() error {
	configPath, err := config.Locate()
	if err != nil {
		return err
	}
	log.Infof("configuration found at %s", configPath)

	config, err := config.Load(configPath)
	if err != nil {
		return err
	}

	if _, ok := config.Environments[cmd.EnvName]; !ok {
		return fmt.Errorf("undefined environment '%s'", cmd.EnvName)
	}

	distName := config.Name + "-" + cmd.EnvName
	envExists := wsllib.WslIsDistributionRegistered(distName)
	if cmd.Recreate && envExists {
		log.Warnf("distribution with name already '%s' exists. recreating it", distName)

		err := wsllib.WslUnregisterDistribution(distName)
		if err != nil {
			return err
		}

		envExists = false
	}

	if !envExists {
		imageRef, err := name.ParseReference(config.Environments[cmd.EnvName].Image)
		if err != nil {
			return err
		}
		image, err := remote.Image(imageRef)
		if err != nil {
			return err
		}
		imageDigest, err := image.Digest()
		if err != nil {
			return err
		}

		weaselPath := filepath.Join(filepath.Dir(configPath), ".weasel")
		cachePath := filepath.Join(weaselPath, "cache", imageRef.Context().Name())
		archivePath := filepath.Join(cachePath, imageDigest.Hex[:12]+".tar.gz")

		if _, err := os.Stat(archivePath); os.IsNotExist(err) {
			err = os.MkdirAll(cachePath, os.ModePerm)
			if err != nil {
				return err
			}

			pullPath, err := ioutil.TempDir("", "weasel")
			if err != nil {
				return err
			}
			defer os.RemoveAll(pullPath)

			tarPath := filepath.Join(pullPath, "image.tar.gz")

			log.Infoln("pulling image tarball")
			log.Debugf("storing tarball at '%s'", tarPath)

			err = tarball.WriteToFile(tarPath, imageRef, image)
			if err != nil {
				return err
			}

			untaredPath := filepath.Join(pullPath, "content")
			err = utils.UntarPattern(tarPath, untaredPath)
			if err != nil {
				return err
			}

			archivePathCandidates, err := filepath.Glob(filepath.Join(untaredPath, "*.tar.gz"))
			if err != nil {
				return err
			} else if len(archivePathCandidates) == 0 {
				return fmt.Errorf("archive not found")
			}

			utils.CopyFile(archivePathCandidates[0], archivePath)
		} else {
			log.Info("image tarball already found in cache")
		}

		workspacePath := filepath.Join(weaselPath, "workspaces", imageRef.Context().Name(), imageDigest.Hex[:12])
		err = os.MkdirAll(workspacePath, os.ModePerm)
		if err != nil {
			return err
		}

		log.Infof("importing distribution into WSL as '%s'", distName)
		log.Debugf("workspace of distribution at '%s'", workspacePath)

		err = wsl.Import(distName, workspacePath, archivePath)
		if err != nil {
			return err
		}
	}

	wsllib.WslLaunchInteractive(distName, "", true)
	return nil
}
