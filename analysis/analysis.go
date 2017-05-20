package analysis

import (
	"bytes"
	"io"
	"os"
)

type GoStat struct {
	funcCt      int
	interfaceCt int
	fileCt      int
}

func Analize(filename string) {

}

// LineCount counts the number of lines in a file given a filename.
func LineCount(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := file.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}
