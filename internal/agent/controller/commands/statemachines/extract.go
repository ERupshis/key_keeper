package statemachines

import (
	"errors"
	"path/filepath"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/common/data"
)

type stateExtract int

const (
	extractInitialState  = stateExtract(0)
	extractSavePathState = stateExtract(1)
	extractFinishState   = stateExtract(2)
)

type ExtractConfig struct {
	Record   *data.Record
	FileSave func(record *data.Record, savePath string) error
}

func (s *StateMachines) Extract(cfg ExtractConfig) error {
	currentState := extractInitialState

	var pathToFile string
	var err error
	for currentState != extractFinishState {
		switch currentState {
		case extractInitialState:
			{
				pathToFile, err = s.extractFilePath()
				if err != nil {
					if errors.Is(err, errs.ErrInterruptedByUser) {
						return err
					} else {
						continue
					}
				}
				currentState = extractSavePathState
			}
		case extractSavePathState:
			{
				if err = cfg.FileSave(cfg.Record, pathToFile); err != nil {
					if errors.Is(err, errs.ErrInterruptedByUser) {
						return err
					} else {
						continue
					}
				}
				currentState = extractFinishState
			}
		}
	}

	return nil
}

// MAIN DATA STATE MACHINE.
type extractPathState int

const (
	extractPathInitialState  = extractPathState(0)
	extractPathFilePathState = extractPathState(1)
	extractPathFinishState   = extractPathState(2)
)

func (s *StateMachines) extractFilePath() (string, error) {
	currentState := extractPathInitialState

	var pathToFile string
	var err error
	for currentState != extractPathFinishState {
		switch currentState {
		case extractPathInitialState:
			currentState = s.stateInitial()
		case extractPathFilePathState:
			{
				currentState, pathToFile, err = s.stateFilePath()
				if err != nil {
					return "", err
				}
			}
		}
	}

	return pathToFile, nil
}

func (s *StateMachines) stateInitial() extractPathState {
	s.iactr.Printf("enter absolute path to file: ")
	return extractPathFilePathState
}

func (s *StateMachines) stateFilePath() (extractPathState, string, error) {
	pathToFile, ok, err := s.iactr.GetUserInputAndValidate(nil)
	if !ok {
		return extractPathFilePathState, "", err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return extractPathFilePathState, "", err
	}

	if !filepath.IsAbs(pathToFile) {
		s.iactr.Printf("entered local path. Try to set absolute path: ")
		return extractPathFilePathState, "", nil
	}

	return extractPathFinishState, pathToFile, nil
}
