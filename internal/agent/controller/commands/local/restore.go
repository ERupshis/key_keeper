package local

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/erupshis/key_keeper/internal/agent/errs"
	"github.com/erupshis/key_keeper/internal/agent/storage/inmemory"
	"github.com/erupshis/key_keeper/internal/agent/storage/local"
	"github.com/erupshis/key_keeper/internal/agent/utils"
)

func (l *Local) ProcessRestore(ctx context.Context, exist bool, inmemory *inmemory.Storage, localStorage *local.FileManager) error {
	return l.handleRestore(ctx, exist, inmemory, localStorage)
}

type restoreState int

const (
	restoreInitialState   = restoreState(0)
	restorePassPhrase     = restoreState(1)
	restoreStorageDecode  = restoreState(2)
	restoreNewStoragePath = restoreState(3)
	restoreNewPassPhrase  = restoreState(4)
	restoreFinishState    = restoreState(5)
)

func (l *Local) handleRestore(ctx context.Context, exist bool, inmemory *inmemory.Storage, localStorage *local.FileManager) error {
	currentState := restoreInitialState
	if !exist {
		currentState = restoreNewPassPhrase
	}

	var err error
	var passPhrase string
	for currentState != restoreFinishState {
		switch currentState {
		case restoreInitialState:
			currentState = l.stateInitial(localStorage)
		case restorePassPhrase:
			{
				currentState, passPhrase, err = l.statePassPhrase()
				if err != nil {
					return err
				}
			}
		case restoreStorageDecode:
			{
				currentState, err = l.stateStorageDecode(ctx, inmemory, localStorage, passPhrase)
				if err != nil {
					return err
				}
			}
		case restoreNewStoragePath:
			{
				currentState, err = l.stateNewStorage(localStorage)
				if err != nil {
					return err
				}
			}
		case restoreNewPassPhrase:
			{
				currentState, err = l.stateNewPassPhrase(ctx, localStorage)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (l *Local) stateInitial(localStorage *local.FileManager) restoreState {
	l.iactr.Printf("enter passphrase to decode local storage (%s): ", localStorage.Path())
	return restorePassPhrase
}

func (l *Local) statePassPhrase() (restoreState, string, error) {
	passPhrase, ok, err := l.iactr.GetUserInputAndValidate(nil)

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		return restoreNewStoragePath, "", nil
	}

	return restoreStorageDecode, passPhrase, err
}

func (l *Local) stateStorageDecode(ctx context.Context, inmemory *inmemory.Storage, localStorage *local.FileManager, passPhrase string) (restoreState, error) {
	localStorage.SetPassPhrase(passPhrase)
	records, err := localStorage.RestoreUserData(ctx)
	if err != nil {
		l.iactr.Printf("failed to decode storage, reenter passphrase or '%s' to create new storage:\n", utils.CommandCancel)
		return restorePassPhrase, nil
	}

	if err = inmemory.AddRecords(records); err != nil {
		return restoreStorageDecode, fmt.Errorf("write local storage data in memory: %w", err)
	}

	return restoreFinishState, nil
}

func (l *Local) stateNewStorage(localStorage *local.FileManager) (restoreState, error) {
	l.iactr.Printf("enter new storage path(current: '%s'): ", localStorage.Path()[:strings.LastIndex(localStorage.Path(), "/")])
	newPath, ok, err := l.iactr.GetUserInputAndValidate(nil)
	if newPath == localStorage.Path() {
		l.iactr.Printf("attempt to rewrite existing storage, try again: ")
		return restoreNewStoragePath, nil
	}

	if !ok {
		return restoreNewStoragePath, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		l.iactr.Printf("cancel is not allowed here\n")
		return restoreNewStoragePath, nil
	}

	if !strings.HasSuffix(newPath, "/") {
		newPath = newPath + "/"
	}

	localStorage.SetPath(newPath)
	return restoreNewPassPhrase, nil
}

func (l *Local) stateNewPassPhrase(ctx context.Context, localStorage *local.FileManager) (restoreState, error) {
	l.iactr.Printf("enter passphrase for local storage securing: ")

	newPassPhrase, ok, err := l.iactr.GetUserInputAndValidate(nil)
	if !ok {
		return restoreNewStoragePath, err
	}

	if ok && errors.Is(err, errs.ErrInterruptedByUser) {
		l.iactr.Printf("cancel is not allowed here\n")
		return restoreNewStoragePath, nil
	}

	localStorage.SetPassPhrase(newPassPhrase)
	localStorage.RunAutoSave(ctx)
	return restoreFinishState, nil
}
