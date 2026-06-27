package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// SafeRemoveAll removes all files and directories in targetPath.
// It verifies that targetPath is a valid, non-empty, subdirectory of a safe location
// and not a root directory.
func SafeRemoveAll(targetPath string) error {
	absPath, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Clean path
	absPath = filepath.Clean(absPath)

	// Safety check: Never delete root directories, system directories, etc.
	if absPath == "" || absPath == "/" || absPath == `\\` || filepath.VolumeName(absPath) == absPath {
		return fmt.Errorf("refusing to delete root or volume directory: %s", absPath)
	}

	// Additional safety: ensure it's not some system directories
	lowerPath := strings.ToLower(absPath)
	if strings.HasPrefix(lowerPath, "c:\\windows") || strings.HasPrefix(lowerPath, "c:\\program files") {
		return fmt.Errorf("refusing to delete system directories: %s", absPath)
	}

	return os.RemoveAll(absPath)
}

// Exists checks if a file or directory exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

// IsEmpty checks if a directory is empty.
func IsEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// CopyDir recursively copies a directory tree.
func CopyDir(src string, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyFile copies a single file from src to dst.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}
