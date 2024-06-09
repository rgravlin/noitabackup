package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	TimeFormat                = "2006-01-02-15-04-05"
	ConfigDefaultAppDataPath  = "..\\LocalLow\\Nolla_Games_Noita"
	ConfigDefaultSavePath     = "save00"
	ConfigDefaultDstPath      = "NoitaBackups"
	ConfigUserProfile         = "USERPROFILE"
	ConfigAppData             = "APPDATA"
	ConfigOverrideSrcPath     = "CONFIG_NOITA_SRC_PATH"
	ConfigOverrideDstPath     = "CONFIG_NOITA_DST_PATH"
	ConfigMaxNumBackupsToKeep = 100
	SteamExe                  = "C:\\Program Files (x86)\\Steam\\steam.exe"
	SteamNoitaFlags           = "steam://rungameid/881100"
	ExplorerExe               = "explorer"
)

var (
	dCounter = 0 // directory copy counter
	fCounter = 0 // file copy counter
)

func BackupNoita(async bool, maxBackups int) {
	if !isNoitaRunning() {
		if phase == stopped {
			phase = started
			// build timestamp
			t := time.Now()
			datePath := t.Format(TimeFormat)

			// build source path
			srcPath, err := getSourcePath()
			if err != nil {
				log.Printf("cannot get source path: %v", err)
				phase = stopped
				return
			}

			// build destination path
			dstPath, err := getDestinationPath()
			if err != nil {
				log.Printf("cannot get destination path: %v", err)
				phase = stopped
				return
			}

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
				if err := cleanBackups(sortedBackupDirs, backupPath, numberOfBackups); err != nil {
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
		} else {
			log.Printf("backup operation already in progress")
		}
	} else {
		log.Printf("noita.exe cannot be running to backup")
	}
}

func resetPhase() {
	phase = stopped
	dCounter = 0
	fCounter = 0
}

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
