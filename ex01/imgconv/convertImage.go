/*
Package convert enables to convert between JPEG, PNG and GIF. Also, it can make monochrome image(PGM) from
color image(JPEG, PNG and GIF).
*/
package imgconv

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

/*
Convert all image files in directory or filepath itself specified as string arg. inForm and outForm are I/O image format. If some sort of error occurs(failed to read directory, invalid format(txt, pdf, etc)), ConvertImage returns proper error and do nothing. Unnecessary formats jpg, png, pgm and gif are ignored if arg is specified as directory.
*/
func ConvertImage(arg string, inForm string, outForm string) error {
	params, err := initParams(arg, inForm, outForm)
	if err != nil {
		return err
	}
	input := make([]*os.File, params.Size)
	output := make([]*os.File, params.Size)

	for i := 0; i < params.Size; i++ {
		input[i], err = os.Open(params.Infile[i])
		if err != nil {
			closeAllFiles(input, output, i, i)
			return getError(err)
		}
		output[i], err = os.Create(params.Outfile[i])
		if err != nil {
			closeAllFiles(input, output, i+1, i)
			return getError(err)
		}
	}
	defer func() {
		closeAllFiles(input, output, params.Size, params.Size)
	}()
	for i := 0; i < params.Size; i++ {
		img, _, err := image.Decode(input[i])
		if err != nil {
			return err
		}
		switch params.Outform {
		case PNG:
			png.Encode(output[i], img)
		case JPEG:
			jpeg.Encode(output[i], img, &jpeg.Options{Quality: 100})
		}
	}
	return nil
}

func closeAllFiles(input []*os.File, output []*os.File, inCnt int, outCnt int) {
	var err error

	for i := 0; i < inCnt; i++ {
		err = input[i].Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
	for i := 0; i < outCnt; i++ {
		err = output[i].Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}
