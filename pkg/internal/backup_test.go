package internal

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	TestBackupPath = "test_backup_path"
)

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
