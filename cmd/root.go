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

// Package cmd defines the helm-subenv CLI commands.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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
func run(_ *cobra.Command, _ []string) error {
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
	// Read the original file content
	originalContent, err := os.ReadFile(filePath) // nolint:gosec
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Read file and substitute environment variables
	newContent, err := envsubst.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read and substitute file: %w", err)
	}

	// Validate for potential issues
	if err := validateSubstitution(filePath, string(originalContent), string(newContent)); err != nil {
		// Log warnings but don't fail - substitution succeeded
		fmt.Fprintf(os.Stderr, "⚠️  Warning processing %s: %v\n", filePath, err)
	}

	// Write the substituted content back to the file
	// Use 0o600 permissions (rw-------) for security
	if err := os.WriteFile(filePath, newContent, 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// validateSubstitution checks for common issues during variable substitution
func validateSubstitution(_ string, originalContent, newContent string) error {
	// Check for empty variable values (potential unexported variables)
	emptyVars := detectEmptyVariables(originalContent, newContent)
	if len(emptyVars) > 0 {
		return fmt.Errorf("variables were replaced with empty strings (may not be exported): %s. "+
			"Make sure to use 'export' when setting environment variables. "+
			"See: https://github.com/hydeenoble/helm-subenv#known-limitations",
			strings.Join(emptyVars, ", "))
	}

	// Check for patterns that might indicate bash arrays
	if detectBashArrayPatterns(originalContent) {
		return fmt.Errorf("file may contain bash array variables which are not supported. " +
			"Bash arrays cannot be substituted. Use space-separated or comma-separated strings instead. " +
			"See: https://github.com/hydeenoble/helm-subenv#known-limitations")
	}

	return nil
}

// detectEmptyVariables finds environment variables that were replaced with empty strings
func detectEmptyVariables(originalContent, _ string) []string {
	var emptyVars []string
	seenVars := make(map[string]bool)

	// Find all variable references in original content
	varPattern := regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_]*)\}|\$([a-zA-Z_][a-zA-Z0-9_]*)`)
	matches := varPattern.FindAllStringSubmatch(originalContent, -1)

	for _, match := range matches {
		varName := match[1]
		if varName == "" {
			varName = match[2]
		}

		if seenVars[varName] {
			continue
		}
		seenVars[varName] = true

		// Check if variable exists and is not empty
		value := os.Getenv(varName)
		if value == "" {
			// Check if the substitution actually resulted in an empty string
			// by looking for the variable pattern in both contents
			if strings.Contains(originalContent, "$"+varName) || strings.Contains(originalContent, "${"+varName+"}") {
				emptyVars = append(emptyVars, varName)
			}
		}
	}

	return emptyVars
}

// detectBashArrayPatterns detects patterns that might indicate bash array definitions
func detectBashArrayPatterns(content string) bool {
	// Pattern to detect bash array-like environment variable values
	// Looks for patterns like VAR=(...) or VAR=[ ... ]
	arrayPattern := regexp.MustCompile(`\$\{[a-zA-Z_][a-zA-Z0-9_]*\}\s*=\s*\(.*\)`)
	return arrayPattern.MatchString(content)
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
