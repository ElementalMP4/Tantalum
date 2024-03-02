package main

import (
	"encoding/json"
	"os"
	"strconv"
)

func check(err error) {
	if err != nil {
		fail(err.Error())
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
		info("Loading differences", magenta(couple.Left), ">", cyan(couple.Right))
		leftComplete := loadDirRecursive(couple.Left)
		collection := TantalumFileCollection{}
		collection.Couple = couple
		collection.Files = leftComplete
		fileCollections = append(fileCollections, collection)
	}

	for _, collection := range fileCollections {
		startMirrorProcess(collection)
	}
	ok("Process Completed!")
}

func startMirrorProcess(collection TantalumFileCollection) {
	info("Mirroring couple", magenta(collection.Couple.Left), ">", cyan(collection.Couple.Right))
	files, dirs := copyFiles(collection.Files, collection.Couple)
	ok("Copied", red(strconv.Itoa(files)), "files and created", red(strconv.Itoa(dirs)), "directories")
}
