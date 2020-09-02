package log

import (
	"fmt"
	"log"
	"os"
)

var (
	debug bool
)

func init() {
	tmp := os.Getenv("DEBUG")
	if tmp == "1" {
		debug = true
	} else {
		debug = false
	}
}

func Error(s string, err error) {
	if debug == true {
		log.Print(s, err)
	}
}

func Info(s string, err error) {
	if debug == true {
		log.Println(s)
	}
}

func Print (s ...interface{}) {
	if debug == true {
		fmt.Println(s)
	}
}