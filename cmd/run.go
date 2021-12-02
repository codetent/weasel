/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/codetent/weasel/pkg/weasel/store"
	"github.com/codetent/weasel/pkg/weasel/wsl"

	log "github.com/sirupsen/logrus"
)

var (
	User        string
	AttachStdin bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a new instance of a distribution",
	Long:  "Create a new instance of an available distribution. This command also opens up an interactive session to it.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		distId := args[0]
		instanceId := RandomId()

		// Get path to distribution archive
		distArchive, err := store.GetRegisteredDistribution(distId)
		if err != nil {
			log.Fatalf("Error finding distribution: %v", err)
		}
		if distArchive == "" {
			log.Fatalf("Distribution with id '%s' not found", distId)
		}

		// Get path to workspace
		workspaceRoot, err := store.GetWorkspaceRoot()
		if err != nil {
			log.Fatalf("Error location workspace root: %v", err)
		}
		distWorkspace := filepath.Join(workspaceRoot, instanceId)
		err = os.MkdirAll(distWorkspace, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating workspace: %v", err)
		}
		defer os.RemoveAll(distWorkspace)

		// Load distribution archive into WSL
		err = wsl.Import(instanceId, distWorkspace, distArchive)
		if err != nil {
			log.Fatalf("Error loading instance: %v", err)
		}
		defer func() {
			err := wsl.Unregister(instanceId)
			if err != nil {
				log.Errorf("Error unloading instance: %v", err)
			}
		}()

		// Register created WSL instance
		err = store.RegisterInstance(instanceId, distId)
		if err != nil {
			log.Fatalf("Error registering instance: %v", err)
		}
		defer func() {
			log.Info("here")
			err := store.UnregisterInstance(instanceId)
			if err != nil {
				log.Errorf("Error unregistering instance: %v", err)
			}
		}()

		// Run instance
		code, err := wsl.Run(instanceId, &wsl.RunOpts{
			User:        User,
			AttachStdin: AttachStdin,
		}, args[1:]...)
		if code < 0 {
			if err != nil {
				log.Fatalf("Error running instance: %v", err)
			}
			code = 1
		}

		os.Exit(code)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&User, "user", "u", "", "Sets the username for the specified command.")
	runCmd.Flags().BoolVarP(&AttachStdin, "stdin", "i", false, "Attach to STDIN.")
}

func RandomId() string {
	data := make([]byte, 10)
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}
	return fmt.Sprintf("%x", sha256.Sum256(data))[:7]
}
