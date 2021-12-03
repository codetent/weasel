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
	"os"

	"github.com/codetent/weasel/pkg/weasel/store"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

// pruneCmd represents the prune command
var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune distribution",
	Long:  "Remove distribution(s) including cache entries and built files.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		code := func() int {
			distIds := args

			if len(args) == 0 {
				dists, err := store.GetRegisteredDistributions()
				if err != nil {
					log.Errorf("Error reading distributions: %v", err)
					return 1
				}

				for _, dist := range dists {
					distIds = append(distIds, dist.Id)
				}
			}

			for _, id := range distIds {
				err := store.UnregisterDistribution(id)
				if err != nil {
					log.Errorf("Error pruning distribution: %v", err)
					return 1
				}
			}

			return 0
		}()
		os.Exit(code)
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)
}
