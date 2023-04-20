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

func JsonToCsv(lines []string) ([]string, [][]string) {
	logs := make(map[string][]string)
	var result map[string]interface{}
	for _, line := range lines {
		json.Unmarshal([]byte(line), &result)
		for k, v := range result {
			logs[k] = append(logs[k], fmt.Sprintf("%v", v))
		}
	}

	columns := make([]string, 0, len(logs))
	rows := make([][]string, 0, len(logs))
	for k, v := range logs {
		columns = append(columns, k)
		rows = append(rows, v)
	}
	return columns, Transpose(rows)
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func GetCurrentDirectory() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}

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

func AddToLogFile(filename, field string, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		AppendToFile(filename, "Marshal err : "+err.Error())
		return
	}
	AppendToFile(filename, field+" : "+string(data))
}

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
