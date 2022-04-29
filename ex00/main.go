package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"strings"
)

func getError(err error, fn string) string {
	errWords := strings.Split(err.Error(), ": ")
	errMsg := errWords[len(errWords) - 1]
	errMsg = strings.ToUpper(errMsg[:1]) + errMsg[1:]
	return "ft_cat: " + fn + ": " + errMsg + "\n"
}


// https://pkg.go.dev/io#Reader
// I/O abstraction with io package
func ft_write_file_content(rd io.Reader, wd io.Writer) error {
	rb := bufio.NewReader(rd)
	wb := bufio.NewWriter(wd)
	defer wb.Flush()

	_, err := rb.WriteTo(wb)
	if err != nil {
		return err
	}
	return nil
}

func ft_cat(args []string) (errFlg bool) {
	var file *os.File
	var err error
	errFlg = false

	for _, arg := range args {
		if arg == "-" {
			file = os.Stdin
		} else {
			file, err = os.Open(arg)
			if err != nil {
				errFlg = true
				io.WriteString(os.Stderr, getError(err, arg))
				continue
			}
		}
		err = ft_write_file_content(file, os.Stdout)
		if err != nil {
			errFlg = true
			io.WriteString(os.Stderr, getError(err, arg))
		}
	}
	return
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}
	errFlg := ft_cat(args)
	if errFlg {
		os.Exit(1)
	}
	os.Exit(0)
}
