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
	"github.com/spf13/cobra"
	"github.com/yuk7/wsllib-go"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove distribution",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		distName := args[0]

		if !wsllib.WslIsDistributionRegistered(distName) {
			return fmt.Errorf("distribution '%s' not found", distName)
		}

		err := wsllib.WslUnregisterDistribution(distName)
		if err != nil {
			return err
		}

		distWorkspace, err := cache.GetWorkspacePath(distName)
		if err != nil {
			return err
		}

		return os.RemoveAll(distWorkspace)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
