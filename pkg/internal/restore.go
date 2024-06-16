package internal

import (
	"fmt"
	"log"
	"os"
	"slices"
	"time"
)

const (
	backupSuffix = ".bak"
)

type Restore struct {
	restoreFile string
	backup      *Backup
}

func NewRestore(restoreFile string, backup *Backup) *Restore {
	return &Restore{
		restoreFile: restoreFile,
		backup:      backup,
	}
}

func (r *Restore) RestoreNoita() {
	if !isNoitaRunning() {
		if r.backup.phase == stopped {
			if r.backup.async {
				go r.restoreNoita()
			} else {
				r.restoreNoita()
			}
		} else {
			log.Print("operation already in progress")
		}
	} else {
		log.Print("noita.exe cannot be running during a restore")
	}
}

func (r *Restore) restoreNoita() {
	var err error
	r.backup.phase = started

	// get sorted backup directories
	r.backup.sortedBackupDirs, err = getBackupDirs(r.backup.dstPath, TimeFormat)
	if err != nil {
		log.Printf("failed to get backup dirs: %v", err)
		r.backup.phase = stopped
		return
	}

	// protect against no backups
	if len(r.backup.sortedBackupDirs) == 0 {
		log.Print("no backup dirs found, cannot restore")
		r.backup.phase = stopped
		return
	}

	// check the backup directory for the specified backup to restore
	if r.restoreFile != "latest" {
		backups := convertTimeSliceToStrings(r.backup.sortedBackupDirs)
		if !slices.Contains(backups, r.restoreFile) {
			log.Printf("backup %s not found in backup directory", r.restoreFile)
			r.backup.phase = stopped
			return
		}
	}

	// process save00
	// 1. delete save00.bak
	// 2. rename save00 -> save00.bak
	if err := r.processSave00(); err != nil {
		log.Printf("error processing save00: %v", err)
		r.backup.phase = stopped
		return
	}

	// restore specified (default latest) backup to destination
	if err := r.restoreSave00(); err != nil {
		log.Printf("error restoring backup file to save00: %v", err)
		r.backup.phase = stopped
		return

	}

	r.backup.phase = stopped
}

func (r *Restore) restoreSave00() error {
	// create destination directory
	log.Printf("creating save00 directory")

	// create directory
	err := os.MkdirAll(r.backup.srcPath, os.ModePerm)
	if err != nil {
		return err
	}

	// recursively copy latest directory to destination
	latest := fmt.Sprintf("%s\\%s", r.backup.dstPath, r.backup.sortedBackupDirs[len(r.backup.sortedBackupDirs)-1].Format(TimeFormat))
	log.Printf("copying latest backup %s to save00", latest)
	var d, f int
	if err := copyDirectory(latest, r.backup.srcPath, &d, &f); err != nil {
		log.Fatal(err)
	}

	log.Printf("successfully restored backup: %s", latest)

	// launch noita after successful restore
	if r.backup.autoLaunchChecked {
		err = LaunchNoita(r.backup.async)
		if err != nil {
			log.Printf("failed to launch noita: %v", err)
		}
	}

	return nil
}

func (r *Restore) deleteSave00Bak() error {
	log.Printf("deleting save00.bak folder")
	err := os.RemoveAll(fmt.Sprintf("%s%s", r.backup.srcPath, backupSuffix))
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

	log.Printf("renaming save00 to save00.bak")
	err = os.Rename(r.backup.srcPath, fmt.Sprintf("%s%s", r.backup.srcPath, backupSuffix))
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
