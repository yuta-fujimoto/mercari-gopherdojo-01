package imgconv

import (
	"fmt"
	"strings"
)

// valid image format(PGM is for output only)
const (
	JPEG = ".jpg"
	PNG  = ".png"
)

// ConvertImage at first searches directory and specifies formats to store them into Params
type Params struct {
	Infile  []string
	Outfile []string
	Inform  string
	Outform string
	Size    int
}

func getError(err error) error {
	newErrorMsg := err.Error()
	if strings.Contains(err.Error(), " ") {
		newErrorMsg = "error:" + err.Error()[strings.Index(err.Error(), " "):]
	}
	return fmt.Errorf("%s", newErrorMsg)
}
