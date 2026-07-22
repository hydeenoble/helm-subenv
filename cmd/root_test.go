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
		if err := os.Setenv("TEST_IMAGE", "nginx"); err != nil {
			t.Fatalf("Failed to set TEST_IMAGE: %v", err)
		}
		if err := os.Setenv("TEST_TAG", "latest"); err != nil {
			t.Fatalf("Failed to set TEST_TAG: %v", err)
		}
		defer func() {
			_ = os.Unsetenv("TEST_IMAGE")
			_ = os.Unsetenv("TEST_TAG")
		}()

		// Write test file
		err := os.WriteFile(testFile, []byte(content), 0o600)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		// Run expandEnv
		if err := expandEnv(testFile); err != nil {
			t.Fatalf("expandEnv failed: %v", err)
		}

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
		_ = os.Unsetenv("NONEXISTENT_VAR")

		// Write test file
		err := os.WriteFile(testFile, []byte(content), 0o600)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		// Run expandEnv
		if err := expandEnv(testFile); err != nil {
			t.Fatalf("expandEnv failed: %v", err)
		}

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
		if err := os.Setenv("PROTOCOL", "https"); err != nil {
			t.Fatalf("Failed to set PROTOCOL: %v", err)
		}
		if err := os.Setenv("HOST", "example.com"); err != nil {
			t.Fatalf("Failed to set HOST: %v", err)
		}
		if err := os.Setenv("PORT", "8080"); err != nil {
			t.Fatalf("Failed to set PORT: %v", err)
		}
		defer func() {
			_ = os.Unsetenv("PROTOCOL")
			_ = os.Unsetenv("HOST")
			_ = os.Unsetenv("PORT")
		}()

		// Write test file
		if err := os.WriteFile(testFile, []byte(content), 0o600); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		// Run expandEnv
		if err := expandEnv(testFile); err != nil {
			t.Fatalf("expandEnv failed: %v", err)
		}

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
		if err := os.Setenv("REGISTRY", "docker.io"); err != nil {
			t.Fatalf("Failed to set REGISTRY: %v", err)
		}
		if err := os.Setenv("IMAGE_NAME", "myapp"); err != nil {
			t.Fatalf("Failed to set IMAGE_NAME: %v", err)
		}
		if err := os.Setenv("IMAGE_TAG", "v1.0.0"); err != nil {
			t.Fatalf("Failed to set IMAGE_TAG: %v", err)
		}
		defer func() {
			_ = os.Unsetenv("REGISTRY")
			_ = os.Unsetenv("IMAGE_NAME")
			_ = os.Unsetenv("IMAGE_TAG")
		}()

		// Write test file
		if err := os.WriteFile(testFile, []byte(content), 0o600); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		// Run expandEnv
		if err := expandEnv(testFile); err != nil {
			t.Fatalf("expandEnv failed: %v", err)
		}

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

// TestProcessPath tests the processPath function with a single file
func TestProcessPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "helm-subenv-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.yaml")
	content := "value: $TEST_VAR"

	if err := os.Setenv("TEST_VAR", "success"); err != nil {
		t.Fatalf("Failed to set TEST_VAR: %v", err)
	}
	defer func() {
		_ = os.Unsetenv("TEST_VAR")
	}()

	if err := os.WriteFile(testFile, []byte(content), 0o600); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	if err := processPath(testFile); err != nil {
		t.Fatalf("processPath failed: %v", err)
	}

	result, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read result file: %v", err)
	}

	expected := "value: success"
	if string(result) != expected {
		t.Errorf("Expected %q, got %q", expected, string(result))
	}
}

// TestProcessDirectory tests the processDirectory function
func TestProcessDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "helm-subenv-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	file1 := filepath.Join(tmpDir, "file1.yaml")
	file2 := filepath.Join(tmpDir, "file2.yaml")

	if err := os.Setenv("VAR1", "value1"); err != nil {
		t.Fatalf("Failed to set VAR1: %v", err)
	}
	if err := os.Setenv("VAR2", "value2"); err != nil {
		t.Fatalf("Failed to set VAR2: %v", err)
	}
	defer func() {
		_ = os.Unsetenv("VAR1")
		_ = os.Unsetenv("VAR2")
	}()

	if err := os.WriteFile(file1, []byte("key: $VAR1"), 0o600); err != nil {
		t.Fatalf("Failed to write file1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("key: $VAR2"), 0o600); err != nil {
		t.Fatalf("Failed to write file2: %v", err)
	}

	if err := processDirectory(tmpDir); err != nil {
		t.Fatalf("processDirectory failed: %v", err)
	}

	result1, err := os.ReadFile(file1)
	if err != nil {
		t.Fatalf("Failed to read file1: %v", err)
	}
	result2, err := os.ReadFile(file2)
	if err != nil {
		t.Fatalf("Failed to read file2: %v", err)
	}

	if string(result1) != "key: value1" {
		t.Errorf("File1: Expected %q, got %q", "key: value1", string(result1))
	}
	if string(result2) != "key: value2" {
		t.Errorf("File2: Expected %q, got %q", "key: value2", string(result2))
	}
}
