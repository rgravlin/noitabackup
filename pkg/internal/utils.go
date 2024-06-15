package internal

import (
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// ConfigDefaultAppDataPath is the default path to the Noita application data folder.
// ConfigDefaultSavePath is the default folder name for saving Noita game data.
// ConfigDefaultDstPath is the default folder name for storing Noita backup files.
// ConfigUserProfile is the environment variable for the user profile path.
// ConfigAppData is the environment variable for the application data path.
// ConfigOverrideSrcPath is the environment variable for overriding the default source path for Noita backups.
// ConfigOverrideDstPath is the environment variable for overriding the default destination path for Noita backups.
const (
	ConfigDefaultAppDataPath = "..\\LocalLow\\Nolla_Games_Noita"
	ConfigDefaultSavePath    = "save00"
	ConfigDefaultDstPath     = "NoitaBackups"
	ConfigUserProfile        = "USERPROFILE"
	ConfigAppData            = "APPDATA"
	ConfigOverrideSrcPath    = "CONFIG_NOITA_SRC_PATH"
	ConfigOverrideDstPath    = "CONFIG_NOITA_DST_PATH"
)

func GetDefaultSourcePath() string {
	return buildDefaultSrcPath()
}

func GetDefaultDestinationPath() string {
	return buildDefaultDstPath()
}

func GetSourcePath(path string) (string, error) {
	// check for source path override
	srcPath := os.Getenv(ConfigOverrideSrcPath)
	if srcPath == "" {
		srcPath = path
	}

	// validate source path exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return srcPath, fmt.Errorf("source path does not exist: %s", srcPath)
	}

	return srcPath, nil
}

func GetDestinationPath(path string) (string, error) {
	// check for destination path override
	dstPath := os.Getenv(ConfigOverrideDstPath)
	if dstPath == "" {
		dstPath = path
	}

	// validate destination path exists
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		return "", fmt.Errorf("destination path does not exist: %s", dstPath)
	}

	return dstPath, nil
}

func copyDirectory(src, dst string, dirCounter, fileCounter *int) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		srcInfo, err := os.Stat(srcPath)
		if err != nil {
			return err
		}

		switch srcInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := createIfNotExists(dstPath, 0755); err != nil {
				return err
			}
			if err := copyDirectory(srcPath, dstPath, dirCounter, fileCounter); err != nil {
				return err
			}
			*dirCounter += 1
		default:
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
			*fileCounter += 1
		}

		fInfo, err := entry.Info()
		if err != nil {
			return err
		}

		isSymlink := fInfo.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(dstPath, fInfo.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Printf("error closing file: %v", err)
		}
	}(out)

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			log.Printf("error closing file: %v", err)
		}
	}(in)

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func createIfNotExists(dir string, perm os.FileMode) error {
	if exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: %s, error: %v", dir, err)
	}

	return nil
}

func buildDefaultSrcPath() string {
	path := os.Getenv(ConfigAppData)
	return fmt.Sprintf("%s\\%s\\%s", path, ConfigDefaultAppDataPath, ConfigDefaultSavePath)
}

func buildDefaultDstPath() string {
	path := os.Getenv(ConfigUserProfile)
	return fmt.Sprintf("%s\\%s", path, ConfigDefaultDstPath)
}

func LaunchExplorer() error {
	dstPath := viper.GetString("destination-path")

	// TODO: find out why explorer always returns an error code
	cmd := exec.Command(ExplorerExe, dstPath)
	_ = cmd.Run()
	return nil
}

func LaunchNoita(async bool) error {
	cmd := exec.Command(SteamExe, SteamNoitaFlags)

	if !isNoitaRunning() {
		if async {
			err := cmd.Start()
			if err != nil {
				log.Printf("error running steam: %v", err)
			}
		} else {
			err := cmd.Run()
			if err != nil {
				log.Printf("error running steam: %v", err)
			}
		}
	} else {
		log.Printf("noita.exe is already running")
	}

	return nil
}
