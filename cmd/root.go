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
	"os"
	"path/filepath"

	"github.com/a8m/envsubst"
	"github.com/spf13/cobra"
)

var (
	paths   []string
	version = "1.0.0" // Version can be overridden at build time
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "subenv",
	Short: "Substitutes the values of environment variables.",
	Long: `The plugin allows to substitute the values of environment variables within a CI/CD pipeline.

It supports:
- Single file substitution
- Multiple files substitution
- Directory (recursive) substitution
- Mixed files and directories

Environment variables can be referenced using $VAR or ${VAR} syntax.`,
	Version: version,
	RunE:    run,
}

// run is the main execution function for the root command
func run(cmd *cobra.Command, args []string) error {
	if len(paths) == 0 {
		return fmt.Errorf("at least one file or directory path must be specified using -f flag")
	}

	for _, path := range paths {
		if err := processPath(path); err != nil {
			return fmt.Errorf("error processing path %s: %w", path, err)
		}
	}

	return nil
}

// processPath processes a single path (file or directory)
func processPath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat path: %w", err)
	}

	if info.IsDir() {
		return processDirectory(path)
	}

	return expandEnv(path)
}

// processDirectory recursively processes all files in a directory
func processDirectory(dirPath string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %w", path, err)
		}

		// Skip directories, only process files
		if info.IsDir() {
			return nil
		}

		if err := expandEnv(path); err != nil {
			return fmt.Errorf("error processing file %s: %w", path, err)
		}

		return nil
	})
}

// expandEnv reads a file, substitutes environment variables, and writes back
func expandEnv(filePath string) error {
	// Read file and substitute environment variables
	newContent, err := envsubst.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read and substitute file: %w", err)
	}

	// Write the substituted content back to the file
	// Use 0644 permissions (rw-r--r--) instead of 0777 for security
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()

	// Define flags
	rootCmd.Flags().StringArrayVarP(&paths, "file", "f", []string{},
		"specify path to values file or directory. You can configure the flag multiple times for different files/directories.")

	// Mark the file flag as required
	if err := rootCmd.MarkFlagRequired("file"); err != nil {
		fmt.Fprintf(os.Stderr, "Error marking flag as required: %v\n", err)
		os.Exit(1)
	}
}
