package main

import (
	"fmt"
	"strings"
)

const OK string = "  \x1b[1m\x1b[32mOK\x1b[0m  "
const FAIL string = " \x1b[1m\x1b[31mFAIL\x1b[0m "
const INFO string = " \x1b[1m\x1b[33mINFO\x1b[0m "

const RESET string = "\x1b[0m"
const MAGENTA string = "\x1b[35m"
const CYAN string = "\x1b[36m"
const BRIGHT_RED string = "\x1b[91m"

func print(level string, msg string) {
	if outputEnabled {
		fmt.Printf("%s %s\n", level, msg)
	}
}

func info(msg ...string) {
	print(INFO, strings.Join(msg, " "))
}

func ok(msg ...string) {
	print(OK, strings.Join(msg, " "))
}

func fail(msg ...string) {
	print(FAIL, strings.Join(msg, " "))
}

func magenta(in string) string {
	return MAGENTA + in + RESET
}

func cyan(in string) string {
	return CYAN + in + RESET
}

func red(in string) string {
	return BRIGHT_RED + in + RESET
}
