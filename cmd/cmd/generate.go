/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	rbx_manifest "rbxmanifest"

	"github.com/spf13/cobra"
)

var (
	InDir, OutDir                                  string
	Verbose, GeneratePkgManifest, GenerateManifest bool
)

func VerboseLog(format string, arguments ...any) {
	if Verbose {
		fmt.Printf(format+"\n", arguments...)
	}
}

func ValidateExistance(folder string) {
	VerboseLog("Checking if %s exist", folder)

	if _, err := os.Stat(folder); err == nil || os.IsExist(err) {
		return
	}

	fmt.Printf("The path %s dosent exist \n", folder)
	os.Exit(1)
}

func GetFile(path string) *os.File {
	VerboseLog("Creating file %s", path)
	stream, err := os.Create(path)
	if err != nil {
		fmt.Printf("Couldnt create file at %s \n", path)
		os.Exit(1)
	}
	return stream
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a Manifest info",
	Long: `This command can generate Manifest info.
By default it will create them for the current path.
This can be configured by the cli.`,
	Run: func(cmd *cobra.Command, args []string) {

		inDir, err := filepath.Abs(InDir)
		if err != nil {
			fmt.Printf("Couldnt get abs path for %s", InDir)
			os.Exit(1)
		}
		outDir, err := filepath.Abs(OutDir)
		if err != nil {
			fmt.Printf("Couldnt get abs path for %s", OutDir)
			os.Exit(1)
		}

		InDir = inDir
		OutDir = outDir

		ValidateExistance(InDir)
		ValidateExistance(OutDir)

		if GeneratePkgManifest {
			outLocation := filepath.Join(OutDir, "rbxPkgManifest.txt")
			stream := GetFile(outLocation)

			rbx_manifest.PkgGenerateDirectory(stream, InDir, filepath.Join(OutDir, "package"), Verbose)
		}
		if GenerateManifest {
			outLocation := filepath.Join(OutDir, "rbxManifest.txt")
			stream := GetFile(outLocation)

			rbx_manifest.GenerateDirectory(stream, InDir, Verbose)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&InDir, "in", "i", "./", "Directory to generate values from")
	generateCmd.Flags().StringVarP(&OutDir, "out", "o", "../out", "Directory to generate files from")
	generateCmd.Flags().BoolVarP(&Verbose, "verbose", "v", true, "If uneeded output should be muted")
	generateCmd.Flags().BoolVarP(&GeneratePkgManifest, "pkgManifest", "p", true, "If a PkgManifest should be generated")
	generateCmd.Flags().BoolVarP(&GenerateManifest, "manifest", "m", true, "If a manifest should be generated")
}
