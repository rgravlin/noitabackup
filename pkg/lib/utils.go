package lib

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func getSourcePath() (string, error) {
	// check for source path override
	srcPath := os.Getenv(ConfigOverrideSrcPath)
	if srcPath == "" {
		srcPath = buildDefaultSrcPath()
	}

	// validate source path exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return srcPath, fmt.Errorf("source path does not exist: %s", srcPath)
	}

	return srcPath, nil
}

func getDestinationPath() (string, error) {
	// check for destination path override
	dstPath := os.Getenv(ConfigOverrideDstPath)
	if dstPath == "" {
		dstPath = buildDefaultDstPath()
	}

	// validate destination path exists
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		return "", fmt.Errorf("destination path does not exist: %s", dstPath)
	}

	return dstPath, nil
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
	dstPath, err := getDestinationPath()
	if err != nil {
		return err
	}

	// TODO: find out why explorer always returns an error code
	cmd := exec.Command(ExplorerExe, dstPath)
	_ = cmd.Run()
	return nil
}

func LaunchNoita() error {
	if !isNoitaRunning() {
		cmd := exec.Command(SteamExe, SteamNoitaFlags)

		err := cmd.Run()
		if err != nil {
			return err
		}
	} else {
		log.Printf("noita.exe is already running")
	}

	return nil
}
