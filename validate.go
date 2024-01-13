package rbx_manifest

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
)

const (
	ValidationThreads = 5
)

var (
	ErrInvalidChecksum = errors.New("Checksum was invalid")
)

func validationThread(channel chan error, directory string, files []RobloxFile) {
	defer close(channel)

	for _, file := range files {
		select {
		case _, ok := <-channel:
			if !ok {
				return
			}
		default:
			if err := file.Validate(directory); err != nil {
				channel <- err
			}
		}
	}
}

func (F *RobloxFile) Validate(directory string) error {
	Location := filepath.Join(directory, F.Path)
	if _, err := os.Stat(Location); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	NewHash := md5.New()
	fileStream, err := os.Open(Location)
	if err != nil {
		return err
	}

	if _, err := io.Copy(NewHash, fileStream); err != nil {
		return err
	} else if hex.EncodeToString(NewHash.Sum(nil)) == string(F.Checksum) {
		return nil
	}

	return errors.Join(ErrInvalidChecksum, errors.New(""))
}

// https://stackoverflow.com/questions/35179656/slice-chunking-in-go
func chunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

func closeErrors(errors []chan error) {
	for _, err := range errors {
		close(err)
	}
}

func (M *RobloxManifest) Validate(directory string) error {
	chunks := chunkBy(M.Files, len(M.Files)/ValidationThreads)
	errors := make([]chan error, len(chunks))

	defer closeErrors(errors)

	for index, chunk := range chunks {
		go validationThread(errors[index], directory, chunk)
	}

	closedChannels := 0
	for {
		if closedChannels == len(errors) {
			return nil
		}
		for _, err_chan := range errors {
			select {
			case err, ok := <-err_chan:
				if ok {
					return err
				} else {
					closedChannels++
				}
			}
		}
	}
}
