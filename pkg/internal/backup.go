package internal

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	ConfigMaxNumBackupsToKeep = 64.00
	ConfigMaxWorkers          = 32.00
	ExplorerExe               = "explorer"
	SteamNoitaFlags           = "steam://rungameid/881100"
	TimeFormat                = "2006-01-02-15-04-05"
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
		LogRing:           NewLogRing(1),
	}
}

func (b *Backup) BackupNoita() {
	if !isNoitaRunning() {
		if b.phase == stopped {
			if b.async {
				go func() { _ = b.backupNoita() }()
			} else {
				_ = b.backupNoita()
			}
		} else {
			b.LogRing.LogAndAppend(ErrOperationAlreadyInProgress)
		}
	} else {
		b.LogRing.LogAndAppend(fmt.Sprintf("%s %s", ErrNoitaRunning, ErrDuringBackup))
	}
}

func (b *Backup) backupNoita() error {
	b.timestamp = time.Now()
	b.phase = started
	b.reportStart()

	newBackupPath := fmt.Sprintf("%s\\%s", b.dstPath, b.timestamp.Format(TimeFormat))

	// get current number of backups
	curNumBackups, err := getNumBackups(b.dstPath)
	if err != nil {
		return b.backupPost(newBackupPath, fmt.Sprintf("%s: %v", ErrErrorGettingBackups, err))
	} else {
		b.LogRing.LogAndAppend(fmt.Sprintf("%s: %d", InfoNumberOfBackups, curNumBackups))
	}

	// protect against invalid maxBackups
	if b.maxBackups > ConfigMaxNumBackupsToKeep || b.maxBackups <= 0 {
		b.maxBackups = ConfigMaxNumBackupsToKeep
	}

	// clean up backups
	if curNumBackups >= b.maxBackups {
		b.LogRing.LogAndAppend(ErrMaxBackupsExceeded)

		// get and sort backup directories
		// oldest are first in the sorted slice
		b.sortedBackupDirs, err = getBackupDirs(b.dstPath, TimeFormat)
		if err != nil {
			return b.backupPost(newBackupPath, fmt.Sprintf("%s: %v", ErrErrorGettingBackups, err))
		}

		// clean backup directories to make room for this backup
		if err := b.cleanBackups(); err != nil {
			return b.backupPost(newBackupPath, fmt.Sprintf("%s: %v", ErrFailureDeletingBackups, err))
		}
	}

	// create new backup path
	if err := createIfNotExists(newBackupPath, 0755); err != nil {
		return b.backupPost(newBackupPath, fmt.Sprintf("%s: %v", ErrCannotCreateDestination, err))
	}

	// recursively copy source to destination
	if err := concurrentCopy(b.srcPath, newBackupPath, &b.dirCounter, &b.fileCounter, viper.GetInt("num-workers")); err != nil {
		return b.backupPost(newBackupPath, fmt.Sprintf("%s: %v", ErrWorkerFailed, err))
	}

	b.reportStop()
	b.resetPhase()

	if b.autoLaunchChecked {
		err = LaunchNoita(b.async)
		if err != nil {
			b.LogRing.LogAndAppend(fmt.Sprintf("%s: %v", ErrFailedToLaunch, err))
		}
	}

	return nil
}

func (b *Backup) resetPhase() {
	b.phase = stopped
	b.dirCounter = 0
	b.fileCounter = 0
}

func (b *Backup) reportStart() {
	b.LogRing.LogAndAppend(fmt.Sprintf("%s: %s", InfoTimestamp, b.timestamp.Format(LogRingTimeFormat)))
	b.LogRing.LogAndAppend(fmt.Sprintf("%s: %s", InfoSource, b.srcPath))
	b.LogRing.LogAndAppend(fmt.Sprintf("%s: %s", InfoDestination, fmt.Sprintf("%s\\%s", b.dstPath, b.timestamp.Format(TimeFormat))))
}

func (b *Backup) reportStop() {
	b.LogRing.LogAndAppend(fmt.Sprintf("%s: %s", InfoTimestamp, time.Now().Format(LogRingTimeFormat)))
	b.LogRing.LogAndAppend(fmt.Sprintf("%s: %s", InfoTotalTime, time.Since(b.timestamp)))
	b.LogRing.LogAndAppend(fmt.Sprintf("%s: %d", InfoTotalDirCopied, b.dirCounter))
	b.LogRing.LogAndAppend(fmt.Sprintf("%s: %d", InfoTotalFileCopied, b.fileCounter))
}

func (b *Backup) cleanBackups() error {
	if b.maxBackups <= 0 {
		return fmt.Errorf(ErrInvalidBackups)
	}

	totalBackups := len(b.sortedBackupDirs)
	totalToRemove := totalBackups - (b.maxBackups - 1)

	for i := 0; i < totalToRemove; i++ {
		folder := fmt.Sprintf("%s\\%s", b.dstPath, b.sortedBackupDirs[i].Format(TimeFormat))
		b.LogRing.LogAndAppend(fmt.Sprintf("%s: %s", InfoRemovingBackup, folder))
		err := os.RemoveAll(folder)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Backup) backupPost(backupPath, errorMessage string) error {
	b.LogRing.LogAndAppend(errorMessage)

	if exists(backupPath) {
		err := os.RemoveAll(backupPath)
		if err != nil {
			b.phase = stopped
			return err
		}
	}

	b.phase = stopped
	return fmt.Errorf(errorMessage)
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
