package rbx_manifest

type RobloxPkgFile struct {
	FileName, Checksum string
	ZipSize, RawSize   uint32
}

type RobloxPkgManifest struct {
	Files []RobloxPkgFile
}

type RobloxFile struct {
	Checksum, Path string
}

type RobloxManifest struct {
	Files []RobloxFile
}
