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
	"path"

	"github.com/aelsabbahy/dargs/run"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// generateBinsCmd represents the generateBins command
var generateBinsCmd = &cobra.Command{
	Use:   "generate-bins",
	Short: "Generate wrapper executables for commands",
	Long: `Generate wraper executables for commands.

These can be used as a replacement for:
alias cmd='dargs run -- cmd'

The advantage of these wrapper scripts over aliases is the ability
to use them as a drop-in replacement for the orignal command whithout
worrying about aliases or alias support if being called from another tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := run.GenerateBins(args); err != nil {
			logrus.StandardLogger().Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(generateBinsCmd)
	var binDir string
	if home, err := homedir.Dir(); err == nil {
		binDir = path.Join(home, "bin")
	}
	generateBinsCmd.Flags().StringP("bin-out", "o", binDir, "output directory")
	generateBinsCmd.Flags().BoolP("force-bins", "f", false, "force override files")
	viper.BindPFlags(generateBinsCmd.Flags())
}
