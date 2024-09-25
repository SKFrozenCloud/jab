package main

var (
	DBFile = "hashes.db"
	AESKey = ""
)

func main() {
	// Monitor function
	// - Create empty database
	// - Walk through all files and directories
	// 		- Calculate hashes & append to new database
	// - Load old hash database
	// - Compare the hashes
	// - Save new hash database
	// Periodically call monitor funciton with channel
}
