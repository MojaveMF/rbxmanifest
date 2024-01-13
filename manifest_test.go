package rbx_manifest_test

import (
	"os"
	rbx_manifest "rbxmanifest"
	"testing"
)

/*
	PkgParse and ManifestParse have the same usage so testing is nearly the same
*/

func TestPkgParse(t *testing.T) {
	manifest, err := os.ReadFile("./testdata/rbxPkgManifest.txt")
	if err != nil {
		t.Error("Failed to read testing data")
		return
	}
	decoded, err := rbx_manifest.ParsePkgManifest(manifest)
	if err != nil {
		t.Error(err)
	} else if len(decoded.Files) != 32 {
		t.Errorf("Expected length of 32 got %d", len(decoded.Files))
	}
}

func TestManifestParse(t *testing.T) {
	manifest, err := os.ReadFile("./testdata/rbxManifest.txt")
	if err != nil {
		t.Error("Failed to read testing data")
		return
	}
	decoded, err := rbx_manifest.ParseManifest(manifest)
	if err != nil {
		t.Error(err)
	} else if len(decoded.Files) != 18094 {
		t.Errorf("Expected length of 18094 got %d", len(decoded.Files))
	}
}

func TestReader(t *testing.T) {
	fileStream, err := os.Open("./testdata/rbxPkgManifest.txt")
	if err != nil {
		t.Error(err)
	} else if _, err := rbx_manifest.ParsePkgManifestStream(fileStream); err != nil {
		t.Error(err)
		return
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
