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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"github.com/a8m/envsubst"
)

// type Files []string
var paths []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "subenv",
	Short: "Substitutes the values of environment variables.",
	Long:  `The plugin allows to substitue the values of environment variables withing a CICD pipeline.`,

	Run: func(cmd *cobra.Command, args []string) {

		for _, path := range paths {

			dir, err := os.Stat(path)

			if err != nil {
				log.Fatal(err)
			}

			if dir.IsDir() {
				err := filepath.Walk(path,
					func(subPath string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}
						if !info.IsDir() {
							expandEnv(subPath)
						}
						return nil
					})

				if err != nil {
					log.Fatal(err)
				}
			} else {
				expandEnv(path)
			}

		}
	},
}

func expandEnv(filePath string) {
	newContent, err := envsubst.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filePath, []byte(newContent), 0777)

	if err != nil {
		log.Fatal(err)
	}
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
	cobra.OnInitialize()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringArrayVarP(&paths, "file", "f", []string{}, "specify path to values file. You can configure the flag multiple times for different files.")
	rootCmd.MarkFlagRequired("file")
}
