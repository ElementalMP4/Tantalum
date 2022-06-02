package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type TantalumFileCollection struct {
	Files  []TantalumFile
	Couple TantalumCouple
}

type TantalumFile struct {
	Path string
	Info fs.FileInfo
}

type TantalumCouple struct {
	Left  string
	Right string
}

type TantalumConfig struct {
	Couples []TantalumCouple
	Output  bool
}

func print(text ...string) {
	if outputEnabled {
		fmt.Println(strings.Join(text, " "))
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var outputEnabled bool

func main() {
	file, err := ioutil.ReadFile("config.json")
	check(err)
	config := TantalumConfig{}
	err = json.Unmarshal([]byte(file), &config)
	check(err)

	outputEnabled = config.Output

	mirror(config.Couples)
}

func mirror(couples []TantalumCouple) {
	fileCollections := []TantalumFileCollection{}
	for _, couple := range couples {
		fmt.Println("Loading couple", couple.Left, ">", couple.Right)
		leftComplete := loadDirRecursive(couple.Left)
		collection := TantalumFileCollection{}
		collection.Couple = couple
		collection.Files = leftComplete
		fileCollections = append(fileCollections, collection)
	}
	for _, collection := range fileCollections {
		fmt.Println("Mirroring couple", collection.Couple.Left, ">", collection.Couple.Right)
		files, dirs := copyFiles(collection.Files, collection.Couple)
		fmt.Println("Copied", strconv.Itoa(files), "files and created", strconv.Itoa(dirs), "directories")
	}
}

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

func copyFiles(left []TantalumFile, couple TantalumCouple) (int, int) {
	filesCopied := 0
	dirsCreated := 0
	for _, file := range left {
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
				filesCopied++
				check(err)
			}
		}
	}
	return filesCopied, dirsCreated
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func copy(src, dst string) error {
	print("Copying from", src, "to", dst)
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
	print("Copied", strconv.FormatInt(nBytes, 10), "bytes from", src, "to", dst)
	return err
}
