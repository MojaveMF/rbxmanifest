package rbx_manifest

import "fmt"

var (
	Header  = "v0"
	NewLine = "\r\n"
)

func (M *RobloxPkgManifest) Encode() []byte {
	newManifest := []byte(Header + NewLine)

	for _, file := range M.Files {
		newManifest = append(newManifest, file.Encode()...)
	}

	return newManifest
}

func (F *RobloxPkgFile) Encode() []byte {
	numberData := fmt.Sprintf("%d%s%d%s", F.ZipSize, NewLine, F.RawSize, NewLine)
	fileData := []byte{}
	fileData = append(fileData, []byte(F.FileName+NewLine)...)
	fileData = append(fileData, append(F.Checksum, newLine...)...)
	fileData = append(fileData, []byte(numberData)...)
	return fileData
}
