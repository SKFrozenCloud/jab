package main

import (
	"io/fs"
	"path/filepath"
)

type FilePath string
type FileHash string

type FileHashDatabase map[FilePath]FileHash

func SetupFileHashDatabase(directoryNames []string) (FileHashDatabase, error) {
	fileHashDatabase := FileHashDatabase{}

	// Walk through all files and directories
	for _, v := range directoryNames {
		filepath.WalkDir(
			v,
			func(path string, d fs.DirEntry, err error) error {
				if d.IsDir() || err != nil {
					return nil
				}

				// Calculate hashes & append to new database
				hash, err2 := HashFile(path)
				if err2 != nil {
					return err2
				}

				fileHashDatabase[FilePath(path)] = FileHash(hash)

				return nil
			},
		)
	}

	return fileHashDatabase, nil
}
