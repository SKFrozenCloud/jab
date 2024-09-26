package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestSetupFileHashDatabase(t *testing.T) {
	//// Test normal directory containing files and recursive directories with files
	// Setup
	directoryOne, _ := os.MkdirTemp("", "jab_database_test_one")
	fileOne, _ := os.CreateTemp(directoryOne, "fileone")
	fileOne.WriteString("testcontent1")
	fileTwo, _ := os.CreateTemp(directoryOne, "filetwo")
	fileTwo.WriteString("testcontent2")
	directoryTwo, _ := os.MkdirTemp(directoryOne, "jab_database_test_two")
	fileThree, _ := os.CreateTemp(directoryTwo, "filethree")
	directoryThree, _ := os.MkdirTemp(directoryTwo, "jab_database_test_three")
	fileFour, _ := os.CreateTemp(directoryThree, "filefour")
	fileFour.WriteString("testcontent4")

	directoryFour, _ := os.MkdirTemp("", "jab_database_test_four")
	fileFive, _ := os.CreateTemp(directoryFour, "filefive")
	fileFive.WriteString("testcontent5")

	fileSix, _ := os.CreateTemp("", "filesix")
	fileSix.WriteString("testcontent6")
	fileSeven, _ := os.CreateTemp("", "fileseven")
	fileSeven.WriteString("testcontent7")

	correctDB := FileHashDatabase{}
	correctDB[FilePath(fileOne.Name())] = "19e1c0dda4fbdd75d6f0d31e5e28b2a4103b370049909bb485ebacf297c8952f"
	correctDB[FilePath(fileTwo.Name())] = "f25b7c4a9e88d251b5942c7d7a2c6480eef3e07195f9b7063ca332702edeaff9"
	correctDB[FilePath(fileThree.Name())] = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	correctDB[FilePath(fileFour.Name())] = "d12ebbd2f807bf53517dd8e8eed737a919e4d2f2b3952d62a52b23fcd53a3109"
	correctDB[FilePath(fileFive.Name())] = "2b4a26d3dceacce9b09c4faddd5594d6113235d0fb03ff1bd70bf27a2b7f7dc7"
	correctDB[FilePath(fileSix.Name())] = "c2d098c3d20d3d6e321689f26377b72424e72704edfa331d5c30e08d0ef97b32"
	correctDB[FilePath(fileSeven.Name())] = "f004d515b6bd74a6f052bbf5f24b4c08b0d990e8c5811e3721a4642584b9a754"

	// Execute
	db, err := SetupFileHashDatabase([]string{directoryOne, directoryFour, fileSix.Name(), fileSeven.Name()})
	if err != nil {
		t.Error("could not create file hash database")
	}

	// Validate
	/*
		// Debug
		t.Log("generated db")
		for path, hash := range db {
			t.Logf("path: %s, hash: %s", path, hash)
		}
		t.Log("correct db")
		for path, hash := range correctDB {
			t.Logf("path: %s, hash: %s", path, hash)
		}
	*/

	result := reflect.DeepEqual(correctDB, db)
	if !result {
		t.Error("databases not matching")
	}

	// Teardown
	err = os.RemoveAll(directoryOne)
	if err != nil {
		t.Error("could not remove temporary directory")
	}

	//// Test empty directory
	// Setup
	directoryOne, _ = os.MkdirTemp("", "jab_database_test_one")

	correctDB = FileHashDatabase{}

	// Execute
	db, err = SetupFileHashDatabase([]string{directoryOne})
	if err != nil {
		t.Error("could not create file hash database")
	}

	// Validate
	result = reflect.DeepEqual(correctDB, db)
	if !result {
		t.Error("databases not matching")
	}

	// Teardown
	err = os.RemoveAll(directoryOne)
	if err != nil {
		t.Error("could not remove temporary directory")
	}
}

func TestLoadFileHashDatabase(t *testing.T) {
	//// Test normal database
	// Setup DB
	correctDB := FileHashDatabase{}
	correctDB[FilePath("fileone")] = "19e1c0dda4fbdd75d6f0d31e5e28b2a4103b370049909bb485ebacf297c8952f"
	correctDB[FilePath("filetwo")] = "f25b7c4a9e88d251b5942c7d7a2c6480eef3e07195f9b7063ca332702edeaff9"
	correctDB[FilePath("filethree")] = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	correctDB[FilePath("filefour")] = "d12ebbd2f807bf53517dd8e8eed737a919e4d2f2b3952d62a52b23fcd53a3109"
	correctDB[FilePath("filefive")] = "2b4a26d3dceacce9b09c4faddd5594d6113235d0fb03ff1bd70bf27a2b7f7dc7"
	correctDB[FilePath("filesix")] = "c2d098c3d20d3d6e321689f26377b72424e72704edfa331d5c30e08d0ef97b32"
	correctDB[FilePath("fileseven")] = "f004d515b6bd74a6f052bbf5f24b4c08b0d990e8c5811e3721a4642584b9a754"

	// Manually save DB
	jsonBytes, err := json.Marshal(correctDB)
	if err != nil {
		t.Error("could not manually save db")
	}

	fileDB, err := os.CreateTemp("", "hashes.db")
	if err != nil {
		t.Error("could not save db")
	}
	fileDB.Write(jsonBytes)

	// Load DB with function
	db, err := LoadFileHashDatabase(fileDB.Name())
	if err != nil {
		t.Error("could not load db")
	}

	// Validate
	result := reflect.DeepEqual(correctDB, db)
	if !result {
		t.Error("databases not matching")
	}

	// Teardown
	err = os.Remove(fileDB.Name())
	if err != nil {
		t.Error("could not remove temporary file")
	}

	//// Test empty database
	// Setup DB
	correctDB = FileHashDatabase{}

	// Manually save DB
	jsonBytes, err = json.Marshal(correctDB)
	if err != nil {
		t.Error("could not manually save db")
	}

	fileDB, err = os.CreateTemp("", "hashes.db")
	if err != nil {
		t.Error("could not save db")
	}
	fileDB.Write(jsonBytes)

	// Load DB with function
	db, err = LoadFileHashDatabase(fileDB.Name())
	if err != nil {
		t.Error("could not load db")
	}

	// Validate
	result = reflect.DeepEqual(correctDB, db)
	if !result {
		t.Error("databases not matching")
	}

	// Teardown
	err = os.Remove(fileDB.Name())
	if err != nil {
		t.Error("could not remove temporary file")
	}
}

func TestSaveFileHashDatabase(t *testing.T) {
	//// Test normal database
	// Setup DB
	correctDB := FileHashDatabase{}
	correctDB[FilePath("fileone")] = "19e1c0dda4fbdd75d6f0d31e5e28b2a4103b370049909bb485ebacf297c8952f"
	correctDB[FilePath("filetwo")] = "f25b7c4a9e88d251b5942c7d7a2c6480eef3e07195f9b7063ca332702edeaff9"
	correctDB[FilePath("filethree")] = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	correctDB[FilePath("filefour")] = "d12ebbd2f807bf53517dd8e8eed737a919e4d2f2b3952d62a52b23fcd53a3109"
	correctDB[FilePath("filefive")] = "2b4a26d3dceacce9b09c4faddd5594d6113235d0fb03ff1bd70bf27a2b7f7dc7"
	correctDB[FilePath("filesix")] = "c2d098c3d20d3d6e321689f26377b72424e72704edfa331d5c30e08d0ef97b32"
	correctDB[FilePath("fileseven")] = "f004d515b6bd74a6f052bbf5f24b4c08b0d990e8c5811e3721a4642584b9a754"

	// Save DB with function
	fileDB, _ := os.CreateTemp("", "hashes.db")
	err := SaveFileHashDatabase(correctDB, fileDB.Name())
	if err != nil {
		t.Error("could not save db")
	}

	// Manually load DB
	dbBytes, err := os.ReadFile(fileDB.Name())
	if err != nil {
		t.Error("could not load db")
	}

	var db FileHashDatabase
	err = json.Unmarshal([]byte(dbBytes), &db)
	if err != nil {
		t.Error("could not load db")
	}

	// Validate
	result := reflect.DeepEqual(correctDB, db)
	if !result {
		t.Error("databases not matching")
	}

	// Teardown
	err = os.Remove(fileDB.Name())
	if err != nil {
		t.Error("could not remove temporary file")
	}

	//// Test empty database
	// Setup DB
	correctDB = FileHashDatabase{}

	// Save DB with function
	fileDB, _ = os.CreateTemp("", "hashes.db")
	err = SaveFileHashDatabase(correctDB, fileDB.Name())
	if err != nil {
		t.Error("could not save db")
	}

	// Manually load DB
	dbBytes, err = os.ReadFile(fileDB.Name())
	if err != nil {
		t.Error("could not load db")
	}

	var dbEmpty FileHashDatabase
	err = json.Unmarshal([]byte(dbBytes), &dbEmpty)
	if err != nil {
		t.Error("could not load db")
	}

	// Validate
	result = reflect.DeepEqual(correctDB, dbEmpty)
	if !result {
		t.Error("databases not matching")
	}

	// Teardown
	err = os.Remove(fileDB.Name())
	if err != nil {
		t.Error("could not remove temporary file")
	}
}
