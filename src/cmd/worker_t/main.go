package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var input string
	fmt.Println("23344")
	flag.StringVar(&input, "input", "fuck", "haha")
	env := os.Getenv("INPUT_ADDR")
	flag.Parse()
	fmt.Println(input)
	fmt.Println(env)
}
