package binary

import (
	"github.com/erupshis/key_keeper/internal/agent/controller/commands/statemachines"
	"github.com/erupshis/key_keeper/internal/agent/interactor"
	"github.com/erupshis/key_keeper/internal/common/crypt/ska"
	"github.com/erupshis/key_keeper/internal/common/hasher"
)

type Config struct {
	Iactr     *interactor.Interactor
	Sm        *statemachines.StateMachines
	Hash      *hasher.Hasher
	Cryptor   *ska.SKA
	StorePath string
}

type Binary struct {
	iactr     *interactor.Interactor
	sm        *statemachines.StateMachines
	hash      *hasher.Hasher
	cryptor   *ska.SKA
	storePath string
}

func NewBinary(cfg *Config) *Binary {
	return &Binary{
		iactr:     cfg.Iactr,
		sm:        cfg.Sm,
		hash:      cfg.Hash,
		cryptor:   cfg.Cryptor,
		storePath: cfg.StorePath,
	}
}
