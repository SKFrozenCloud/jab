package main

type IntegrityChangesAdded struct {
	FilePath FilePath
	FileHash FileHash
}

type IntegrityChangesModified struct {
	FilePath    FilePath
	NewFileHash FileHash
	OldFileHash FileHash
}

type IntegrityChangesRemoved struct {
	FilePath FilePath
	FileHash FileHash
}

type IntegrityChanges struct {
	Added    []IntegrityChangesAdded
	Modified []IntegrityChangesModified
	Removed  []IntegrityChangesRemoved
}

func CheckIntegrity(dbPath string, checkPaths []string) (IntegrityChanges, error) {
	//// Setup a new database
	newDB, err := SetupFileHashDatabase(checkPaths)
	if err != nil {
		return IntegrityChanges{}, err
	}

	//// Load old hash database
	oldDB, oldDBError := LoadFileHashDatabase(dbPath)

	//// Compare the hashes
	integrityChanges := IntegrityChanges{
		Added:    []IntegrityChangesAdded{},
		Modified: []IntegrityChangesModified{},
		Removed:  []IntegrityChangesRemoved{},
	}

	if oldDBError == nil {
		//// Database did exist before
		// Check for added and modified files
		for filePath, fileHash := range newDB {
			oldDBFileHash, ok := oldDB[filePath]

			if !ok {
				// New file
				integrityChanges.Added = append(integrityChanges.Added, IntegrityChangesAdded{
					FilePath: filePath,
					FileHash: fileHash,
				})
			} else if fileHash != oldDBFileHash {
				// Modified file
				integrityChanges.Modified = append(integrityChanges.Modified, IntegrityChangesModified{
					FilePath:    filePath,
					OldFileHash: oldDBFileHash,
					NewFileHash: fileHash,
				})
			}
		}

		// Check for removed filed
		for filePath, fileHash := range oldDB {
			_, ok := newDB[filePath]

			if !ok {
				// Removed file
				integrityChanges.Removed = append(integrityChanges.Removed, IntegrityChangesRemoved{
					FilePath: filePath,
					FileHash: fileHash,
				})
			}
		}
	} else {
		//// Database did not exist before
		// Check for added files
		for filePath, fileHash := range newDB {
			// New file
			integrityChanges.Added = append(integrityChanges.Added, IntegrityChangesAdded{
				FilePath: filePath,
				FileHash: fileHash,
			})
		}
	}

	//// Save new hash database
	err = SaveFileHashDatabase(newDB, dbPath)
	if err != nil {
		return IntegrityChanges{}, err
	}

	return integrityChanges, nil
}
