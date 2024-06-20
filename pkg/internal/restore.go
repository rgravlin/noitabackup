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
				go r.restoreNoita()
			} else {
				r.restoreNoita()
			}
		} else {
			r.Backup.LogRing.LogAndAppend("operation already in progress")
		}
	} else {
		r.Backup.LogRing.LogAndAppend("noita.exe cannot be running during a restore")
	}
}

func (r *Restore) restoreNoita() {
	var err error
	r.Backup.phase = started

	// get sorted Backup directories
	r.Backup.sortedBackupDirs, err = getBackupDirs(r.Backup.dstPath, TimeFormat)
	if err != nil {
		r.Backup.LogRing.LogAndAppend(fmt.Sprintf("failed to get backup dirs: %v", err))
		r.Backup.phase = stopped
		return
	}

	// protect against no Backups
	if len(r.Backup.sortedBackupDirs) == 0 {
		r.Backup.LogRing.LogAndAppend("no backup dirs found, cannot restore")
		r.Backup.phase = stopped
		return
	}

	// check the Backup directory for the specified Backup to restore
	if r.RestoreFile != "latest" {
		Backups := convertTimeSliceToStrings(r.Backup.sortedBackupDirs)
		if !slices.Contains(Backups, r.RestoreFile) {
			r.Backup.LogRing.LogAndAppend(fmt.Sprintf("backup %s not found in backup directory", r.RestoreFile))
			r.Backup.phase = stopped
			return
		}
	}

	// process save00
	// 1. delete save00.bak
	// 2. rename save00 -> save00.bak
	if err := r.processSave00(); err != nil {
		r.Backup.LogRing.LogAndAppend(fmt.Sprintf("error processing save00: %v", err))
		r.Backup.phase = stopped
		return
	}

	// restore specified (default latest) Backup to destination
	if err := r.restoreSave00(); err != nil {
		r.Backup.LogRing.LogAndAppend(fmt.Sprintf("error restoring backup file to save00: %v", err))
		r.Backup.phase = stopped
		return

	}

	r.Backup.phase = stopped
}

func (r *Restore) restoreSave00() error {
	// create destination directory
	r.Backup.LogRing.LogAndAppend("creating save00 directory")

	// create directory
	err := os.MkdirAll(r.Backup.srcPath, os.ModePerm)
	if err != nil {
		return err
	}

	// recursively copy latest directory to destination
	latest := fmt.Sprintf("%s\\%s", r.Backup.dstPath, r.Backup.sortedBackupDirs[len(r.Backup.sortedBackupDirs)-1].Format(TimeFormat))
	r.Backup.LogRing.LogAndAppend(fmt.Sprintf("copying latest backup %s to save00", latest))
	var d, f int
	if err := copyDirectory(latest, r.Backup.srcPath, &d, &f); err != nil {
		r.Backup.LogRing.LogAndAppend(fmt.Sprintf("error copying latest backup %s to save00: %v", latest, err))
	}

	r.Backup.LogRing.LogAndAppend(fmt.Sprintf("successfully restored backup: %s", latest))

	// launch noita after successful restore
	if r.Backup.autoLaunchChecked {
		err = LaunchNoita(r.Backup.async)
		if err != nil {
			r.Backup.LogRing.LogAndAppend(fmt.Sprintf("failed to launch noita: %v", err))
		}
	}

	return nil
}

func (r *Restore) deleteSave00Bak() error {
	r.Backup.LogRing.LogAndAppend("deleting save00.bak folder")
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

	r.Backup.LogRing.LogAndAppend("renaming save00 to save00.bak")
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
