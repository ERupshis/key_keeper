// Package config parses config settings from flags and environments.
// Environments have higher priority than flags.
package config

import (
	"errors"
	"flag"

	"github.com/caarlos0/env"
	"github.com/erupshis/key_keeper/internal/common/utils/configutils"
)

// Config agent's settings.
type Config struct {
	Host        string
	DatabaseDSN string
	JWT         string
	HashKey     string
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
	flagServerHost  = "addr"
	flagDatabaseDSN = "rdsn"
	flagJWT         = "jwt"
	flagHashKey     = "hk"
)

// checkFlags checks flags of app's launch.
func checkFlags(config *Config) {
	flag.StringVar(&config.Host, flagServerHost, ":8081", "server host")
	flag.StringVar(&config.DatabaseDSN, flagDatabaseDSN, "postgres://postgres:postgres@localhost:5432/key_db?sslmode=disable", "records database DSN")
	flag.StringVar(&config.JWT, flagJWT, "SECRET_KEY", "jwt token generation key")
	flag.StringVar(&config.HashKey, flagHashKey, "SECRET_KEY", "user passwords hasher key")

	flag.Parse()
}

// ENVIRONMENTS PARSING.
// envConfig struct of environments suitable for agent.
type envConfig struct {
	Host        string `env:"HOST"`
	DatabaseDSN string `env:"DATABASE_DSN"`
	JWT         string `env:"JWT_KEY"`
	HashKey     string `env:"HASH_KEY"`
}

// checkEnvironments checks environments suitable for agent.
func checkEnvironments(config *Config) error {
	var envs = envConfig{}
	err := env.Parse(&envs)
	if err != nil {
		return configutils.ErrCheckEnvsWrapper(err)
	}

	var errs []error
	errs = append(errs, configutils.SetEnvToParamIfNeed(&config.Host, envs.Host))
	errs = append(errs, configutils.SetEnvToParamIfNeed(&config.DatabaseDSN, envs.DatabaseDSN))
	errs = append(errs, configutils.SetEnvToParamIfNeed(&config.JWT, envs.JWT))
	errs = append(errs, configutils.SetEnvToParamIfNeed(&config.HashKey, envs.HashKey))

	resErr := errors.Join(errs...)
	if resErr != nil {
		return configutils.ErrCheckEnvsWrapper(resErr)
	}

	return nil
}
