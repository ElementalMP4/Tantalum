package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func print(text ...string) {
	if outputEnabled {
		fmt.Println(strings.Join(text, " "))
	}
}

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}

var outputEnabled bool

func main() {
	file, err := os.ReadFile("config.json")
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
		startMirrorProcess(collection)
	}
	fmt.Println("Process Completed!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func startMirrorProcess(collection TantalumFileCollection) {
	fmt.Println("Mirroring couple", collection.Couple.Left, ">", collection.Couple.Right)
	files, dirs := copyFiles(collection.Files, collection.Couple)
	fmt.Println("Copied", strconv.Itoa(files), "files and created", strconv.Itoa(dirs), "directories")
}
