package lib

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	backupSuffix = ".bak"
)

func deleteSave00Bak() error {
	srcPath := buildDefaultSrcPath()
	log.Printf("deleting save00.bak folder")
	err := os.RemoveAll(fmt.Sprintf("%s%s", srcPath, backupSuffix))
	if err != nil {
		return err
	}

	return nil
}

func processSave00() error {
	srcPath, err := getSourcePath()
	if err != nil {
		if strings.Contains(err.Error(), "source path does not exist") {
			return nil
		}
		return err
	}

	err = deleteSave00Bak()
	if err != nil {
		return err
	}

	log.Printf("renaming save00 to save00.bak")
	err = os.Rename(srcPath, fmt.Sprintf("%s%s", srcPath, backupSuffix))
	if err != nil {
		return err
	}

	return nil
}

func restoreSave00(file, dstPath string, backupDirs []time.Time) error {
	// TODO: implement specified file restore
	_ = file

	// create destination directory
	log.Printf("creating save00 directory")
	srcPath, err := getSourcePath()
	if err != nil {
		if strings.Contains(err.Error(), "source path does not exist") {
			log.Printf("save00 directory does not exist")
		} else {
			return err
		}
	}

	// create directory
	err = os.MkdirAll(srcPath, os.ModePerm)
	if err != nil {
		return err
	}

	// recursively copy latest directory to destination
	latest := fmt.Sprintf("%s\\%s", dstPath, backupDirs[len(backupDirs)-1].Format(TimeFormat))
	log.Printf("copying latest backup %s to save00", latest)
	if err := copyDirectory(latest, srcPath); err != nil {
		log.Fatal(err)
	}

	log.Printf("successfully restored backup: %s", latest)
	phase = stopped

	// launch noita after successful restore
	// TODO: identify if called via cli
	err = LaunchNoita(true)
	if err != nil {
		log.Printf("failed to launch noita: %v", err)
	}

	return nil
}

func RestoreNoita(file string) error {
	// TODO: make this a channel/wait group as the logging is coming in incorrectly
	if !isNoitaRunning() {
		if phase == stopped {
			go func() {
				phase = started
				// get destination path
				dstPath, err := getDestinationPath()
				if err != nil {
					log.Printf("failed to get destination path: %v", err)
					return
				}

				// get sorted backup directories
				backupDirs, err := getBackupDirs(dstPath)
				if err != nil {
					log.Printf("failed to get backup dirs: %v", err)
					return
				}

				// protect against no backups
				if len(backupDirs) == 0 {
					log.Print("no backup dirs found, cannot restore")
					return
				}

				// process save00
				// 1. delete save00.bak
				// 2. rename save00 -> save00.bak
				if err := processSave00(); err != nil {
					log.Printf("error processing save00: %v", err)
					phase = stopped
					return
				}

				// restore specified (default latest) backup to destination
				if err := restoreSave00(file, dstPath, backupDirs); err != nil {
					log.Printf("error restoring backup file to save00: %v", err)
					phase = stopped
					return
				}

				return
			}()
		}
	} else {
		log.Print("noita.exe cannot be running during a restore")
		return nil
	}

	return nil
}
