package rbx_manifest

import (
	"errors"
	"fmt"
	"io"
)

func ValidateStream(stream io.Reader, directory string) error {
	Reader := NewStream(stream)
	for {

		path, err := Reader.ReadLine()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return err
			}
		}
		checkSum, err := Reader.ReadLine()
		if err != nil {
			return err
		}

		file := RobloxFile{
			Path:     string(path),
			Checksum: checkSum,
		}

		if err := file.Validate(directory); err != nil {
			return err
		}
	}
}

func (M *RobloxPkgManifest) EncodeStream(stream io.Writer) error {
	if _, err := stream.Write([]byte(Header + NewLine)); err != nil {
		return err
	}
	for _, file := range M.Files {
		if err := file.EncodeStream(stream); err != nil {
			return err
		}
	}
	return nil
}

func (F *RobloxPkgFile) EncodeStream(stream io.Writer) error {
	numberData := fmt.Sprintf("%d%s%d%s", F.ZipSize, NewLine, F.RawSize, NewLine)

	if _, err := stream.Write([]byte(numberData)); err != nil {
		return err
	} else if _, err := stream.Write([]byte(F.FileName + NewLine)); err != nil {
		return err
	} else if _, err := stream.Write(append(F.Checksum, newLine...)); err != nil {
		return err
	} else if _, err := stream.Write([]byte(numberData)); err != nil {
		return err
	}

	return nil
}

func (M *RobloxManifest) EncodeStream(stream io.Writer) error {
	for _, file := range M.Files {
		if err := file.EncodeStream(stream); err != nil {
			return err
		}
	}
	return nil
}

func (F *RobloxFile) EncodeStream(file io.Writer) error {
	if _, err := file.Write([]byte(F.Path + NewLine)); err != nil {
		return err
	} else if _, err := file.Write(append(F.Checksum, newLine...)); err != nil {
		return err
	}
	return nil
}
