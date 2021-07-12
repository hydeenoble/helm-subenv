/*
Copyright © 2020 hydeenoble
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
	"github.com/spf13/cobra"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

var cfgFile string
var file string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "subenv",
	Short: "Substitutes the values of environment variables.",
	Long: `The plugin allows to substitue the values of environment variables withing a CICD pipeline.`,
		
		Run: func(cmd *cobra.Command, args []string) { 
			exCmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("envsubst < %[1]s > %[1]s.temp && mv %[1]s.temp %[1]s", file))
			stdout, err := exCmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Print(string(stdout))
		},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tctl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("file", "f", false, "Please specify values file.")
	rootCmd.Flags().StringVarP(&file, "file", "f", "", "specify values file.")
	rootCmd.MarkFlagRequired("file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".tctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}