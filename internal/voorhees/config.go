package voorhees

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Nivl/voorhees/internal/errutil"
	"gopkg.in/yaml.v2"
)

const (
	week  = 7 * 24 * time.Hour
	month = 30 * 24 * time.Hour
)

// ConfigFile represent the configuration file
type ConfigFile struct {
	Version int `yaml:"version"`
	Default struct {
		Limit string `yaml:"limit"`
	} `yaml:"default"`
	Rules map[string]string `yaml:"rules"`
}

// Config represents the configuration of
type Config struct {
	toSkip       map[string]struct{}
	limits       map[string]time.Duration
	defaultLimit time.Duration
}

// IsIgnored check if a package should be ignored
func (cfg *Config) IsIgnored(pkg string) bool {
	_, ok := cfg.toSkip[pkg]
	return ok
}

// Duration returns the duration of a package
func (cfg *Config) Duration(pkg string) time.Duration {
	if limit, ok := cfg.limits[pkg]; ok {
		return limit
	}
	return cfg.defaultLimit
}

// NewConfig return a new Config from a a reader
func NewConfig(r io.Reader) (*Config, error) {
	cf := &ConfigFile{
		Version: 1,
		Rules:   map[string]string{},
	}
	dec := yaml.NewDecoder(r)
	dec.SetStrict(true) // Strict mode prevents duplicates and unknown fields
	if err := dec.Decode(cf); err != nil {
		return nil, fmt.Errorf("could not parse config: %w", err)
	}

	cfg := newDefaultConfig()
	switch cf.Version {
	case 1:
		if cf.Default.Limit != "" {
			limit, err := parseConfigDuration(cf.Default.Limit)
			if err != nil {
				return nil, fmt.Errorf("could not parse config: default.limit %w", err)
			}
			cfg.defaultLimit = limit
		}

		for pkg, v := range cf.Rules {
			// We want to keep it lowercase in case the user doesn't
			// enter the exact same value as what is in the gomod
			pkg = strings.ToLower(pkg)

			value := strings.ToLower(v)
			switch value {
			case "ignore", "skip":
				cfg.toSkip[pkg] = struct{}{}
			default:
				limit, err := parseConfigDuration(value)
				if err != nil {
					return nil, fmt.Errorf("could not parse config: package %s: %w", pkg, err)
				}
				cfg.limits[pkg] = limit
			}
		}
		return cfg, nil
	default:
		return nil, fmt.Errorf("unsupported config version: %d", cf.Version)
	}
}

func newDefaultConfig() *Config {
	return &Config{
		toSkip:       map[string]struct{}{},
		limits:       map[string]time.Duration{},
		defaultLimit: 6 * month,
	}
}

// parseConfigDuration parses a duration such as: 6 weeks
func parseConfigDuration(line string) (time.Duration, error) {
	duration := strings.Split(line, " ")
	if len(duration) != 2 {
		return 0, fmt.Errorf("invalid rule value: %s", line)
	}
	n, err := strconv.Atoi(duration[0])
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("expected a number > 0: %s", duration[0])
	}
	switch duration[1] {
	case "week", "weeks":
		return time.Duration(n) * week, nil
	case "month", "months":
		return time.Duration(n) * month, nil
	default:
		return 0, fmt.Errorf("unexpected duration type: %s", duration[1])
	}
}

// LoadConfigFile load the configuration file located at the given path
func LoadConfigFile(path string) (cfg *Config, err error) {
	f, err := os.Open(path)
	if err != nil {
		// if the config file doesn't exist and the path points to the
		// default path, then we assume the user doesn't use a config
		// file
		if errors.Is(err, os.ErrNotExist) && path == "./.voorhees.yml" {
			return newDefaultConfig(), nil
		}
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer errutil.Close(f, &err)
	return NewConfig(f)
}
