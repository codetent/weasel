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
	"os"

	"github.com/codetent/weasel/pkg/weasel/store"
	"github.com/codetent/weasel/pkg/weasel/wsl"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove instance",
	Long:  "Remove an instance that is in any state by its id.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		code := func() int {
			instanceId := args[0]

			err := store.UnregisterInstance(instanceId)
			if err != nil {
				log.Errorf("Error unregistering instance: %v", err)
			}

			err = wsl.Unregister(instanceId)
			if err != nil {
				log.Errorf("Error unloading instance: %v", err)
			}

			if err == nil {
				return 0
			} else {
				return 1
			}
		}()
		os.Exit(code)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
