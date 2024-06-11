package lib

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	TimeFormat                = "2006-01-02-15-04-05"
	ConfigMaxNumBackupsToKeep = 100
	SteamExe                  = "C:\\Program Files (x86)\\Steam\\steam.exe"
	SteamNoitaFlags           = "steam://rungameid/881100"
	ExplorerExe               = "explorer"
)

var (
	dCounter = 0 // directory copy counter
	fCounter = 0 // file copy counter
)

// BackupNoita performs the backup operation for the Noita game.
// It checks if Noita is running and if the backup operation is already in progress.
// If Noita is not running and the backup operation is not already in progress, it proceeds with the backup process.
// The function builds the timestamp, source path, and destination path.
// It then retrieves the number of existing backups and checks if it exceeds the maximum backup threshold.
// If the number of backups exceeds the maximum threshold, it sorts the backup directories from oldest to newest.
// It then cleans up the oldest backup directories to make room for the new backup.
// After that, it creates the destination path and copies the source directory contents to the destination.
// Lastly, it logs the backup statistics and resets the phase.
// If the auto-launch feature is enabled, it launches Noita after a successful backup.
//
// Parameters:
//   - async (bool): Determines whether to launch Noita asynchronously after the backup (true) or synchronously (false).
//   - maxBackups (int): The maximum number of backups to keep. If the number of backups exceeds this limit,
//     the oldest backups will be deleted to make room for new backups.
func BackupNoita(async bool, maxBackups int) {
	if !isNoitaRunning() {
		if phase == stopped {
			go func() {
				phase = started
				// build timestamp
				t := time.Now()
				datePath := t.Format(TimeFormat)

				// build source path
				srcPath := viper.GetString("source-path")

				// build destination path
				dstPath := viper.GetString("destination-path")

				// mutate destination with timestamp
				backupPath := dstPath
				dstPath = fmt.Sprintf("%s\\%s", dstPath, datePath)

				// report start
				log.Printf("timestamp: %s\n", datePath)
				log.Printf("source: %s\n", srcPath)
				log.Printf("destination: %s\n", dstPath)

				numberOfBackups, err := getNumBackups(backupPath)
				if err != nil {
					log.Printf("error getting backups: %v", err)
					phase = stopped
					return
				} else {
					log.Printf("number of backups: %d", numberOfBackups)
				}

				// protect against invalid maxBackups
				// cannot breach maximum (100)
				// cannot breach minimum (1)
				if maxBackups > ConfigMaxNumBackupsToKeep || maxBackups <= 0 {
					maxBackups = ConfigMaxNumBackupsToKeep
				}

				if numberOfBackups >= maxBackups {
					log.Printf("maximum backup threshold reached")

					// get and sort backup directories
					// oldest are first in the sorted slice
					sortedBackupDirs, err := getBackupDirs(backupPath)
					if err != nil {
						log.Printf("error getting backups: %v", err)
						phase = stopped
						return
					}

					// clean backup directories - 1 to make room for this backup
					if err := cleanBackups(sortedBackupDirs, backupPath, maxBackups-1); err != nil {
						log.Printf("failure deleting backups: %v", err)
						phase = stopped
						return
					}
				}

				// create destination path
				if err := createIfNotExists(dstPath, 0755); err != nil {
					log.Printf("cannot create destination path: %v", err)
					phase = stopped
					return
				}

				// recursively copy source to destination
				if err := copyDirectory(srcPath, dstPath); err != nil {
					log.Printf("cannot copy source to destination: %v", err)
					phase = stopped
					return
				}

				// return stats
				log.Printf("timestamp: %s\n", time.Now().Format(TimeFormat))
				log.Printf("total time: %s\n", time.Since(t))
				log.Printf("total dirs copied: %d\n", dCounter)
				log.Printf("total files copied: %d\n", fCounter)

				// reset phase
				resetPhase()

				// launch noita automatically after successful backup
				if autoLaunchChecked {
					err = LaunchNoita(async)
					if err != nil {
						log.Printf("failed to launch noita: %v", err)
					}
				}

			}()
		} else {
			log.Printf("backup operation already in progress")
		}
	} else {
		log.Printf("noita.exe cannot be running to backup")
	}
}

// resetPhase resets the phase, dCounter, and fCounter variables to their initial values.
func resetPhase() {
	phase = stopped
	dCounter = 0
	fCounter = 0
}

// cleanBackups removes the oldest backup directories to make room for new backups.
// It receives a list of backup directories, the backup path, and the number of backups to keep.
// It calculates the number of directories to remove based on the difference between the total number of backups and the number of backups to keep.
// Then it iterates through the backup directories from oldest to newest and removes the oldest directories from the file system.
// The function returns an error if any deletion operation fails.
//
// Parameters:
//   - backupDirs ([]time.Time): A list of backup directories represented by timestamps.
//   - backupPath (string): The path to the backup directory.
//   - numToKeep (int): The maximum number of backups to keep. If the number of backups exceeds this limit,
//     the oldest backups will be deleted to make room for new backups.
func cleanBackups(backupDirs []time.Time, backupPath string, numToKeep int) error {
	totalBackups := len(backupDirs)
	totalToRemove := totalBackups - numToKeep

	for i := 0; i < totalToRemove; i++ {
		folder := fmt.Sprintf("%s\\%s", backupPath, backupDirs[i].Format(TimeFormat))
		log.Printf("removing backup folder: %s", folder)
		err := os.RemoveAll(folder)
		if err != nil {
			return err
		}
	}

	return nil
}

// getBackupDirs retrieves the list of backup directories in the specified backupPath.
// If backupPath does not exist, an empty slice is returned.
// If an error occurs while reading the backupPath directory, the error is returned.
// For each directory entry in backupPath, the function checks if it is a directory.
// If it is a directory, it parses the directory name as a time value using the TimeFormat constant.
// The parsed time value is appended to the backupDirs slice.
// Finally, the backupDirs slice is sorted in ascending order based on the time values and returned.
//
// Parameters:
// - backupPath (string): The path to the backup directory.
//
// Returns:
// - ([]time.Time): The list of backup directories sorted by ascending time values.
// - (error): An error that occurred during the process, or nil if successful.
func getBackupDirs(backupPath string) ([]time.Time, error) {
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
				nameDate, err := time.Parse(TimeFormat, name[len(name)-1])
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

// getNumBackups retrieves the number of existing backups in the specified backupPath.
// It calls the getBackupDirs function to retrieve the list of backup directories in the backupPath.
// The number of backup directories is equal to the number of existing backups.
//
// Parameters:
// - backupPath (string): The path to the backup directory.
//
// Returns:
// - (int): The number of existing backups.
// - (error): An error that occurred during the process, or nil if successful.
func getNumBackups(backupPath string) (int, error) {
	numBackups := 0
	backupDirs, err := getBackupDirs(backupPath)
	if err != nil {
		return numBackups, err
	}

	return len(backupDirs), nil
}

type ByDate []time.Time

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Before(a[j]) }
