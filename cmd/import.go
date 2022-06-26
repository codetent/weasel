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
	"github.com/codetent/weasel/pkg/weasel/docker"
	"github.com/codetent/weasel/pkg/weasel/wsl"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
	"github.com/spf13/cobra"
	"github.com/yuk7/wsllib-go"
)

var (
	ImportDisableCache bool
	ImportForce        bool
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import distribution",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		imageRef := args[0]
		distName := imageRef
		if len(args) > 1 {
			distName = args[1]
		}

		// Check distribution state
		if wsllib.WslIsDistributionRegistered(distName) {
			if ImportForce {
				fmt.Printf("Distribution with name already '%s' found. Unregistering it.\n", distName)

				err := wsllib.WslUnregisterDistribution(distName)
				if err != nil {
					return err
				}
			} else {
				fmt.Printf("Distribution '%s' already installed\n", distName)
				return nil
			}
		}

		// Pull docker image
		stream, err := docker.ImagePull(imageRef)
		if err != nil {
			return err
		}
		defer stream.Close()
		termFd, isTerm := term.GetFdInfo(os.Stdout)
		jsonmessage.DisplayJSONMessagesStream(stream, os.Stdout, termFd, isTerm, nil)

		// Lock cache
		lock, err := cache.GetLock()
		if err != nil {
			return err
		}
		defer lock.Unlock()

		// Export docker image rootfs
		archivePath, err := cache.GetDistPath(imageRef)
		if err != nil {
			return err
		}

		err = docker.ImageExport(imageRef, archivePath)
		if err != nil {
			return err
		}
		defer os.Remove(archivePath)

		// Get path to workspace
		distWorkspace, err := cache.GetWorkspacePath(distName)
		if err != nil {
			return err
		}

		// Import distribution
		err = wsl.Import(distName, distWorkspace, archivePath)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().BoolVar(&ImportDisableCache, "no-cache", false, "Do not use cache when building the image.")
	importCmd.Flags().BoolVar(&ImportForce, "force", false, "Overwrite already existing distribution")
}
