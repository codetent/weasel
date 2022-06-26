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
package cache

import (
	"github.com/codetent/weasel/pkg/weasel/cache"
	"github.com/spf13/cobra"
)

type CleanCmd struct{}

func NewCleanCmd() *cobra.Command {
	cmd := &CleanCmd{}

	cleanCmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean weasel cache",
		Args:  cobra.NoArgs,
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.Run()
		},
	}

	return cleanCmd
}

func (cmd *CleanCmd) Run() error {
	err := cache.CleanDistCache()
	if err != nil {
		return err
	}

	return cache.CleanWorkspaceCache()
}
