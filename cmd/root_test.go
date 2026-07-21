package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// TestExpandEnv tests the expandEnv function
func TestExpandEnv(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "helm-subenv-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test case 1: Simple environment variable substitution
	t.Run("SimpleSubstitution", func(t *testing.T) {
		testFile := filepath.Join(tmpDir, "test1.yaml")
		content := "image: $TEST_IMAGE\ntag: $TEST_TAG"
		
		// Set environment variables
		os.Setenv("TEST_IMAGE", "nginx")
		os.Setenv("TEST_TAG", "latest")
		defer os.Unsetenv("TEST_IMAGE")
		defer os.Unsetenv("TEST_TAG")

		// Write test file
		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		// Run expandEnv
		expandEnv(testFile)

		// Read result
		result, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("Failed to read result file: %v", err)
		}

		expected := "image: nginx\ntag: latest"
		if string(result) != expected {
			t.Errorf("Expected %q, got %q", expected, string(result))
		}
	})

	// Test case 2: Missing environment variable (should be empty)
	t.Run("MissingVariable", func(t *testing.T) {
		testFile := filepath.Join(tmpDir, "test2.yaml")
		content := "value: $NONEXISTENT_VAR"
		
		// Ensure variable doesn't exist
		os.Unsetenv("NONEXISTENT_VAR")

		// Write test file
		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		// Run expandEnv
		expandEnv(testFile)

		// Read result
		result, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("Failed to read result file: %v", err)
		}

		expected := "value: "
		if string(result) != expected {
			t.Errorf("Expected %q, got %q", expected, string(result))
		}
	})

	// Test case 3: Multiple variables in one line
	t.Run("MultipleVariables", func(t *testing.T) {
		testFile := filepath.Join(tmpDir, "test3.yaml")
		content := "url: $PROTOCOL://$HOST:$PORT"
		
		// Set environment variables
		os.Setenv("PROTOCOL", "https")
		os.Setenv("HOST", "example.com")
		os.Setenv("PORT", "8080")
		defer os.Unsetenv("PROTOCOL")
		defer os.Unsetenv("HOST")
		defer os.Unsetenv("PORT")

		// Write test file
		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		// Run expandEnv
		expandEnv(testFile)

		// Read result
		result, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("Failed to read result file: %v", err)
		}

		expected := "url: https://example.com:8080"
		if string(result) != expected {
			t.Errorf("Expected %q, got %q", expected, string(result))
		}
	})

	// Test case 4: Braces syntax ${VAR}
	t.Run("BracesSyntax", func(t *testing.T) {
		testFile := filepath.Join(tmpDir, "test4.yaml")
		content := "image: ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
		
		// Set environment variables
		os.Setenv("REGISTRY", "docker.io")
		os.Setenv("IMAGE_NAME", "myapp")
		os.Setenv("IMAGE_TAG", "v1.0.0")
		defer os.Unsetenv("REGISTRY")
		defer os.Unsetenv("IMAGE_NAME")
		defer os.Unsetenv("IMAGE_TAG")

		// Write test file
		err := os.WriteFile(testFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		// Run expandEnv
		expandEnv(testFile)

		// Read result
		result, err := os.ReadFile(testFile)
		if err != nil {
			t.Fatalf("Failed to read result file: %v", err)
		}

		expected := "image: docker.io/myapp:v1.0.0"
		if string(result) != expected {
			t.Errorf("Expected %q, got %q", expected, string(result))
		}
	})
}

// TestRootCommandWithSingleFile tests the root command with a single file
func TestRootCommandWithSingleFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "helm-subenv-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "values.yaml")
	content := "image: $TEST_IMAGE"
	
	// Set environment variable
	os.Setenv("TEST_IMAGE", "nginx:latest")
	defer os.Unsetenv("TEST_IMAGE")

	// Write test file
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Set paths and execute
	paths = []string{testFile}
	rootCmd.Run(rootCmd, []string{})

	// Read result
	result, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read result file: %v", err)
	}

	expected := "image: nginx:latest"
	if string(result) != expected {
		t.Errorf("Expected %q, got %q", expected, string(result))
	}
}

// TestRootCommandWithDirectory tests the root command with a directory
func TestRootCommandWithDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "helm-subenv-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create multiple test files
	file1 := filepath.Join(tmpDir, "values1.yaml")
	file2 := filepath.Join(tmpDir, "values2.yaml")
	
	content1 := "image: $TEST_IMAGE1"
	content2 := "image: $TEST_IMAGE2"
	
	// Set environment variables
	os.Setenv("TEST_IMAGE1", "nginx:1.0")
	os.Setenv("TEST_IMAGE2", "redis:2.0")
	defer os.Unsetenv("TEST_IMAGE1")
	defer os.Unsetenv("TEST_IMAGE2")

	// Write test files
	err = os.WriteFile(file1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file 1: %v", err)
	}
	err = os.WriteFile(file2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file 2: %v", err)
	}

	// Set paths and execute
	paths = []string{tmpDir}
	rootCmd.Run(rootCmd, []string{})

	// Read results
	result1, err := os.ReadFile(file1)
	if err != nil {
		t.Fatalf("Failed to read result file 1: %v", err)
	}
	result2, err := os.ReadFile(file2)
	if err != nil {
		t.Fatalf("Failed to read result file 2: %v", err)
	}

	expected1 := "image: nginx:1.0"
	expected2 := "image: redis:2.0"
	
	if string(result1) != expected1 {
		t.Errorf("File 1: Expected %q, got %q", expected1, string(result1))
	}
	if string(result2) != expected2 {
		t.Errorf("File 2: Expected %q, got %q", expected2, string(result2))
	}
}

// TestRootCommandWithMultipleFiles tests the root command with multiple files
func TestRootCommandWithMultipleFiles(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "helm-subenv-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	file1 := filepath.Join(tmpDir, "file1.yaml")
	file2 := filepath.Join(tmpDir, "file2.yaml")
	
	content1 := "service: $SERVICE_NAME"
	content2 := "port: $SERVICE_PORT"
	
	// Set environment variables
	os.Setenv("SERVICE_NAME", "api")
	os.Setenv("SERVICE_PORT", "8080")
	defer os.Unsetenv("SERVICE_NAME")
	defer os.Unsetenv("SERVICE_PORT")

	// Write test files
	err = os.WriteFile(file1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file 1: %v", err)
	}
	err = os.WriteFile(file2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file 2: %v", err)
	}

	// Set paths and execute
	paths = []string{file1, file2}
	rootCmd.Run(rootCmd, []string{})

	// Read results
	result1, err := os.ReadFile(file1)
	if err != nil {
		t.Fatalf("Failed to read result file 1: %v", err)
	}
	result2, err := os.ReadFile(file2)
	if err != nil {
		t.Fatalf("Failed to read result file 2: %v", err)
	}

	expected1 := "service: api"
	expected2 := "port: 8080"
	
	if string(result1) != expected1 {
		t.Errorf("File 1: Expected %q, got %q", expected1, string(result1))
	}
	if string(result2) != expected2 {
		t.Errorf("File 2: Expected %q, got %q", expected2, string(result2))
	}
}

// TestRootCommandWithNestedDirectory tests the root command with nested directories
func TestRootCommandWithNestedDirectory(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir, err := os.MkdirTemp("", "helm-subenv-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create nested directory
	nestedDir := filepath.Join(tmpDir, "nested")
	err = os.Mkdir(nestedDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested directory: %v", err)
	}

	file1 := filepath.Join(tmpDir, "root.yaml")
	file2 := filepath.Join(nestedDir, "nested.yaml")
	
	content := "env: $TEST_ENV"
	
	// Set environment variable
	os.Setenv("TEST_ENV", "production")
	defer os.Unsetenv("TEST_ENV")

	// Write test files
	err = os.WriteFile(file1, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file 1: %v", err)
	}
	err = os.WriteFile(file2, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file 2: %v", err)
	}

	// Set paths and execute
	paths = []string{tmpDir}
	rootCmd.Run(rootCmd, []string{})

	// Read results
	result1, err := os.ReadFile(file1)
	if err != nil {
		t.Fatalf("Failed to read result file 1: %v", err)
	}
	result2, err := os.ReadFile(file2)
	if err != nil {
		t.Fatalf("Failed to read result file 2: %v", err)
	}

	expected := "env: production"
	
	if string(result1) != expected {
		t.Errorf("Root file: Expected %q, got %q", expected, string(result1))
	}
	if string(result2) != expected {
		t.Errorf("Nested file: Expected %q, got %q", expected, string(result2))
	}
}
