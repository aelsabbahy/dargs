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

// completionsCmd represents the completions command
var completionsCmd = &cobra.Command{
	Use:   "completions [flags] -- cmd prev_arg cur_arg",
	Short: "Generate completions for command",
	Long: `This command is used by bash/zsh completions to generate
completions for the given command.

It can be used standalone to debug completers`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := run.Completions(args); err != nil {
			logrus.StandardLogger().Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(completionsCmd)
	completionsCmd.Flags().StringSliceP("completers", "c", []string{}, "comma separated list of completers to use")
	viper.BindPFlags(completionsCmd.Flags())
}
