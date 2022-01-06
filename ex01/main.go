package main

import (
	"convert/imgconv"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "error: invalid argument\n")
		return
	}
	err := imgconv.ConvertImage(args[0], "jpg", "png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}
