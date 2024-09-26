package main

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

type FilePath string
type FileHash string

type FileHashDatabase map[FilePath]FileHash

func SetupFileHashDatabase(paths []string) (FileHashDatabase, error) {
	fileHashDatabase := FileHashDatabase{}

	//// Walk through all files and directories

	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			panic(err)
		}

		if !stat.IsDir() {
			// Files
			// Calculate hashes & append to new database
			databaseHashFile(path, fileHashDatabase)
		} else {
			// Directories
			filepath.WalkDir(
				path,
				func(path string, d fs.DirEntry, err error) error {
					if d.IsDir() || err != nil {
						return nil
					}

					// Calculate hashes & append to new database
					databaseHashFile(path, fileHashDatabase)

					return nil
				},
			)
		}
	}

	return fileHashDatabase, nil
}

func databaseHashFile(path string, db FileHashDatabase) error {
	hash, err := HashFile(path)
	if err != nil {
		return err
	}

	db[FilePath(path)] = FileHash(hash)
	return nil
}

func LoadFileHashDatabase(databasePath string) (FileHashDatabase, error) {
	dbBytes, err := os.ReadFile(databasePath)
	if err != nil {
		return nil, err
	}

	var db FileHashDatabase
	err = json.Unmarshal([]byte(dbBytes), &db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SaveFileHashDatabase(db FileHashDatabase, databasePath string) error {
	jsonBytes, err := json.Marshal(db)
	if err != nil {
		return err
	}

	fileDB, err := os.Create(databasePath)
	if err != nil {
		return err
	}
	fileDB.Write(jsonBytes)

	return nil
}
