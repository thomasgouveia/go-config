package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type format string

const (
	JSON format = "json"
	YAML format = "yaml"
)

var (
	ErrInvalidConfigFormat = errors.New("invalid config format. valid values are: json, yaml")
	ErrMissingDefaultValue = errors.New("default value is nil. you must provide a valid default value for your configuration struct")
)

type configLoader[T interface{}] struct {
	opts *Options[T]
	vip  *viper.Viper
}

type Options[T interface{}] struct {
	// Format for the configuration file. Must be one of the following : json, `yaml`
	// If not set, will default to `yaml`.
	Format format
	// Default is the default value for the configuration.
	Default *T
	// FileName is the name of the configuration file without the extension, e.g: config
	FileName string
	// FileLocations is a list of paths on disk where the loader should look
	// for configuration files matching the pattern : Filename.Format
	//
	// For example, if you define `Format: "json" and `FileName: "config"`, the loader load every `config.json` file
	// in those directories.
	FileLocations []string
	// EnvEnabled is a boolean that indicates whether the configuration must check for
	// environment variables or not during the loading process.
	EnvEnabled bool
	// EnvPrefix is an optional string used as prefix for all environment variables.
	// For example, if specify `EnvPrefix: "APP`, all the environments variables for your configuration must
	// match the following pattern : `APP_*`
	EnvPrefix string
}

// NewLoader create a new configuration loader for the given type T.
func NewLoader[T interface{}](opts *Options[T]) (*configLoader[T], error) {
	v := viper.New()

	// If we don't have a default value,
	// the loader will not work as expected as
	// we want to merge the default configuration with
	// overrides that come from the different sources.
	if opts.Default == nil {
		return nil, ErrMissingDefaultValue
	}

	// Check that the configuration format is valid
	// By default if no option was provided into the options,
	// the default format used is YAML
	if opts.Format == "" {
		opts.Format = YAML
	}

	if !isValidFormat(opts.Format) {
		return nil, ErrInvalidConfigFormat
	}

	v.SetConfigType(string(opts.Format))

	if opts.FileName == "" {
		opts.FileName = "config"
	}
	v.SetConfigName(opts.FileName)

	// Register locations for the configuration file
	for _, path := range opts.FileLocations {
		v.AddConfigPath(path)
	}

	if opts.EnvEnabled {
		// If enabled by the user, we should check for environment
		// variables that match the pattern `EnvPrefix_*` if EnvPrefix isn't empty
		if opts.EnvPrefix != "" {
			v.SetEnvPrefix(strings.ToUpper(opts.EnvPrefix))
		}
		v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		v.AutomaticEnv()
	}

	return &configLoader[T]{opts, v}, nil
}

// Load will try to retrieve the configuration from the different sources and return the merged result.
func (cl *configLoader[T]) Load() (*T, error) {
	// Set the default values into Viper.
	// Viper needs to know if a key exists in order to
	// be able to override it.
	b, err := cl.marshal(cl.opts.Default)
	if err != nil {
		return nil, err
	}

	re := bytes.NewReader(b)
	if err := cl.vip.MergeConfig(re); err != nil {
		return nil, err
	}

	// Parse configuration for all additional sources
	if err := cl.vip.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// Finally, we can unmarshal the viper loaded configuration into our struct
	config := cl.opts.Default
	if err := cl.vip.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

// marshal is a simple wrapper function around the different implementations of `Marshal`, as we depend on the
// configuration format provided in loader options.
func (cl *configLoader[T]) marshal(v any) ([]byte, error) {
	switch cl.opts.Format {
	case JSON:
		return json.Marshal(v)
	case YAML:
		return yaml.Marshal(v)
	default:
		return nil, ErrInvalidConfigFormat
	}
}
