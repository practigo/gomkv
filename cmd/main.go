package main

import (
	"fmt"
	"os"

	"github.com/practigo/gomkv"
)

func run(filename string) error {
	f, err := gomkv.Open(filename)
	if err != nil {
		return err
	}

	return gomkv.View(f)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing argument, provide an mkv/webm file!")
		os.Exit(1)
	}

	if err := run(os.Args[1]); err != nil {
		panic(err)
	}
}
