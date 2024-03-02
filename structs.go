package main

import "io/fs"

type TantalumFileCollection struct {
	Files  []TantalumFile
	Couple TantalumCouple
}

type TantalumFile struct {
	Path string
	Info fs.FileInfo
}

type TantalumCouple struct {
	Left        string
	Right       string
	ForceUpdate bool
}

type TantalumConfig struct {
	Couples []TantalumCouple
	Output  bool
}
