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

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [flags] -- command args..",
	Short: "Transform arguments and run command",
	Long:  `Transform arguments and run command`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := run.Run(args); err != nil {
			logrus.StandardLogger().Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("noop", "n", false, "noop, only print what would have been executed")
	runCmd.Flags().BoolP("verbose", "v", false, "verbose, print transformed command before executing")
	runCmd.Flags().StringSliceP("transformers", "t", []string{}, "comma separated list of transformers to use")
	viper.BindPFlags(runCmd.Flags())
}
