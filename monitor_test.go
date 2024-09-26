package main

import (
	"os"
	"sort"
	"testing"
)

func TestMonitorIntegrityOnce(t *testing.T) {
	//// Test a normal scenario
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

	correctAddedChanges := []IntegrityChangesAdded{
		{
			FilePath: FilePath(fileOne.Name()),
			FileHash: "19e1c0dda4fbdd75d6f0d31e5e28b2a4103b370049909bb485ebacf297c8952f",
		},
		{
			FilePath: FilePath(fileTwo.Name()),
			FileHash: "f25b7c4a9e88d251b5942c7d7a2c6480eef3e07195f9b7063ca332702edeaff9",
		},
		{
			FilePath: FilePath(fileThree.Name()),
			FileHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			FilePath: FilePath(fileFour.Name()),
			FileHash: "d12ebbd2f807bf53517dd8e8eed737a919e4d2f2b3952d62a52b23fcd53a3109",
		},
		{
			FilePath: FilePath(fileFive.Name()),
			FileHash: "2b4a26d3dceacce9b09c4faddd5594d6113235d0fb03ff1bd70bf27a2b7f7dc7",
		},
		{
			FilePath: FilePath(fileSix.Name()),
			FileHash: "c2d098c3d20d3d6e321689f26377b72424e72704edfa331d5c30e08d0ef97b32",
		},
		{
			FilePath: FilePath(fileSeven.Name()),
			FileHash: "f004d515b6bd74a6f052bbf5f24b4c08b0d990e8c5811e3721a4642584b9a754",
		},
	}

	// Execute
	fileDB, _ := os.CreateTemp("", "hashes.db")
	integrityChanges, err := MonitorIntegrityOnce(fileDB.Name(), []string{
		directoryOne,
		directoryFour,
		fileSix.Name(),
		fileSeven.Name(),
	})
	if err != nil {
		t.Error("could not monitor integrity")
	}

	if len(integrityChanges.Added) != 7 || len(integrityChanges.Modified) != 0 || len(integrityChanges.Removed) != 0 {
		t.Errorf("bad integrity change; added %v; modified %v; removed %v", len(integrityChanges.Added), len(integrityChanges.Modified), len(integrityChanges.Removed))
	}

	// Validate
	sort.Slice(correctAddedChanges, func(i, j int) bool {
		return correctAddedChanges[i].FilePath < correctAddedChanges[j].FilePath
	})
	sort.Slice(integrityChanges.Added, func(i, j int) bool {
		return integrityChanges.Added[i].FilePath < integrityChanges.Added[j].FilePath
	})
	for index, value := range integrityChanges.Added {
		if value.FilePath != correctAddedChanges[index].FilePath || value.FileHash != correctAddedChanges[index].FileHash {
			t.Error("integrity not preserved")
		}
	}

	//// Test a changed scenario with files added, removed and modified
	// Setup
	fileEight, _ := os.CreateTemp(directoryTwo, "fileeigth")
	fileEight.WriteString("testcontent8")
	fileTwo.WriteString("testcontent2")
	fileThree.WriteString("testcontent3")
	fileSix.WriteString("testcontent6")
	os.Remove(fileFour.Name())
	os.Remove(fileSeven.Name())

	correctAddedChanges = []IntegrityChangesAdded{
		{
			FilePath: FilePath(fileEight.Name()),
			FileHash: "87302530342d37cddef11a7d6f06808c7c1625907c60083ca0f71e48d1114537",
		},
	}

	correctModifiedChanges := []IntegrityChangesModified{
		{
			FilePath:    FilePath(fileTwo.Name()),
			OldFileHash: "f25b7c4a9e88d251b5942c7d7a2c6480eef3e07195f9b7063ca332702edeaff9",
			NewFileHash: "aeb3d6b8fb1c94ea4e1a08b8d26eb43dc53d911c149a9039e8869e51b5ff7daa",
		},
		{
			FilePath:    FilePath(fileThree.Name()),
			OldFileHash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			NewFileHash: "e10a2b08a5b491fa24d934edff560bd8e0154b31165b8ef1214c59ab02ccc838",
		},
		{
			FilePath:    FilePath(fileSix.Name()),
			OldFileHash: "c2d098c3d20d3d6e321689f26377b72424e72704edfa331d5c30e08d0ef97b32",
			NewFileHash: "09ab606d95f819e8371f12bda5dd0af4d18507b648946a5f72a3f282317c8bd2",
		},
	}

	correctRemovedChanges := []IntegrityChangesRemoved{
		{
			FilePath: FilePath(fileFour.Name()),
			FileHash: "d12ebbd2f807bf53517dd8e8eed737a919e4d2f2b3952d62a52b23fcd53a3109",
		},
		{
			FilePath: FilePath(fileSeven.Name()),
			FileHash: "f004d515b6bd74a6f052bbf5f24b4c08b0d990e8c5811e3721a4642584b9a754",
		},
	}

	// Execute
	integrityChanges, err = MonitorIntegrityOnce(fileDB.Name(), []string{
		directoryOne,
		directoryFour,
		fileSix.Name(),
		fileSeven.Name(),
	})
	if err != nil {
		t.Error("could not monitor integrity")
	}

	if len(integrityChanges.Added) != 1 || len(integrityChanges.Modified) != 3 || len(integrityChanges.Removed) != 2 {
		t.Errorf("bad integrity change; added %v; modified %v; removed %v", len(integrityChanges.Added), len(integrityChanges.Modified), len(integrityChanges.Removed))
	}

	// Validate
	sort.Slice(correctAddedChanges, func(i, j int) bool {
		return correctAddedChanges[i].FilePath < correctAddedChanges[j].FilePath
	})
	sort.Slice(integrityChanges.Added, func(i, j int) bool {
		return integrityChanges.Added[i].FilePath < integrityChanges.Added[j].FilePath
	})
	for index, value := range integrityChanges.Added {
		if value.FilePath != correctAddedChanges[index].FilePath || value.FileHash != correctAddedChanges[index].FileHash {
			t.Error("integrity not preserved")
		}
	}

	sort.Slice(correctModifiedChanges, func(i, j int) bool {
		return correctModifiedChanges[i].FilePath < correctModifiedChanges[j].FilePath
	})
	sort.Slice(integrityChanges.Modified, func(i, j int) bool {
		return integrityChanges.Modified[i].FilePath < integrityChanges.Modified[j].FilePath
	})
	for index, value := range integrityChanges.Modified {
		if value.FilePath != correctModifiedChanges[index].FilePath ||
			value.OldFileHash != correctModifiedChanges[index].OldFileHash ||
			value.NewFileHash != correctModifiedChanges[index].NewFileHash {
			t.Error("integrity not preserved")
		}
	}

	sort.Slice(correctRemovedChanges, func(i, j int) bool {
		return correctRemovedChanges[i].FilePath < correctRemovedChanges[j].FilePath
	})
	sort.Slice(integrityChanges.Removed, func(i, j int) bool {
		return integrityChanges.Removed[i].FilePath < integrityChanges.Removed[j].FilePath
	})
	for index, value := range integrityChanges.Removed {
		if value.FilePath != correctRemovedChanges[index].FilePath || value.FileHash != correctRemovedChanges[index].FileHash {
			t.Error("integrity not preserved")
		}
	}

	//// Teardown
	os.Remove(fileDB.Name())
	os.RemoveAll(directoryOne)
	os.RemoveAll(directoryFour)
	os.Remove(fileSix.Name())
	os.Remove(fileSeven.Name())
}
