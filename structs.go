package rbx_manifest

type RobloxPkgFile struct {
	Checksum         []byte
	FileName         string
	ZipSize, RawSize int
}

type RobloxPkgManifest struct {
	Files []RobloxPkgFile
}

type RobloxFile struct {
	Checksum []byte
	Path     string
}

type RobloxManifest struct {
	Files []RobloxFile
}
