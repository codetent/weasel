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

	"github.com/codetent/weasel/pkg/weasel/store"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List built distributions",
	Long:  "List all distributions built available for creating instances from them.",
	Run: func(cmd *cobra.Command, args []string) {
		code := func() int {
			dists, err := store.GetRegisteredDistributions()
			if err != nil {
				log.Errorf("Error reading distributions: %v", err)
				return 1
			}

			writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			fmt.Fprintln(writer, "ID")

			for _, dist := range dists {
				fmt.Fprintln(writer, dist.Id)
			}

			writer.Flush()
			return 0
		}()
		os.Exit(code)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
