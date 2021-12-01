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
	"path/filepath"

	//"fmt"

	"os"

	"github.com/spf13/cobra"

	"github.com/codetent/weasel/pkg/weasel/store"
	"github.com/codetent/weasel/pkg/weasel/wsl"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Extract parameters
		distId := args[0]
		instanceId, _ := cmd.Flags().GetString("name")
		autoRemove, _ := cmd.Flags().GetBool("rm")
		user, _ := cmd.Flags().GetString("user")
		attachStdin, _ := cmd.Flags().GetBool("stdin")

		if instanceId == "" {
			instanceId = RandomId()
		}

		// Get path to distribution archive
		distArchive, err := store.GetRegisteredDistribution(distId)
		if err != nil {
			panic(err)
		}
		if distArchive == "" {
			fmt.Printf("Distribution with id '%s' not found.\n", distId)
			os.Exit(1)
		}

		// Get path to workspace
		workspaceRoot, err := store.GetWorkspaceRoot()
		if err != nil {
			panic(err)
		}
		distWorkspace := filepath.Join(workspaceRoot, instanceId)
		err = os.MkdirAll(distWorkspace, os.ModePerm)
		if err != nil {
			panic(err)
		}
		if autoRemove {
			defer os.RemoveAll(distWorkspace)
		}

		// Load distribution archive into WSL
		err = wsl.Import(instanceId, distWorkspace, distArchive)
		if err != nil {
			panic(err)
		}

		// Register created WSL instance
		err = store.RegisterInstance(instanceId, distId)
		if err != nil {
			panic(err)
		}

		// Run instance
		err = wsl.Run(instanceId, &wsl.RunOpts{
			User:        user,
			AttachStdin: attachStdin,
		}, args[1:]...)

		// Fixme check error

		// Remove instance after using it
		if autoRemove {
			err = wsl.Unregister(instanceId)
			if err != nil {
				panic(err)
			}

			err = store.UnregisterInstance(instanceId)
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP("user", "u", "", "Sets the username for the specified command.")
	runCmd.Flags().String("name", "", "Sets the name of the instance.")
	runCmd.Flags().BoolP("stdin", "i", false, "Attach to STDIN.")
	runCmd.Flags().Bool("rm", false, "Automatically remove when it exits.")
}

func RandomId() string {
	data := make([]byte, 10)
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}
	return fmt.Sprintf("%x", sha256.Sum256(data))[:7]
}
