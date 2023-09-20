//go:build windows

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	TimeFormat               = "2006-01-02-15-04-05"
	ConfigDefaultAppDataPath = "..\\LocalLow\\Nolla_Games_Noita"
	ConfigDefaultSavePath    = "save00"
	ConfigDefaultDstPath     = "NoitaBackups"
	ConfigUserProfile        = "USERPROFILE"
	ConfigAppData            = "APPDATA"
	ConfigOverrideSrcPath    = "CONFIG_NOITA_SRC_PATH"
	ConfigOverrideDstPath    = "CONFIG_NOITA_DST_PATH"
)

var (
	dCounter = 0
	fCounter = 0
)

func main() {
	t := time.Now()
	datePath := t.Format(TimeFormat)

	// check for source path override
	srcPath := os.Getenv(ConfigOverrideSrcPath)
	if srcPath == "" {
		srcPath = buildDefaultSrcPath()
	}

	// check for destination path override
	dstPath := os.Getenv(ConfigOverrideDstPath)
	if dstPath == "" {
		dstPath = buildDefaultDstPath()
	}

	dstPath = fmt.Sprintf("%s\\%s", dstPath, datePath)

	if err := createIfNotExists(dstPath, 0755); err != nil {
		log.Fatal(err)
	}

	if err := copyDirectory(srcPath, dstPath); err != nil {
		log.Fatal(err)
	}

	log.Printf("%s: %s\n", "Destination", dstPath)
	log.Printf("%s: %d\n", "Total dirs copied", dCounter)
	log.Printf("%s: %d\n", "Total files copied", fCounter)
}

func copyDirectory(src, dst string) error {
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
			if err := copyDirectory(srcPath, dstPath); err != nil {
				return err
			}
			dCounter += 1
		case os.ModeSymlink:
			if err := copySymLink(srcPath, dstPath); err != nil {
				return err
			}
		default:
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
			fCounter += 1
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
	defer out.Close()

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

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

func copySymLink(src, dst string) error {
	link, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(link, dst)
}
