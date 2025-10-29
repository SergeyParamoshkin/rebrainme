package testdata

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateGolden = flag.Bool("update", false, "update golden files")

func TestParseConfig_SingleFile(t *testing.T) {
	data, err := os.ReadFile("testdata/valid_config.json")
	require.NoError(t, err, "failed to read test data")

	config, err := ParseConfig(data)
	require.NoError(t, err)

	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, 30, config.Timeout)
	assert.True(t, config.Debug)
}

func TestParseConfig_MultipleFiles(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantHost string
		wantPort int
		wantErr  bool
	}{
		{
			name:     "valid config",
			filename: "testdata/valid_config.json",
			wantHost: "localhost",
			wantPort: 8080,
			wantErr:  false,
		},
		{
			name:     "production config",
			filename: "testdata/production_config.json",
			wantHost: "example.com",
			wantPort: 443,
			wantErr:  false,
		},
		{
			name:     "invalid config - empty host",
			filename: "testdata/invalid_config.json",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.filename)
			require.NoError(t, err, "failed to read %s", tt.filename)

			config, err := ParseConfig(data)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantHost, config.Host)
			assert.Equal(t, tt.wantPort, config.Port)
		})
	}
}

func TestParseUsers_WithFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/fixtures/users.json")
	require.NoError(t, err)

	users, err := ParseUsers(data)
	require.NoError(t, err)

	assert.Len(t, users, 3)

	assert.Equal(t, 1, users[0].ID)
	assert.Equal(t, "Alice", users[0].Name)
	assert.Equal(t, "alice@example.com", users[0].Email)

	assert.Equal(t, 2, users[1].ID)
	assert.Equal(t, "Bob", users[1].Name)

	assert.Equal(t, 3, users[2].ID)
	assert.Equal(t, "Charlie", users[2].Name)
}

func loadUserFixtures(t *testing.T) []User {
	t.Helper()

	data, err := os.ReadFile("testdata/fixtures/users.json")
	require.NoError(t, err)

	users, err := ParseUsers(data)
	require.NoError(t, err)

	return users
}

func TestUsersValidation_WithHelper(t *testing.T) {
	users := loadUserFixtures(t)

	for _, user := range users {
		t.Run(user.Name, func(t *testing.T) {
			assert.NotEmpty(t, user.Name)
			assert.NotEmpty(t, user.Email)
			assert.Greater(t, user.Age, 0)
		})
	}
}

func TestGenerateReport_GoldenFile(t *testing.T) {
	tests := []struct {
		name   string
		report *Report
		golden string
	}{
		{
			name: "simple report",
			report: &Report{
				Title:   "Monthly Report",
				Date:    "2024-01-15",
				Summary: "This is a summary of monthly activities.",
				Items: []string{
					"Task A completed",
					"Task B in progress",
					"Task C pending",
				},
			},
			golden: "testdata/golden/simple_report.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateReport(tt.report)

			if *updateGolden {
				dir := filepath.Dir(tt.golden)
				err := os.MkdirAll(dir, 0o755)
				require.NoError(t, err)

				err = os.WriteFile(tt.golden, []byte(result), 0o644)
				require.NoError(t, err)
				t.Logf("Updated golden file: %s", tt.golden)
				return
			}

			expected, err := os.ReadFile(tt.golden)
			require.NoError(t, err, "failed to read golden file %s", tt.golden)

			assert.Equal(t, string(expected), result)
		})
	}
}

func TestGenerateReport_MultipleGoldenFiles(t *testing.T) {
	tests := []struct {
		name   string
		report *Report
		golden string
	}{
		{
			name: "empty report",
			report: &Report{
				Title:   "Empty Report",
				Date:    "2024-01-01",
				Summary: "Nothing to report.",
				Items:   []string{},
			},
			golden: "testdata/golden/empty_report.txt",
		},
		{
			name: "report with many items",
			report: &Report{
				Title:   "Quarterly Review",
				Date:    "2024-03-31",
				Summary: "Q1 2024 achievements.",
				Items: []string{
					"Launched new feature A",
					"Improved performance by 50%",
					"Fixed 100+ bugs",
					"Added comprehensive tests",
					"Updated documentation",
				},
			},
			golden: "testdata/golden/quarterly_report.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateReport(tt.report)

			if *updateGolden {
				dir := filepath.Dir(tt.golden)
				err := os.MkdirAll(dir, 0o755)
				require.NoError(t, err)

				err = os.WriteFile(tt.golden, []byte(result), 0o644)
				require.NoError(t, err)
				t.Logf("Updated golden file: %s", tt.golden)
				return
			}

			if _, err := os.Stat(tt.golden); os.IsNotExist(err) {
				t.Skipf("Golden file %s does not exist. Run with -update to create it.", tt.golden)
			}

			expected, err := os.ReadFile(tt.golden)
			require.NoError(t, err, "failed to read golden file %s", tt.golden)

			assert.Equal(t, string(expected), result)
		})
	}
}

func TestParseConfig_AllFilesInDirectory(t *testing.T) {
	files, err := filepath.Glob("testdata/*_config.json")
	require.NoError(t, err)
	require.NotEmpty(t, files, "no config files found")

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			data, err := os.ReadFile(file)
			require.NoError(t, err)

			_, err = ParseConfig(data)

			if filepath.Base(file) == "invalid_config.json" {
				assert.Error(t, err, "expected error for invalid config")
			}
		})
	}
}

func Example() {
	data, _ := os.ReadFile("testdata/valid_config.json")
	config, _ := ParseConfig(data)

	println("Host:", config.Host)
	println("Port:", config.Port)
}

func ExampleGenerateReport() {
	report := &Report{
		Title:   "Test Report",
		Date:    "2024-01-01",
		Summary: "Summary text",
		Items:   []string{"Item 1", "Item 2"},
	}

	result := GenerateReport(report)
	println(result)
}
