package lib

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// getSourcePath retrieves the source path for the backup operation by checking for a source path override in the
// environment variables. If a source path override is not found, it builds the default source path. It then validates
// if the source path exists and returns it along with any error encountered.
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

// getDestinationPath retrieves the destination path for the backup operation by checking for a destination path override in the
// environment variables. If a destination path override is not found, it builds the default destination path. It then validates
// if the destination path exists and returns it along with any error encountered.
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

// copyDirectory recursively copies the contents of the source directory to the destination directory. It first reads
// the entries in the source directory and then loops through each entry. For each entry, it constructs the source and
// destination paths. If the entry is a directory, it creates the corresponding directory in the destination if it
// doesn't already exist, and then recursively calls copyDirectory on the subdirectory. If the entry is a file, it calls
// the copyFile function to copy the file from the source to the destination. After copying each file or directory, it
// sets the mode of the destination path to match the source path, excluding symbolic links. The function returns an
// error if any operation fails.
//
// Parameters:
// - src (string): The source directory path.
// - dst (string): The destination directory path.
//
// Returns:
// - error: An error if any operation fails, otherwise nil.
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

// copyFile copies the contents of the source file to the destination file. It creates or truncates the destination file,
// then opens the source file. It copies the contents of the source file to the destination file using the io.Copy function.
// After copying is complete, it closes both the source and destination files. If any error is encountered during this
// process, it returns the error. Otherwise, it returns nil.
//
// Parameters:
// - src (string): The path of the source file.
// - dst (string): The path of the destination file.
//
// Returns:
// - error: An error if any operation fails, otherwise nil.
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

// exists checks if the specified file or directory exists by checking the error returned by os.Stat().
// If the error indicates that the file or directory does not exist, it returns false.
// Otherwise, it returns true.
//
// Parameters:
// - filePath (string): The path of the file or directory to check.
//
// Returns:
// - (bool): true if the file or directory exists, false otherwise.
func exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

// createIfNotExists checks if the specified directory exists.
// If the directory does not exist, it creates the directory with the specified permissions.
// If the directory already exists, it returns nil.
//
// Parameters:
// - dir (string): The path of the directory to create.
// - perm (os.FileMode): The permissions to set for the directory.
//
// Returns:
// - error: An error if the directory creation fails, otherwise nil.
func createIfNotExists(dir string, perm os.FileMode) error {
	if exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: %s, error: %v", dir, err)
	}

	return nil
}

// buildDefaultSrcPath constructs and returns the default source path for the backup operation.
// The source path is built by retrieving the value of the environment variable "APPDATA"
// and concatenating it with the default application data path "../LocalLow/Nolla_Games_Noita"
// and the default save path "save00". The constructed source path is returned as a string.
func buildDefaultSrcPath() string {
	path := os.Getenv(ConfigAppData)
	return fmt.Sprintf("%s\\%s\\%s", path, ConfigDefaultAppDataPath, ConfigDefaultSavePath)
}

// buildDefaultDstPath builds the default destination path for the backup operation by concatenating the
// user's profile path with the value of ConfigDefaultDstPath constant.
// It then returns the formatted path as a string.
func buildDefaultDstPath() string {
	path := os.Getenv(ConfigUserProfile)
	return fmt.Sprintf("%s\\%s", path, ConfigDefaultDstPath)
}

// LaunchExplorer opens the file explorer at the specified destination path.
// It first retrieves the destination path for the backup operation by calling the getDestinationPath function.
// If an error occurs during the retrieval, it is returned.
// The function then creates a new exec.Command with the ExplorerExe constant and the destination path as arguments.
// The command is executed by calling cmd.Run(). Any error encountered is discarded.
// Finally, nil is returned to indicate that the operation completed successfully.
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

// LaunchNoita launches the Noita game using Steam.
// It takes a boolean parameter async, which determines if the game should be launched asynchronously.
// If async is true, the game is launched using cmd.Start(), otherwise it is launched using cmd.Run().
// If the Noita game is already running, it logs a message stating that the game is already running.
// It returns an error if any occurs during the execution of the function.
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
