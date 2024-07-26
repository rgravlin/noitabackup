package internal

import (
	"fmt"
	"github.com/spf13/viper"
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
	r.Backup.timestamp = time.Now()
	r.Backup.phase = started
	r.Backup.reportStart()

	// get sorted Backup directories
	r.Backup.sortedBackupDirs, err = getBackupDirs(r.Backup.dstPath, TimeFormat)
	if err != nil {
		return r.restorePost(fmt.Sprintf("%s: %v", ErrFailedGettingBackupDirs, err), false)
	}

	// protect against no Backups
	if len(r.Backup.sortedBackupDirs) == 0 {
		return r.restorePost(ErrNoBackupDirs, false)
	}

	// check the Backup directory for the specified Backup to restore
	if r.RestoreFile != "latest" {
		Backups := convertTimeSliceToStrings(r.Backup.sortedBackupDirs)
		if !slices.Contains(Backups, r.RestoreFile) {
			return r.restorePost(fmt.Sprintf(ErrBackupNotFound, r.RestoreFile), false)
		}
	}

	// process save00
	// 1. delete save00.bak
	// 2. rename save00 -> save00.bak
	if err := r.processSave00(); err != nil {
		return r.restorePost(fmt.Sprintf("%s: %v", ErrProcessingSave00, err), false)
	}

	// restore specified (default latest) Backup to destination
	if err := r.restoreSave00(); err != nil {
		return r.restorePost(fmt.Sprintf("%s: %v", ErrRestoringToSave00, err), true)
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

	// recursively copy source to destination
	latest := fmt.Sprintf("%s\\%s", r.Backup.dstPath, r.Backup.sortedBackupDirs[len(r.Backup.sortedBackupDirs)-1].Format(TimeFormat))
	r.Backup.LogRing.LogAndAppend(fmt.Sprintf(InfoCopyBackup, latest))
	if err := concurrentCopy(latest, r.Backup.srcPath, &r.Backup.dirCounter, &r.Backup.fileCounter, viper.GetInt("num-workers")); err != nil {
		r.Backup.LogRing.LogAndAppend(fmt.Sprintf("%s: %v", ErrCopyingToSave00, err))
		r.Backup.phase = stopped
		return err
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

func (r *Restore) restorePost(errorMessage string, cleanup bool) error {
	r.Backup.LogRing.LogAndAppend(errorMessage)

	if cleanup {
		// delete save00
		if exists(r.Backup.srcPath) {
			if err := deletePath(r.Backup.LogRing, r.Backup.srcPath); err != nil {
				r.Backup.phase = stopped
				return err
			}
		}

		// restore save00.bak due to failure
		if exists(fmt.Sprintf(r.Backup.srcPath, backupSuffix)) {
			r.Backup.LogRing.LogAndAppend(InfoRenameRestore)
			if err := os.Rename(fmt.Sprintf("%s%s", r.Backup.srcPath, backupSuffix), r.Backup.srcPath); err != nil {
				r.Backup.phase = stopped
				return err
			}
		}
	}

	r.Backup.phase = stopped
	return fmt.Errorf(errorMessage)
}

func convertTimeSliceToStrings(timeSlice []time.Time) []string {
	var stringSlice []string
	for _, t := range timeSlice {
		stringSlice = append(stringSlice, t.Format(TimeFormat))
	}
	return stringSlice
}
