// Package config parses config settings from flags and environments.
// Environments have higher priority than flags.
package config

import (
	"errors"
	"flag"
	"time"

	"github.com/caarlos0/env"
	"github.com/erupshis/key_keeper/internal/common/utils/configutils"
)

// Config agent's settings.
type Config struct {
	LocalStoragePath   string
	LocalStoreInterval time.Duration
}

// Parse main func to parse variables.
func Parse() (Config, error) {
	var config = Config{}
	checkFlags(&config)
	err := checkEnvironments(&config)
	return config, err
}

// FLAGS PARSING.
const (
	flagLocalStoragePath   = "lsp"
	flagLocalStoreInterval = "lsi"
)

// checkFlags checks flags of app's launch.
func checkFlags(config *Config) {
	flag.StringVar(&config.LocalStoragePath, flagLocalStoragePath, "C:/data/local_storage.json", "local storage path")
	flag.DurationVar(&config.LocalStoreInterval, flagLocalStoreInterval, 10*time.Second, "local store interval. 0 - means store on data change")

	flag.Parse()
}

// ENVIRONMENTS PARSING.
// envConfig struct of environments suitable for agent.
type envConfig struct {
	LocalStoragePath   string `env:"LOCAL_STORAGE_PATH"`
	LocalStoreInterval string `env:"LOCAL_STORE_INTERVAL"`
}

// checkEnvironments checks environments suitable for agent.
func checkEnvironments(config *Config) error {
	var envs = envConfig{}
	err := env.Parse(&envs)
	if err != nil {
		return configutils.ErrCheckEnvsWrapper(err)
	}

	var errs []error
	errs = append(errs, configutils.SetEnvToParamIfNeed(&config.LocalStoragePath, envs.LocalStoragePath))
	errs = append(errs, configutils.SetEnvToParamIfNeed(&config.LocalStoreInterval, envs.LocalStoreInterval))

	resErr := errors.Join(errs...)
	if resErr != nil {
		return configutils.ErrCheckEnvsWrapper(resErr)
	}

	return nil
}
