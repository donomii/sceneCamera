package assets

//Assets
//
// The assets module works with an embedded directory that you provide.  Due to the terrible design of the embed module, you must do your own embedding.  The assets module provides a few helper functions to work with the embedded directory.
//
// To embed files in your program, you must use the go:embed directive.  For example:
//
// //go:embed assets/*
// var embeddedFS embed.FS
//
// You can then pass the variable embeddedFS to the assets module.  For example:
//
// // Get the content of a file by name.
// content, err := assets.GetFileByName(embeddedFS, "assets/myfile.txt")
// if err != nil {
// 	log.Fatal(err)
// }

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
)

// GetFileByName retrieves a file's content by name from the provided embedded file system.
// It panics if the file cannot be opened, read, or closed.
func GetFileByName(embeddedFS embed.FS, fileName string) []byte {
	if fileName == "" {
		panic(errors.New("file name cannot be empty"))
	}

	file, err := embeddedFS.Open(fileName)
	if err != nil {
		panic(fmt.Errorf("open embedded file %q: %w", fileName, err))
	}

	content, readErr := io.ReadAll(file)
	closeErr := file.Close()
	if readErr != nil {
		panic(fmt.Errorf("read embedded file %q: %w", fileName, readErr))
	}
	if closeErr != nil {
		panic(fmt.Errorf("close embedded file %q: %w", fileName, closeErr))
	}

	return content
}

// ListFiles returns a list of FileInfo for all files in the provided embedded file system.
func ListFiles(embeddedFS embed.FS) []string {
	dir, err := fs.ReadDir(embeddedFS, ".")
	if err != nil {
		panic(err)
	}

	files := make([]string, 0, len(dir))
	for _, entry := range dir {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files
}
