package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	TestBackupPath = "test_backup_path"
	TestSourcePath = "test_source_path"
)

type Node struct {
	Name     string
	Children []*Node
}

type DirectoryTree struct {
	Root *Node
}

func TestByDate(t *testing.T) {
	tests := []struct {
		name string
		ByDate
		expected bool
	}{
		{"test1", ByDate{time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC), time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC)}, true},
		{"test2", ByDate{time.Date(2023, 01, 01, 0, 0, 0, 0, time.UTC), time.Date(2022, 01, 01, 0, 0, 0, 0, time.UTC)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ByDate.Less(0, 1); got != tt.expected {
				t.Errorf("ByDate.Less() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestCleanBackups(t *testing.T) {
	mockBackupDirs := createMockBackupDirs(t)

	tests := []struct {
		name      string
		dirs      []time.Time
		numToKeep int
		expectErr bool
	}{
		{name: "Keep all backups", dirs: mockBackupDirs, numToKeep: len(mockBackupDirs), expectErr: false},
		{name: "Keep all except delete one", dirs: mockBackupDirs, numToKeep: len(mockBackupDirs) - 1, expectErr: false},
		{name: "Clean all backups except one", dirs: mockBackupDirs, numToKeep: 1, expectErr: false},
		{name: "Clean all backups and fail", dirs: mockBackupDirs, numToKeep: 0, expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBackup(false, false, tt.numToKeep, "", TestBackupPath)
			err := b.cleanBackups()
			if (err != nil) != tt.expectErr {
				t.Fatalf("cleanBackups() error = %v, expected %v", err, tt.expectErr)
			}
		})
	}

	err := os.RemoveAll(TestBackupPath)
	if err != nil {
		t.Fatal(err)
	}
}

func createMockBackupDirs(t *testing.T) []time.Time {
	t.Helper()

	dirs := make([]time.Time, 10)
	for i := 0; i < 10; i++ {
		dirs[i] = time.Now().Add(time.Duration(i) * time.Minute)
		dirPath := filepath.Join(TestBackupPath, dirs[i].Format(TimeFormat))
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			t.Fatal(err)
		}
	}
	return dirs
}

func TestRestore_RestoreNoita(t *testing.T) {
	// create a mock source save00 directory structure
	if err := createMockSourceDir(t); err != nil {
		t.Fatal(err)
	}
	if err := createMockBackupDir(t); err != nil {
		t.Fatal(err)
	}

	// create a mock backup directory
	if err := os.MkdirAll(TestBackupPath, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// create a new restore
	backup := NewBackup(false, false, 16, TestSourcePath, TestBackupPath)
	restore := NewRestore("latest", backup)
	restore.restoreNoita()

	// cleanup
	if err := os.RemoveAll(TestSourcePath); err != nil {
		t.Fatal(err)
	}

	if err := os.RemoveAll(TestBackupPath); err != nil {
		t.Fatal(err)
	}
}

func createMockBackupDir(t *testing.T) error {
	t.Helper()
	tree := DirectoryTree{
		Root: &Node{
			Name: TestBackupPath,
			Children: []*Node{
				{
					Name: fmt.Sprintf("%s", time.Now().Format(TimeFormat)),
					Children: []*Node{
						{
							Name: "world",
						},
						{
							Name: "persistent",
							Children: []*Node{
								{Name: "bones"},
								{Name: "bones_new"},
								{Name: "flags"},
								{Name: "orbs_new"},
							},
						},
						{
							Name: "stats",
							Children: []*Node{
								{Name: "sessions"},
							},
						},
					},
				},
			},
		},
	}

	if err := tree.Root.createDirectories(""); err != nil {
		return err
	}

	return nil
}

func createMockSourceDir(t *testing.T) error {
	t.Helper()
	tree := DirectoryTree{
		Root: &Node{
			Name: TestSourcePath,
			Children: []*Node{
				{
					Name: "save00",
					Children: []*Node{
						{
							Name: "world",
						},
						{
							Name: "persistent",
							Children: []*Node{
								{Name: "bones"},
								{Name: "bones_new"},
								{Name: "flags"},
								{Name: "orbs_new"},
							},
						},
						{
							Name: "stats",
							Children: []*Node{
								{Name: "sessions"},
							},
						},
					},
				},
			},
		},
	}

	if err := tree.Root.createDirectories(""); err != nil {
		return err
	}

	return nil
}

func (n *Node) createDirectories(path string) error {
	newPath := filepath.Join(path, n.Name)
	err := os.MkdirAll(newPath, 0755)
	if err != nil {
		return err
	}

	for _, child := range n.Children {
		err = child.createDirectories(newPath)
		if err != nil {
			return err
		}
	}

	return nil
}
