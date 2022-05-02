package imgconv_test

import (
	"convert/imgconv"
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	JPG   = "jpg"
	PNG  = "png"
	TESTDIR = "../testdata"
)

func TestDirectory(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		path  string
		files []string
	}{
		{
			name:  "normal",
			path:  "directory1",
			files: []string{"test1", "test2"},
		},
		{
			name:  "image in subdirectory",
			path:  "directory2",
			files: []string{"test1", "subdir/test2"},
		},
	}
	for _, td := range cases {
		td := td
		td.path = filepath.Join(TESTDIR, td.path)
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()
			filepath.WalkDir(td.path, func(path string, info fs.DirEntry, err error) error {
				if !info.IsDir() && filepath.Ext(info.Name())[1:] == PNG {
					os.Remove(path)
				}
				return nil
			})
			err := imgconv.ConvertImage(td.path, JPG, PNG)
			if err != nil {
				t.Fatal("ConvertImage failed by the error")
			}
			for _, file := range td.files {
				assertSameImages(filepath.Join(td.path, file), t)
			}
		})
	}
}

func TestFilePath(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		filepath string
	}{
		{
			name:     "normal",
			filepath: "filepath1/test1",
		},
	}
	for _, td := range cases {
		td := td
		td.filepath = filepath.Join(TESTDIR, td.filepath)
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()
			os.Remove(td.filepath + "." + PNG)

			err := imgconv.ConvertImage(td.filepath + "." + JPG, JPG, PNG)
			if err != nil {
				t.Fatal("ConvertImage failed by the error")
			}
			assertSameImages(td.filepath, t)
		})
	}
}

func TestError(t *testing.T) {
	permTestFilePath := filepath.Join(TESTDIR, "error2", "test1.jpg")
	os.Chmod(permTestFilePath, 0333)
	t.Parallel()
	cases := []struct {
		name string
		args [3]string
		msg  string
	}{
		{
			name: "no such direcory",
			args: [3]string{"nodir", JPG, PNG},
			msg:  "nodir: no such file or directory",
		},
		{
			name: "invalid format 1",
			args: [3]string{TESTDIR, "noformat", PNG},
			msg:  "noformat: invalid format",
		},
		{
			name: "invalid format 2",
			args: [3]string{TESTDIR, JPG, "noformat"},
			msg:  "noformat: invalid format",
		},
		{
			name: "contain text file",
			args: [3]string{filepath.Join(TESTDIR, "error1"), JPG, PNG},
			msg:  filepath.Join(TESTDIR, "error1", "text.txt") + " is not a valid file",
		},
		{
			name: "JPG and PNG formats are same",
			args: [3]string{TESTDIR, JPG, JPG},
			msg:  JPG + ": input and output formats are same",
		},
		{
			name: "no permission",
			args: [3]string{filepath.Join(TESTDIR, "error2", "test1") + "." + JPG, JPG, PNG},
			msg:  permTestFilePath + ": permission denied",
		},
	}
	for _, td := range cases {
		td := td
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()
			got := imgconv.ConvertImage(td.args[0], td.args[1], td.args[2])
			want := "error: " + td.msg
			if got == nil {
				t.Fatal("no error")
			}
			if got.Error() != want {
				t.Fatalf("got: [%s] want: [%s]", got, want)
			}
			if strings.Contains(td.name, "permission") {
				os.Chmod(td.args[0], 0755)
			}
		})
	}
}

func assertSameImages(file string, t *testing.T) {
	t.Helper()
	infile, err := os.Open(file + "." + JPG)
	if err != nil {
		t.Fatal("failed to open " + file + "." + PNG)
	}
	defer infile.Close()
	outfile, err := os.Open(file + "." + PNG)
	if err != nil {
		t.Fatal("failed to open " + file + "." + PNG)
	}
	defer outfile.Close()
	inimg, _, err := image.Decode(infile)
	if err != nil {
		t.Fatal("failed to decode " + file + "." + JPG)
	}
	outimg, _, err := image.Decode(outfile)
	if err != nil {
		t.Fatal("failed to decode " + file + "." + PNG)
	}
	if !inimg.Bounds().Eq(outimg.Bounds()) {
		t.Fatal("sizes of two image files differ :(")
	}
	
	// checking all pixels is time consuming, so it checks lower than 10000 pixels(100 * 100)
	bounds := inimg.Bounds()
	strideX := bounds.Dx() / 100
	strideY := bounds.Dy() / 100
	if strideX == 0 {
		strideX = 1
	}
	if strideY == 0 {
		strideY = 1
	}
	for y := 0; y < bounds.Dy(); y += strideY {
		for x := 0; x < bounds.Dx(); x += strideX {
			ir, ig, ib, ia := inimg.At(x, y).RGBA()
			or, og, ob, oa := outimg.At(x, y).RGBA()
			if ir != or || ig != og || ib != ob || ia != oa {
				t.Fatal("diff detected between two image files :(")
			}
		}
	}
}
