package rbx_manifest

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"
)

var (
	newLine = []byte("\r\n")
)

var (
	ErrTooShort         = errors.New("Provided Pkg-manifest is too short to parse")
	ErrBadSizedManifest = errors.New("The size of the provided manifest is not of the expected values")
	ErrBadSizedChunk    = errors.New("A chunk inside the manifest is too small")
	ErrBadSizedChecksum = errors.New("Checksum isnt of size 32")
)

type pkgSplit = [][4][]byte
type manifestSplit = [][2][]byte

func isEven(number int) bool {
	return (number | 1) > number
}

func splitPkgManifest(manifest []byte) (pkgSplit, error) {
	splitData := bytes.Split(manifest, newLine)
	if len(splitData) < 1 || !isEven(len(splitData)-1) {
		return nil, ErrBadSizedManifest
	}
	splitData = splitData[1:] /* This is the same as pop in other languages */

	foundData := make(pkgSplit, len(splitData)/4)
	for i := 0; i < len(splitData)/4; i++ {
		index := i * 4
		if len(splitData) < index+3 {
			return nil, ErrBadSizedChunk
		}
		foundData[i][0] = splitData[index]
		foundData[i][1] = splitData[index+1]
		foundData[i][2] = splitData[index+2]
		foundData[i][3] = splitData[index+3]
	}
	return foundData, nil
}

func deserializePkgManifest(data pkgSplit) (*RobloxPkgManifest, error) {
	if len(data) < 1 || !isEven(len(data)) {
		return nil, ErrBadSizedManifest
	}
	decodedFiles := make([]RobloxPkgFile, len(data))
	for i, chunk := range data {
		/* Decode data and remove any whitespace */
		name, hash, rZipSize, rSize := string(chunk[0]), string(chunk[1]), string(chunk[2]), string(chunk[3])

		/* Validate length before copying into a smaller buffer */
		if len(hash) != md5.Size*2 {
			return nil, ErrBadSizedChecksum
		}

		ZipSize, err := strconv.Atoi(string(rZipSize))
		if err != nil {
			return nil, err
		}
		RawSize, err := strconv.Atoi(string(rSize))
		if err != nil {
			return nil, err
		}

		/* We do this to not create more memory */
		decodedFiles[i].FileName = name
		decodedFiles[i].Checksum = hash
		decodedFiles[i].ZipSize = uint32(ZipSize)
		decodedFiles[i].RawSize = uint32(RawSize)
	}

	return &RobloxPkgManifest{
		decodedFiles,
	}, nil
}

func ParsePkgManifest(manifest []byte) (*RobloxPkgManifest, error) {
	splitData, err := splitPkgManifest(manifest)
	if err != nil {
		return nil, err
	}

	return deserializePkgManifest(splitData)
}

func splitManifest(manifest []byte) (manifestSplit, error) {
	splitData := bytes.Split(manifest, newLine)
	if !isEven(len(splitData)) {
		return nil, ErrBadSizedManifest
	}
	foundData := make(manifestSplit, len(splitData)/2)
	for i := 0; i < len(splitData)/2; i++ {
		index := i * 2
		if len(splitData) < index+1 {
			return nil, ErrBadSizedChunk
		}
		name, hash := splitData[index], splitData[index+1]
		foundData[i] = [2][]byte{name, hash}
	}
	return foundData, nil
}

func deserializeManifest(data manifestSplit) (*RobloxManifest, error) {
	if len(data) < 1 || !isEven(len(data)) {
		return nil, ErrBadSizedManifest
	}
	fmt.Println(data)
	decodedFiles := make([]RobloxFile, len(data))
	for i, chunk := range data {
		/* Decode data and remove any whitespace */
		path, hash := string(chunk[0]), string(chunk[1])

		/* Validate length before copying into a smaller buffer */
		if len(hash) != md5.Size*2 {
			return nil, ErrBadSizedChecksum
		}

		/* We do this to not create more memory */
		decodedFiles[i].Path = path
		decodedFiles[i].Checksum = hash
	}

	return &RobloxManifest{
		decodedFiles,
	}, nil
}

func ParseManifest(manifest []byte) (*RobloxManifest, error) {
	splitData, err := splitManifest(manifest)
	if err != nil {
		return nil, err
	}

	return deserializeManifest(splitData)
}
