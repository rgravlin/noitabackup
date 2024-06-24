package internal

import (
	"fmt"
	"os"
	"slices"
	"time"
)

const (
	backupSuffix = ".bak"
)

type Restore struct {
	RestoreFile string
	Backup      *Backup
}

func NewRestore(restoreFile string, backup *Backup) *Restore {
	return &Restore{
		RestoreFile: restoreFile,
		Backup:      backup,
	}
}

func (r *Restore) RestoreNoita() {
	if !isNoitaRunning() {
		if r.Backup.phase == stopped {
			if r.Backup.async {
				go func() { _ = r.restoreNoita() }()
			} else {
				_ = r.restoreNoita()
			}
		} else {
			r.Backup.LogRing.LogAndAppend(ErrOperationAlreadyInProgress)
		}
	} else {
		r.Backup.LogRing.LogAndAppend(fmt.Sprintf("%s %s", ErrNoitaRunning, ErrDuringRestore))
	}
}

func (r *Restore) restoreNoita() error {
	var err error
	t := time.Now()
	r.Backup.timestamp = t
	r.Backup.phase = started
	r.Backup.reportStart()

	// get sorted Backup directories
	r.Backup.sortedBackupDirs, err = getBackupDirs(r.Backup.dstPath, TimeFormat)
	if err != nil {
		r.Backup.LogRing.LogAndAppend(fmt.Sprintf("%s: %v", ErrFailedGettingBackupDirs, err))
		r.Backup.phase = stopped
		return err
	}

	// protect against no Backups
	if len(r.Backup.sortedBackupDirs) == 0 {
		r.Backup.LogRing.LogAndAppend(ErrNoBackupDirs)
		r.Backup.phase = stopped
		return err
	}

	// check the Backup directory for the specified Backup to restore
	if r.RestoreFile != "latest" {
		Backups := convertTimeSliceToStrings(r.Backup.sortedBackupDirs)
		if !slices.Contains(Backups, r.RestoreFile) {
			r.Backup.LogRing.LogAndAppend(fmt.Sprintf(ErrBackupNotFound, r.RestoreFile))
			r.Backup.phase = stopped
			return err
		}
	}

	// process save00
	// 1. delete save00.bak
	// 2. rename save00 -> save00.bak
	if err := r.processSave00(); err != nil {
		r.Backup.LogRing.LogAndAppend(fmt.Sprintf("%s: %v", ErrProcessingSave00, err))
		r.Backup.phase = stopped
		return err
	}

	// restore specified (default latest) Backup to destination
	if err := r.restoreSave00(); err != nil {
		r.Backup.LogRing.LogAndAppend(fmt.Sprintf("%s: %v", ErrRestoringToSave00, err))
		r.Backup.phase = stopped
		return err

	}

	r.Backup.reportStop()
	r.Backup.resetPhase()

	return nil
}

func (r *Restore) restoreSave00() error {
	// create destination directory
	r.Backup.LogRing.LogAndAppend(InfoCreatingSave00)

	// create directory
	err := os.MkdirAll(r.Backup.srcPath, os.ModePerm)
	if err != nil {
		return err
	}

	// recursively copy latest directory to destination
	latest := fmt.Sprintf("%s\\%s", r.Backup.dstPath, r.Backup.sortedBackupDirs[len(r.Backup.sortedBackupDirs)-1].Format(TimeFormat))
	r.Backup.LogRing.LogAndAppend(fmt.Sprintf(InfoCopyBackup, latest))
	if err := copyDirectory(latest, r.Backup.srcPath, &r.Backup.dirCounter, &r.Backup.fileCounter); err != nil {
		r.Backup.LogRing.LogAndAppend(fmt.Sprintf(ErrCopyingToSave00, latest, err))
	}

	r.Backup.LogRing.LogAndAppend(fmt.Sprintf("%s: %s", InfoSuccessfulRestore, latest))

	// launch noita after successful restore
	if r.Backup.autoLaunchChecked {
		err = LaunchNoita(r.Backup.async)
		if err != nil {
			r.Backup.LogRing.LogAndAppend(fmt.Sprintf("%s: %v", ErrFailedToLaunch, err))
		}
	}

	return nil
}

func (r *Restore) deleteSave00Bak() error {
	r.Backup.LogRing.LogAndAppend(InfoDeletingSave00Bak)
	err := os.RemoveAll(fmt.Sprintf("%s%s", r.Backup.srcPath, backupSuffix))
	if err != nil {
		return err
	}

	return nil
}

func (r *Restore) processSave00() error {
	err := r.deleteSave00Bak()
	if err != nil {
		return err
	}

	r.Backup.LogRing.LogAndAppend(InfoRename)
	err = os.Rename(r.Backup.srcPath, fmt.Sprintf("%s%s", r.Backup.srcPath, backupSuffix))
	if err != nil {
		return err
	}

	return nil
}

func convertTimeSliceToStrings(timeSlice []time.Time) []string {
	var stringSlice []string
	for _, t := range timeSlice {
		stringSlice = append(stringSlice, t.Format(TimeFormat))
	}
	return stringSlice
}
