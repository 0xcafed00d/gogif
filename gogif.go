package main

import (
	"fmt"
	//"image"
	"image/gif"
	"os"
)

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "usage: gogif <filename>")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	exitOnError(err)

	img, err := gif.DecodeAll(f)
	exitOnError(err)

	fmt.Println(len(img.Image))
}
