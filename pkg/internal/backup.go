package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	TimeFormat                = "2006-01-02-15-04-05"
	ConfigMaxNumBackupsToKeep = 64.00
	SteamExe                  = "C:\\Program Files (x86)\\Steam\\steam.exe"
	SteamNoitaFlags           = "steam://rungameid/881100"
	ExplorerExe               = "explorer"
)

type Backup struct {
	async             bool
	autoLaunchChecked bool
	maxBackups        int
	dirCounter        int
	fileCounter       int
	srcPath           string
	dstPath           string
	phase             int
	timestamp         time.Time
	sortedBackupDirs  []time.Time
	LogRing           *LogRing
}

func NewBackup(async, autoLaunchChecked bool, maxBackups int, srcPath, dstPath string) *Backup {
	return &Backup{
		async:             async,
		autoLaunchChecked: autoLaunchChecked,
		maxBackups:        maxBackups,
		srcPath:           srcPath,
		dstPath:           dstPath,
	}
}

func (b *Backup) BackupNoita() {
	if !isNoitaRunning() {
		if b.phase == stopped {
			if b.async {
				go b.backupNoita()
			} else {
				b.backupNoita()
			}
		} else {
			b.LogRing.LogAndAppend("operation already in progress")
		}
	} else {
		b.LogRing.LogAndAppend("noita.exe cannot be running to backup")
	}
}

func (b *Backup) backupNoita() {
	t := time.Now()
	b.timestamp = t
	b.phase = started
	b.reportStart()

	newBackupPath := fmt.Sprintf("%s\\%s", b.dstPath, b.timestamp.Format(TimeFormat))

	// get current number of backups
	curNumBackups, err := getNumBackups(b.dstPath)
	if err != nil {
		b.LogRing.LogAndAppend(fmt.Sprintf("error getting backups: %v", err))
		b.phase = stopped
		return
	} else {
		b.LogRing.LogAndAppend(fmt.Sprintf("number of backups: %d", curNumBackups))
	}

	// protect against invalid maxBackups
	if b.maxBackups > ConfigMaxNumBackupsToKeep || b.maxBackups <= 0 {
		b.maxBackups = ConfigMaxNumBackupsToKeep
	}

	// clean up backups
	if curNumBackups >= b.maxBackups {
		b.LogRing.LogAndAppend("maximum backup threshold reached")

		// get and sort backup directories
		// oldest are first in the sorted slice
		b.sortedBackupDirs, err = getBackupDirs(b.dstPath, TimeFormat)
		if err != nil {
			b.LogRing.LogAndAppend(fmt.Sprintf("error getting backups: %v", err))
			b.phase = stopped
			return
		}

		// clean backup directories to make room for this backup
		if err := b.cleanBackups(); err != nil {
			b.LogRing.LogAndAppend(fmt.Sprintf("failure deleting backups: %v", err))
			b.phase = stopped
			return
		}
	}

	// create new backup path
	if err := createIfNotExists(newBackupPath, 0755); err != nil {
		b.LogRing.LogAndAppend(fmt.Sprintf("cannot create destination path: %v", err))
		b.phase = stopped
		return
	}

	// recursively copy source to destination
	if err := copyDirectory(b.srcPath, newBackupPath, &b.dirCounter, &b.fileCounter); err != nil {
		b.LogRing.LogAndAppend(fmt.Sprintf("cannot copy source to destination: %v", err))
		b.phase = stopped
		return
	}

	b.reportStop()
	b.resetPhase()

	if b.autoLaunchChecked {
		err = LaunchNoita(b.async)
		if err != nil {
			b.LogRing.LogAndAppend(fmt.Sprintf("failed to launch noita: %v", err))
		}
	}
}

func (b *Backup) resetPhase() {
	b.phase = stopped
	b.dirCounter = 0
	b.fileCounter = 0
}

func (b *Backup) reportStart() {
	b.LogRing.LogAndAppend(fmt.Sprintf("timestamp: %s", b.timestamp))
	b.LogRing.LogAndAppend(fmt.Sprintf("source: %s", b.srcPath))
	b.LogRing.LogAndAppend(fmt.Sprintf("destination: %s", fmt.Sprintf("%s\\%s", b.dstPath, b.timestamp.Format(TimeFormat))))
}

func (b *Backup) reportStop() {
	b.LogRing.LogAndAppend(fmt.Sprintf("timestamp: %s", time.Now()))
	b.LogRing.LogAndAppend(fmt.Sprintf("total time: %s", time.Since(b.timestamp)))
	b.LogRing.LogAndAppend(fmt.Sprintf("total dirs copied: %d", b.dirCounter))
	b.LogRing.LogAndAppend(fmt.Sprintf("total files copied: %d", b.fileCounter))
}

func (b *Backup) cleanBackups() error {
	totalBackups := len(b.sortedBackupDirs)
	totalToRemove := totalBackups - (b.maxBackups - 1)

	for i := 0; i < totalToRemove; i++ {
		folder := fmt.Sprintf("%s\\%s", b.dstPath, b.sortedBackupDirs[i].Format(TimeFormat))
		b.LogRing.LogAndAppend(fmt.Sprintf("removing backup folder: %s", folder))
		err := os.RemoveAll(folder)
		if err != nil {
			return err
		}
	}

	return nil
}

func getBackupDirs(backupPath, timePattern string) ([]time.Time, error) {
	var backupDirs []time.Time
	if !exists(backupPath) {
		return backupDirs, nil
	} else {
		entries, err := os.ReadDir(backupPath)
		if err != nil {
			return backupDirs, err
		}

		for _, entry := range entries {
			srcPath := filepath.Join(backupPath, entry.Name())
			srcInfo, err := os.Stat(srcPath)
			if err != nil {
				return backupDirs, err
			}

			switch srcInfo.Mode() & os.ModeType {
			case os.ModeDir:
				name := strings.Split(srcPath, "\\")
				nameDate, err := time.Parse(timePattern, name[len(name)-1])
				if err != nil {
					return backupDirs, err
				}
				backupDirs = append(backupDirs, nameDate)
			default:
			}
		}
		sort.Sort(ByDate(backupDirs))
		return backupDirs, nil
	}
}

func getNumBackups(backupPath string) (int, error) {
	numBackup := 0
	backupDirs, err := getBackupDirs(backupPath, TimeFormat)
	if err != nil {
		return numBackup, err
	}

	return len(backupDirs), nil
}

type ByDate []time.Time

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Before(a[j]) }
