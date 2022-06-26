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

	"github.com/spf13/cobra"

	"github.com/codetent/weasel/cmd/cache"
	log "github.com/sirupsen/logrus"
)

type RootCmd struct {
	ErrorsOnly bool
	Verbose    bool
}

func NewRootCmd() *cobra.Command {
	cmd := &RootCmd{}

	rootCmd := &cobra.Command{
		Use:  "weasel",
		Long: "Tool for managing WSL distributions like docker containers.",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		PersistentPreRunE: func(cobraCmd *cobra.Command, args []string) error {
			return cmd.PreRun()
		},
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return cobraCmd.Help()
		},
	}

	rootCmd.PersistentFlags().BoolVar(&cmd.ErrorsOnly, "errors-only", false, "Only show errors")
	rootCmd.PersistentFlags().BoolVarP(&cmd.Verbose, "verbose", "v", false, "Enable verbose output")
	return rootCmd
}

func (cmd *RootCmd) PreRun() error {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp:       true,
		ForceColors:            true,
		DisableLevelTruncation: true,
	})

	if cmd.ErrorsOnly {
		log.SetLevel(log.ErrorLevel)
	} else if cmd.Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	return nil
}

func Execute() {
	rootCmd := NewRootCmd()

	rootCmd.AddCommand(NewImportCmd())
	rootCmd.AddCommand(NewRmCmd())
	rootCmd.AddCommand(NewEnterCmd())
	rootCmd.AddCommand(cache.NewCacheCmd())

	cobra.CheckErr(rootCmd.Execute())
}
