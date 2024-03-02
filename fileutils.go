package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func loadDirRecursive(filePath string) []TantalumFile {
	completeFileList := []TantalumFile{}

	err := filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		file := TantalumFile{}
		file.Info = info
		file.Path = path
		completeFileList = append(completeFileList, file)
		return nil
	})

	check(err)

	return completeFileList
}

func copyOrUpdate(couple TantalumCouple, file TantalumFile, filesCopied int, dirsCreated int) (int, int) {
	rightSidePath := couple.Right + strings.ReplaceAll(file.Path, couple.Left, "")
	if file.Info.IsDir() {
		if !fileExists(rightSidePath) {
			os.Mkdir(rightSidePath, 0755)
			dirsCreated++
		}
	} else {
		if fileExists(rightSidePath) {
			rightSideFile, err := os.Stat(rightSidePath)
			check(err)
			if file.Info.ModTime().After(rightSideFile.ModTime()) {
				copy(file.Path, rightSidePath)
				filesCopied++
			}
		} else {
			err := copy(file.Path, rightSidePath)
			check(err)
			filesCopied++
		}
	}
	return filesCopied, dirsCreated
}

func copyFiles(left []TantalumFile, couple TantalumCouple) (int, int) {
	filesCopied := 0
	dirsCreated := 0

	for _, file := range left {
		filesCopied, dirsCreated = copyOrUpdate(couple, file, filesCopied, dirsCreated)
	}
	return filesCopied, dirsCreated
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func copy(src, dst string) error {
	info("Copying from", magenta(src), "to", cyan(dst))
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	defer ok("Copied", red(strconv.FormatInt(nBytes, 10)), "bytes from", magenta(src), "to", cyan(dst))
	return err
}
