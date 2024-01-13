package rbx_manifest_test

import (
	"os"
	rbx_manifest "rbxmanifest"
	"testing"
)

/*
	If these two functions work all parsing functions work.
*/

func TestPkgParse(t *testing.T) {
	_, err := rbx_manifest.ParsePkgManifestFile("./testdata/rbxPkgManifest.txt")
	if err != nil {
		t.Error(err)
	}
}

func TestManifestParse(t *testing.T) {
	_, err := rbx_manifest.ParseManifestFile("./testdata/rbxManifest.txt")
	if err != nil {
		t.Error(err)
	}
}

func TestPkgEncode(t *testing.T) {
	rManifest, err := os.ReadFile("./testdata/rbxPkgManifest.txt")
	if err != nil {
		t.Error(err)
		return
	}
	decoded, err := rbx_manifest.ParsePkgManifest(rManifest)
	if err != nil {
		t.Error(err)
		return
	}

	file, err := os.Create("./test")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()
	file.Write([]byte(decoded.Encode()))
}

func TestValidate(t *testing.T) {

}
