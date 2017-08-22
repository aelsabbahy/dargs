// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/aelsabbahy/dargs/run"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// generateCompletionsCmd represents the generateCompletions command
var generateCompletionsCmd = &cobra.Command{
	Use:   "generate-completions",
	Short: "Generate bash/zsh completion files",
	Long: `Generate bash/zsh completion files.
These files need to be sourced in by your ~/.bashrc or ~/.zshrc
to allow completions to work.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := run.GenerateCompletions(args); err != nil {
			logrus.StandardLogger().Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(generateCompletionsCmd)
	generateCompletionsCmd.Flags().StringP("completions-out", "o", "", "output directory")
	generateCompletionsCmd.Flags().StringP("shell", "s", "bash", "shell to use (default is bash)")
	generateCompletionsCmd.Flags().BoolP("force-completions", "f", false, "force override files")
	viper.BindPFlags(generateCompletionsCmd.Flags())
}
