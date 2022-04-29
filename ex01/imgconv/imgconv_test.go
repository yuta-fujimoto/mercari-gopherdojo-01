package imgconv_test

import (
	"convert/imgconv"
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

const (
	INPUT   = ".jpg"
	OUTPUT  = ".png"
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
			path:  "dir1",
			files: []string{"test1", "test2"},
		},
		{
			name:  "image in subdirectory",
			path:  "dir3",
			files: []string{"test1", "subdir/test2"},
		},
	}
	for _, td := range cases {
		td := td
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()
			filepath.WalkDir(filepath.Join(TESTDIR, td.path), func(path string, info fs.DirEntry, err error) error {
				if !info.IsDir() && filepath.Ext(info.Name()) == OUTPUT {
					os.Remove(path)
				}
				return nil
			})
			err := imgconv.ConvertImage(filepath.Join(TESTDIR, td.path), INPUT[1:], OUTPUT[1:])
			if err != nil {
				t.Fatal("ConvertImage failed by the error")
			}
			for _, file := range td.files {
				assertSameImages(filepath.Join(TESTDIR, td.path, file), t)
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
			filepath: "dir2/test1",
		},
	}
	for _, td := range cases {
		td := td
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()
			os.Remove(td.filepath + OUTPUT)
			err := imgconv.ConvertImage(filepath.Join(TESTDIR, td.filepath)+INPUT, INPUT[1:], OUTPUT[1:])
			if err != nil {
				t.Fatal("ConvertImage failed by the error")
			}
			assertSameImages(filepath.Join(TESTDIR, td.filepath), t)
		})
	}
}

func TestError(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		args [3]string
		msg  string
	}{
		{
			name: "no such direcory",
			args: [3]string{"nodir", INPUT[1:], OUTPUT[1:]},
			msg:  "nodir: no such file or directory",
		},
		{
			name: "invalid format 1",
			args: [3]string{"./", "noformat", OUTPUT[1:]},
			msg:  "noformat: invalid format",
		},
		{
			name: "invalid format 2",
			args: [3]string{"./", INPUT[1:], "noformat"},
			msg:  "noformat: invalid format",
		},
		{
			name: "contain text file",
			args: [3]string{filepath.Join(TESTDIR, "dir4"), INPUT[1:], OUTPUT[1:]},
			msg:  filepath.Join(TESTDIR, "dir4", "text.txt") + " is not a valid file",
		},
		{
			name: "input and output formats are same",
			args: [3]string{"./", INPUT[1:], INPUT[1:]},
			msg:  INPUT[1:] + ": input and output formats are same",
		},
	}
	for _, td := range cases {
		td := td
		t.Run(td.name, func(t *testing.T) {
			t.Parallel()
			got := imgconv.ConvertImage(td.args[0], td.args[1], td.args[2])
			want := "error: " + td.msg
			if got.Error() != want {
				t.Fatalf("got: [%s] want: [%s]", got, want)
			}
		})
	}
}

func assertSameImages(file string, t *testing.T) {
	t.Helper()
	infile, err := os.Open(file + INPUT)
	if err != nil {
		t.Fatal("failed to open " + file + OUTPUT)
	}
	defer infile.Close()
	outfile, err := os.Open(file + OUTPUT)
	if err != nil {
		t.Fatal("failed to open " + file + OUTPUT)
	}
	defer outfile.Close()
	inimg, _, err := image.Decode(infile)
	if err != nil {
		t.Fatal("failed to decode " + file + INPUT)
	}
	outimg, _, err := image.Decode(outfile)
	if err != nil {
		t.Fatal("failed to decode " + file + OUTPUT)
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
