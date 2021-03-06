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

	log "github.com/sirupsen/logrus"
)

var (
	ErrorsOnly    bool
	EnableVerbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "weasel",
	Long: "Tool for managing WSL distributions like docker containers.",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetOutput(os.Stdout)
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp:       true,
			ForceColors:            true,
			DisableLevelTruncation: true,
		})

		if ErrorsOnly {
			log.SetLevel(log.ErrorLevel)
		} else if EnableVerbose {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.PersistentFlags().BoolVar(&ErrorsOnly, "errors-only", false, "Only show errors")
	rootCmd.PersistentFlags().BoolVarP(&EnableVerbose, "verbose", "v", false, "Enable verbose output")

	cobra.CheckErr(rootCmd.Execute())
}
