package main

type FilePath string
type FileHash string

type FileHashDatabase map[FilePath]FileHash

func SetupFileHashDatabase(directoryNames []string) (FileHashDatabase, error) {
	return FileHashDatabase{}, nil
}
