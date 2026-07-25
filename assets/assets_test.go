package assets

import (
	"embed"
	"strings"
	"testing"
)

//go:embed assets.go
var embeddedAssets embed.FS

func TestGetFileByName(t *testing.T) {
	content := GetFileByName(embeddedAssets, "assets.go")
	if !strings.Contains(string(content), "package assets") {
		t.Error("expected embedded assets.go content")
	}
}

func TestGetFileByNamePanics(t *testing.T) {
	testCases := []string{"", "missing"}
	for _, fileName := range testCases {
		t.Run(fileName, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Errorf("expected panic for file name %q", fileName)
				}
			}()
			GetFileByName(embeddedAssets, fileName)
		})
	}
}

func TestListFiles(t *testing.T) {
	files := ListFiles(embeddedAssets)
	if len(files) != 1 || files[0] != "assets.go" {
		t.Errorf("expected [assets.go], got %v", files)
	}
}
