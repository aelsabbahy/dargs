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
	"fmt"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dargs",
	Short: "Dynamic CLI arguments and completions",
	Long: `Dargs is a tool that allows you to define dynamic argument
replacements and completions for any CLI command.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringP("config", "", "", "config file (default is $HOME/.dargs.yaml)")
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "print debug output")

	viper.BindPFlags(RootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("DARGS")
	viper.AutomaticEnv() // read in environment variables that match

	logrus.SetFormatter(&prefixed.TextFormatter{
		DisableTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}
	cfgFile := viper.GetString("config")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		viper.Set("root", path.Dir(cfgFile))
		viper.Set("config", cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.Set("root", home)
		viper.Set("config", path.Join(home, ".dargs.yml"))

		//// Search config in home directory with name ".dargs" (without extension).
		//viper.AddConfigPath(home)
		//viper.SetConfigName(".dargs")
	}

	// If a config file is found, read it in.
	//viper.ReadInConfig()
	viper.Set("cacheFile", path.Join(viper.GetString("root"), ".dargs", "cache.yml"))
	viper.Set("downloadRoot", path.Join(viper.GetString("root"), ".dargs", "downloads"))
}
