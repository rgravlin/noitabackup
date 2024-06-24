package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	TestSourcePath = "test_source_path\\save00"
)

type Node struct {
	Name     string
	Children []*Node
}

type DirectoryTree struct {
	Root *Node
}

func TestRestore_RestoreNoita(t *testing.T) {
	// create a mock source save00 directory structure
	if err := newNoitaSourceDirs(); err != nil {
		t.Fatal(err)
	}
	if err := newNoitaBackupDirs(); err != nil {
		t.Fatal(err)
	}

	// create a mock backup directory
	if err := os.MkdirAll(TestBackupPath, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	// create a new restore
	backup := NewBackup(false, false, 16, TestSourcePath, TestBackupPath)
	restore := NewRestore("latest", backup)
	if err := restore.restoreNoita(); err != nil {
		t.Fatal(err)
	}

	// cleanup
	if err := os.RemoveAll(strings.Split(TestSourcePath, "\\")[0]); err != nil {
		t.Fatal(err)
	}

	if err := os.RemoveAll(TestBackupPath); err != nil {
		t.Fatal(err)
	}
}

func newNoitaSourceDirs() error {
	if err := newNoitaSourceTree().Root.createDirectories(""); err != nil {
		return err
	}

	return nil
}

func newNoitaBackupDirs() error {
	if err := newNoitaBackupTree().Root.createDirectories(""); err != nil {
		return err
	}

	return nil
}

func newNoitaSourceTree() *DirectoryTree {
	return &DirectoryTree{
		Root: &Node{
			Name:     TestSourcePath,
			Children: newNoitaDirNode("save00"),
		},
	}
}

func newNoitaBackupTree() *DirectoryTree {
	return &DirectoryTree{
		Root: &Node{
			Name:     TestBackupPath,
			Children: newNoitaDirNode(time.Now().Format(TimeFormat)),
		},
	}

}

func newNoitaDirNode(name string) []*Node {
	return []*Node{
		{
			Name: name,
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
	}
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
