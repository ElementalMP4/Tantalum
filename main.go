package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
)

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

	var mirrorGroup sync.WaitGroup
	var i int = 0

	for _, collection := range fileCollections {
		mirrorGroup.Add(i)
		go startMirrorProcess(&mirrorGroup, collection)
		i++
	}
	mirrorGroup.Wait()
	fmt.Println("Process Completed")
}

func startMirrorProcess(workerGroup *sync.WaitGroup, collection TantalumFileCollection) {
	defer workerGroup.Done()
	fmt.Println("Mirroring couple", collection.Couple.Left, ">", collection.Couple.Right)
	files, dirs := copyFiles(collection.Files, collection.Couple)
	fmt.Println("Copied", strconv.Itoa(files), "files and created", strconv.Itoa(dirs), "directories")
}
