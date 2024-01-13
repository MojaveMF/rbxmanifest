package rbx_manifest

import (
	"fmt"
	"io"
)

func VerboseLog(Verbose bool, format string, arguments ...any) {
	if Verbose {
		fmt.Printf(format+"\n", arguments...)
	}
}

func PkgGenerateDirectory(output io.Writer, directory, outdir string, verbose bool) error {
	VerboseLog(verbose, "Generating rbxPkgManifest.txt")

	return nil
}

func GenerateDirectory(output io.Writer, directory string, verbose bool) error {
	VerboseLog(verbose, "Generating rbxManifest.txt")

	return nil
}
