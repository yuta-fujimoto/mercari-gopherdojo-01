package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"strings"
)

func getError(err error, errFlg *bool) string {
	*errFlg = true
	errCmps := strings.Split(err.Error(), ": ")
	return "ft_cat:" + errCmps[0][strings.Index(errCmps[0], " "):]+ ": " + strings.ToUpper(errCmps[1][:1]) + errCmps[1][1:] + "\n"
}

func ft_write(rd io.Reader, wd io.Writer) error {
	rb := bufio.NewReader(rd)
	wb := bufio.NewWriter(wd)
	defer wb.Flush()

	for {
		n, err := rb.WriteTo(wb)
		if err != nil {
			return err
		}
		if n == 0 {
			break
		}
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
				io.WriteString(os.Stderr, getError(err, &errFlg))
				continue
			}
		}
		err = ft_write(file, os.Stdout)
		if err != nil {
			io.WriteString(os.Stderr, getError(err, &errFlg))
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
