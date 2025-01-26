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

var DatabaseEncryptionKey string = "ffa321e848eb4fef817376988bbeff80"

func SetupFileHashDatabase(paths []string) (FileHashDatabase, error) {
	fileHashDatabase := FileHashDatabase{}

	//// Walk through all files and directories

	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			continue
		}

		if !stat.IsDir() {
			// Files
			// Calculate hashes & append to new database
			addHashFileToDatabase(path, fileHashDatabase)
		} else {
			// Directories
			filepath.WalkDir(
				path,
				func(path string, d fs.DirEntry, err error) error {
					if d.IsDir() || err != nil {
						return nil
					}

					// Calculate hashes & append to new database
					addHashFileToDatabase(path, fileHashDatabase)

					return nil
				},
			)
		}
	}

	return fileHashDatabase, nil
}

func addHashFileToDatabase(path string, db FileHashDatabase) error {
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

	dbBytesDecrypted, err := DecryptAndVerify(string(dbBytes), DatabaseEncryptionKey)
	if err != nil {
		return nil, err
	}

	var db FileHashDatabase
	err = json.Unmarshal([]byte(dbBytesDecrypted), &db)
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

	jsonBytesEncrypted, err := SignAndEncrypt(string(jsonBytes), DatabaseEncryptionKey)
	if err != nil {
		return err
	}

	fileDB.Write([]byte(jsonBytesEncrypted))

	return nil
}
