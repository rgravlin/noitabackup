package lib

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

const (
	backupSuffix = ".bak"
)

// deleteSave00Bak deletes the save00.bak folder by removing the directory at the source path with the backup suffix.
// It returns an error if the operation fails.
func deleteSave00Bak() error {
	srcPath := buildDefaultSrcPath()
	log.Printf("deleting save00.bak folder")
	err := os.RemoveAll(fmt.Sprintf("%s%s", srcPath, backupSuffix))
	if err != nil {
		return err
	}

	return nil
}

// processSave00 renames the save00 folder to save00.bak by performing the following steps:
// 1. Get the source path using GetSourcePath.
// 2. If the source path does not exist, return nil.
// 3. Delete the save00.bak folder by calling deleteSave00Bak.
// 4. If an error occurs during deleteSave00Bak, return the error.
// 5. Rename the save00 folder to save00.bak using os.Rename, using the source path and backupSuffix.
// 6. If an error occurs during the rename operation, return the error.
// 7. Return nil to indicate successful completion.
func processSave00() error {
	srcPath := viper.GetString("src")

	err := deleteSave00Bak()
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

// restoreSave00 restores the save00 directory by performing the following steps:
// 1. Creates the save00 directory at the source path.
// 2. Recursively copies the latest backup directory to the save00 directory.
// 3. Updates the phase variable to stopped.
// 4. Launches Noita if autoLaunchChecked is true.
// It returns an error if any of the operations fail.
func restoreSave00(file, dstPath string, backupDirs []time.Time, async bool) error {
	// TODO: implement specified file restore
	_ = file

	// create destination directory
	log.Printf("creating save00 directory")
	srcPath := viper.GetString("src")

	// create directory
	err := os.MkdirAll(srcPath, os.ModePerm)
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
	if autoLaunchChecked {
		err = LaunchNoita(async)
		if err != nil {
			log.Printf("failed to launch noita: %v", err)
		}
	}

	return nil
}

// RestoreNoita restores the save00 directory by performing the following steps:
//  1. Checks if Noita is running using isNoitaRunning.
//  2. If Noita is not running, it checks the phase variable and starts the restore process using a goroutine:
//     a. Sets the phase variable to started.
//     b. Gets the destination path using getDestinationPath.
//     c. Gets the sorted backup directories using getBackupDirs.
//     d. Checks if any backup directories exist and returns if none are found.
//     e. Calls processSave00 to rename the save00 folder to save00.bak.
//     f. Calls restoreSave00 to restore the specified backup or the latest backup to the save00 directory.
//     - Calls GetSourcePath to get the source path.
//     - Deletes the save00.bak folder using deleteSave00Bak.
//     - Renames the save00 folder to save00.bak.
//     - Creates the save00 directory at the source path.
//     - Copies the latest backup directory to the save00 directory.
//     - Sets the phase variable to stopped.
//     - Launches Noita if autoLaunchChecked is true.
//  3. If Noita is running, it returns an error and logs a message.
//
// It returns an error if any of the operations fail.
func RestoreNoita(file string, async bool) error {
	// TODO: make this a channel/wait group as the logging is coming in incorrectly
	if !isNoitaRunning() {
		if phase == stopped {
			go func() {
				phase = started
				// get destination path
				dstPath := viper.GetString("dst")

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
				if err := restoreSave00(file, dstPath, backupDirs, async); err != nil {
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
