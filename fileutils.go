package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
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

	check(err, false)

	return completeFileList
}

func copyOrUpdate(couple TantalumCouple, file TantalumFile, filesCopied int, dirsCreated int) (int, int) {
	rightSidePath := couple.Right + strings.ReplaceAll(file.Path, couple.Left, "")
	if file.Info.IsDir() {
		if !fileExists(rightSidePath) {
			os.Mkdir(rightSidePath, file.Info.Mode().Perm())
			dirsCreated++
		}
	} else {
		if fileExists(rightSidePath) {
			rightSideFile, err := os.Stat(rightSidePath)
			check(err, false)
			if file.Info.ModTime().After(rightSideFile.ModTime()) || couple.ForceUpdate {
				copy(file, rightSidePath)
				filesCopied++
			}
		} else {
			err := copy(file, rightSidePath)
			check(err, false)
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

func copy(srcFile TantalumFile, dst string) error {
	info("Copying from", magenta(srcFile.Path), "to", cyan(dst))
	sourceFileStat, err := os.Stat(srcFile.Path)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", srcFile.Path)
	}

	source, err := os.Open(srcFile.Path)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.OpenFile(dst, syscall.O_RDWR|syscall.O_CREAT|syscall.O_TRUNC, srcFile.Info.Mode())
	if err != nil {
		return err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	defer ok("Copied", red(strconv.FormatInt(nBytes, 10)), "bytes from", magenta(srcFile.Path), "to", cyan(dst))
	return err
}
