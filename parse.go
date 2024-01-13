package rbx_manifest

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io"
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

func ParsePkgManifest(manifest []byte) (*RobloxPkgManifest, error) {
	buffer := bytes.Buffer{}
	if _, err := buffer.Write(manifest); err != nil {
		return nil, err
	}

	return ParsePkgManifestStream(&buffer)
}

func splitManifest(manifest []byte) (manifestSplit, error) {
	splitData := bytes.Split(manifest, newLine)
	//splitData = splitData[:len(splitData)-1]

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

	decodedFiles := make([]RobloxFile, len(data))
	for i, chunk := range data {
		/* Decode data and remove any whitespace */
		path, hash := string(chunk[0]), chunk[1]

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

type manifestStream struct {
	bytesLeft       int
	currentLineData []byte
	originStream    io.Reader
}

func NewStream(stream io.Reader) manifestStream {
	return manifestStream{
		bytesLeft:       0,
		currentLineData: []byte{},
		originStream:    stream,
	}
}

func (S *manifestStream) ReadLine() ([]byte, error) {
	for {
		/* Read 1 byte at a time */
		chunk := make([]byte, 32)
		i, err := S.originStream.Read(chunk)
		if errors.Is(err, io.EOF) {
			S.bytesLeft = i
		} else if err != nil {
			return nil, err
		} else {
			S.currentLineData = append(S.currentLineData, chunk...)
		}

		if index := bytes.Index(S.currentLineData, newLine); index != -1 {
			fetchedData := S.currentLineData[:index]
			S.currentLineData = S.currentLineData[index+len(newLine):]
			S.bytesLeft = index - len(newLine)

			return fetchedData, nil
		} else if err != nil {
			S.bytesLeft = 0
			return S.currentLineData, io.EOF
		}
	}
}

func ReadLine(stream io.Reader) ([]byte, error) {
	lineData := make([]byte, 0)
	for {
		/* Read 1 byte at a time */
		chunk := make([]byte, 1)
		_, err := stream.Read(chunk)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
		lineData = append(lineData, chunk...)

		if index := bytes.Index(lineData, newLine); index != -1 {
			return lineData[:index], err
		} else if err != nil {
			return lineData, err
		}
	}
}

func ParsePkgManifestStream(stream io.Reader) (*RobloxPkgManifest, error) {
	parsedFiles := []RobloxPkgFile{}
	Reader := NewStream(stream)
	for {
		startLine, err := Reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return &RobloxPkgManifest{parsedFiles}, nil
			}
			return nil, err
		} else if string(startLine) == Header {
			continue
		}

		checkSum, err := Reader.ReadLine()
		if err != nil {
			return nil, err
		}
		rawZipSize, err := Reader.ReadLine()
		if err != nil {
			return nil, err
		}
		rawRawSize, err := Reader.ReadLine()
		if err != nil {
			return nil, err
		}

		zipSize, err := strconv.Atoi(string(rawZipSize))
		if err != nil {
			return nil, err
		}
		rawSize, err := strconv.Atoi(string(rawRawSize))
		if err != nil {
			return nil, err
		}

		file := RobloxPkgFile{
			FileName: string(startLine),
			Checksum: checkSum,
			RawSize:  rawSize,
			ZipSize:  zipSize,
		}
		parsedFiles = append(parsedFiles, file)
	}
}
