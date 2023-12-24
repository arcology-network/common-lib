// Package common provides common utility functions for file operations.

package common

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// IsPath checks if the given path ends with a forward slash.
// It returns true if the path is not empty and ends with a forward slash, false otherwise.
func IsPath(path string) bool {
	return len(path) > 0 && path[len(path)-1] == '/'
}

// GetParentPath returns the parent path of the given key.
// If the key is empty or the root ("/"), it returns the key itself.
// Otherwise, it returns the substring of the key up to the last occurrence of "/".
func GetParentPath(key string) string {
	if len(key) == 0 || key == "/" { //Root or empty
		return key
	}
	path := key[:strings.LastIndex(key[:len(key)-1], "/")+1]
	return path
}

// FileToLines reads the contents of the file with the given fileName and returns them as a slice of strings.
func FileToLines(fileName string) []string {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return lines
}

// FileExists checks if the file with the given filename exists.
// It returns true if the file exists, false otherwise.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if the directory with the given path exists.
// It returns true if the directory exists, false otherwise.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// GetCurrentDirectory returns the current working directory.
// It uses the filepath package to get the absolute path of the directory.
func GetCurrentDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}

// AppendToFile appends the given content to the file with the given filename.
// If the file does not exist, it creates a new file.
func AppendToFile(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content + "\n")
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return err
	}
	file.Sync()

	return nil
}

// AddToLogFile appends the given field and value to the log file with the given filename.
// It marshals the value to JSON format before appending.
func AddToLogFile(filename, field string, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		AppendToFile(filename, "Marshal err : "+err.Error())
		return
	}
	AppendToFile(filename, field+" : "+string(data))
}

// CopyFile copies the file from the source path to the destination path.
// It returns an error if the copy operation fails.
func CopyFile(src, dst string) error {
	if src == dst {
		return nil
	}

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}

	out, err := os.Create(dst)
	if err != nil {
		in.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	in.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}

	err = out.Sync()
	if err != nil {
		return fmt.Errorf("Sync error: %s", err)
	}

	si, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("Stat error: %s", err)
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return fmt.Errorf("Chmod error: %s", err)
	}

	return nil
}

// MoveFile moves the file from the source path to the destination path.
// It returns an error if the move operation fails.
func MoveFile(src, dst string) error {
	if src == dst {
		return nil
	}

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}

	out, err := os.Create(dst)
	if err != nil {
		in.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	in.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}

	err = out.Sync()
	if err != nil {
		return fmt.Errorf("Sync error: %s", err)
	}

	si, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("Stat error: %s", err)
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return fmt.Errorf("Chmod error: %s", err)
	}

	err = os.Remove(src)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}
