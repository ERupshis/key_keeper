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
	ServerHost         string
	LocalStoragePath   string
	LocalStoreInterval time.Duration
	HashKey            string
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
	flagServerHost         = "addr"
	flagLocalStoragePath   = "lsp"
	flagLocalStoreInterval = "lsi"
	flagHashKey            = "h"
)

// checkFlags checks flags of app's launch.
func checkFlags(config *Config) {
	flag.StringVar(&config.ServerHost, flagServerHost, "http://localhost:8080", "server host")
	flag.StringVar(&config.LocalStoragePath, flagLocalStoragePath, "C:/key_keeper/data/", "folder for local storage")
	flag.DurationVar(&config.LocalStoreInterval, flagLocalStoreInterval, 10*time.Second, "local store interval. 0 - means store on models change")
	flag.StringVar(&config.HashKey, flagHashKey, "", "hash key for binary files hash sum calculation")

	flag.Parse()
}

// ENVIRONMENTS PARSING.
// envConfig struct of environments suitable for agent.
type envConfig struct {
	ServerHost         string `env:"SERVER_HOST"`
	LocalStoragePath   string `env:"LOCAL_STORAGE_PATH"`
	LocalStoreInterval string `env:"LOCAL_STORE_INTERVAL"`
	HashKey            string `env:"HASH_KEY"`
}

// checkEnvironments checks environments suitable for agent.
func checkEnvironments(config *Config) error {
	var envs = envConfig{}
	err := env.Parse(&envs)
	if err != nil {
		return configutils.ErrCheckEnvsWrapper(err)
	}

	var errs []error
	errs = append(errs, configutils.SetEnvToParamIfNeed(&config.ServerHost, envs.ServerHost))
	errs = append(errs, configutils.SetEnvToParamIfNeed(&config.LocalStoragePath, envs.LocalStoragePath))
	errs = append(errs, configutils.SetEnvToParamIfNeed(&config.LocalStoreInterval, envs.LocalStoreInterval))
	errs = append(errs, configutils.SetEnvToParamIfNeed(&config.HashKey, envs.HashKey))

	resErr := errors.Join(errs...)
	if resErr != nil {
		return configutils.ErrCheckEnvsWrapper(resErr)
	}

	return nil
}
