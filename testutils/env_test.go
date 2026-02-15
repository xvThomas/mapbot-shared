package testutils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetupTestEnv(t *testing.T) {
	// Test that SetupTestEnv does not panic and handles missing .env.test gracefully
	SetupTestEnv(t)

	// Test with existing environment variable - should not be overwritten
	testVar := "TEST_SETUP_ENV_VAR"
	expectedValue := "existing_value"
	os.Setenv(testVar, expectedValue)
	defer os.Unsetenv(testVar)

	// Create .env.test with different value for same variable
	projectRoot := getProjectRoot()
	envFile := filepath.Join(projectRoot, ".env.test")
	testContent := testVar + "=new_value\n"
	
	if err := os.WriteFile(envFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test .env.test: %v", err)
	}
	defer os.Remove(envFile)

	// Call SetupTestEnv - should not overwrite existing env var
	SetupTestEnv(t)

	got := os.Getenv(testVar)
	if got != expectedValue {
		t.Errorf("SetupTestEnv() overwrote existing env var: got %q, want %q", got, expectedValue)
	}
}

func TestSetupTestEnvWithRequiredVarsOrSkipTest(t *testing.T) {
	t.Run("skips test when required variable is missing", func(t *testing.T) {
		mockT := &mockTestingT{TB: t}
		os.Unsetenv("REQUIRED_VAR_THAT_DOES_NOT_EXIST")

		SetupTestEnvWithRequiredVarsOrSkipTest(mockT, "REQUIRED_VAR_THAT_DOES_NOT_EXIST")

		if !mockT.skipped {
			t.Error("SetupTestEnvWithRequiredVarsOrSkipTest() should have skipped test")
		}
	})

	t.Run("continues when all required variables are set", func(t *testing.T) {
		mockT := &mockTestingT{TB: t}
		os.Setenv("REQUIRED_VAR_TEST", "some_value")
		defer os.Unsetenv("REQUIRED_VAR_TEST")

		SetupTestEnvWithRequiredVarsOrSkipTest(mockT, "REQUIRED_VAR_TEST")

		if mockT.skipped {
			t.Errorf("SetupTestEnvWithRequiredVarsOrSkipTest() should not skip, message: %s", mockT.skipMessage)
		}
	})

	t.Run("checks multiple required variables", func(t *testing.T) {
		mockT := &mockTestingT{TB: t}
		os.Setenv("VAR1", "value1")
		defer os.Unsetenv("VAR1")
		os.Unsetenv("VAR2")

		SetupTestEnvWithRequiredVarsOrSkipTest(mockT, "VAR1", "VAR2")

		if !mockT.skipped {
			t.Error("SetupTestEnvWithRequiredVarsOrSkipTest() should skip when any required var is missing")
		}
	})
}

func TestGetProjectRoot(t *testing.T) {
	root := getProjectRoot()

	goModPath := filepath.Join(root, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Errorf("getProjectRoot() returned %q, but go.mod does not exist there", root)
	}

	if !filepath.IsAbs(root) {
		t.Errorf("getProjectRoot() should return absolute path, got %q", root)
	}
}

type mockTestingT struct {
	testing.TB
	skipped     bool
	skipMessage string
}

func (m *mockTestingT) Skipf(format string, args ...interface{}) {
	m.skipped = true
	m.skipMessage = format
}

func (m *mockTestingT) Helper() {}
