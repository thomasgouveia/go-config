package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockConfig struct {
	Foo string `yaml:"foo" json:"foo"`
	Bar string `yaml:"bar" json:"bar"`
	Baz string `yaml:"baz" json:"port"`
}

var defaultMockConfig = &mockConfig{
	Foo: "foo",
	Bar: "bar",
	Baz: "baz",
}

func Test_NewLoader_ErrorWhenInvalidFormat(t *testing.T) {
	t.Parallel()
	_, err := NewLoader(&Options[mockConfig]{
		Format:  "N/A",
		Default: defaultMockConfig,
	})
	assert.Error(t, err, "invalid config format. valid values are: json, yaml")
}

func Test_NewLoader_ErrorWhenNoDefault(t *testing.T) {
	t.Parallel()
	_, err := NewLoader(&Options[mockConfig]{
		Format: JSON,
	})
	assert.Error(t, err, "default value is nil. you must provide a valid default value for your configuration struct")
}

func Test_LoadConfig_ReturnDefaultIfNoOverrides(t *testing.T) {
	t.Parallel()
	cl, _ := NewLoader(&Options[mockConfig]{
		Format:  JSON,
		Default: defaultMockConfig,
	})

	cfg, err := cl.Load()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, defaultMockConfig, cfg)
}

func Test_LoadConfig_EnvironmentOverrides(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("MOCK_FOO")
	})
	opts := &Options[mockConfig]{
		Format:     JSON,
		EnvEnabled: true,
		EnvPrefix:  "mock",
		Default:    defaultMockConfig,
	}

	cl, err := NewLoader(opts)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("MOCK_FOO", "hello-test"); err != nil {
		t.Fatal(err)
	}

	cfg, err := cl.Load()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hello-test", cfg.Foo)
	assert.Equal(t, "bar", cfg.Bar)
	assert.Equal(t, "baz", cfg.Baz)
}

func Test_LoadConfig_FileAndEnvironmentOverrides(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("MOCK_BAR")
	})

	cfgJson := `{"foo": "hello-from-file"}`
	if err := os.WriteFile("/tmp/config.json", []byte(cfgJson), 0644); err != nil {
		t.Fatal(err)
	}

	opts := &Options[mockConfig]{
		Format:     JSON,
		EnvEnabled: true,
		EnvPrefix:  "mock",
		Default:    defaultMockConfig,
		FileLocations: []string{
			".",
		},
	}

	cl, err := NewLoader(opts)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("MOCK_BAR", "hello-from-env"); err != nil {
		t.Fatal(err)
	}

	cfg, err := cl.Load()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hello-from-file", cfg.Foo)
	assert.Equal(t, "hello-from-env", cfg.Bar)
	assert.Equal(t, "baz", cfg.Baz)
}

func Test_LoadConfig_EnvironmentTakesPrecedence(t *testing.T) {
	t.Parallel()
	t.Cleanup(func() {
		os.Unsetenv("MOCK_FOO")
	})

	cfgJson := `{"foo": "hello-from-file"}`
	if err := os.WriteFile("config.json", []byte(cfgJson), 0644); err != nil {
		t.Fatal(err)
	}

	opts := &Options[mockConfig]{
		Format:     JSON,
		EnvEnabled: true,
		EnvPrefix:  "mock",
		Default:    defaultMockConfig,
		FileLocations: []string{
			".",
		},
	}

	cl, err := NewLoader(opts)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("MOCK_FOO", "hello-from-env"); err != nil {
		t.Fatal(err)
	}

	cfg, err := cl.Load()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hello-from-env", cfg.Foo)
}
