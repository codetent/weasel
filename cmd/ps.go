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
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/codetent/weasel/pkg/weasel/store"
	"github.com/codetent/weasel/pkg/weasel/wsl"

	log "github.com/sirupsen/logrus"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List running instances",
	Long:  "Shows all running instances together with their distribution.",
	Run: func(cmd *cobra.Command, args []string) {
		code := func() int {
			allInstances, err := store.GetRegisteredInstances()
			if err != nil {
				log.Errorf("Error reading registered instances: %v", err)
				return 1
			}

			runningInstances, err := wsl.ListRunning()
			if err != nil {
				log.Errorf("Error reading running instances: %v", err)
				return 1
			}

			writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			fmt.Fprintln(writer, "ID\tDISTRIBUTION")

			for _, inst := range allInstances {
				running := false
				for _, other := range runningInstances {
					if other == inst.Id {
						running = true
						break
					}
				}

				if running {
					fmt.Fprintf(writer, "%s\t%s\n", inst.Id, inst.Distribution)
				}
			}

			writer.Flush()
			return 0
		}()
		os.Exit(code)
	},
}

func init() {
	rootCmd.AddCommand(psCmd)
}
