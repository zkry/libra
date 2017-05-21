package main

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/zkry/libra/analysis"
)

// Features:
// TODO: Organize output based off of length
// TODO: Directory summary
// TODO: Line Count
// TODO: File Specific analysis
// Flags:
// TODO: Non-programming language flag
// TODO: Ignore file flag
// TODO: Add help flag

type FileWithContext struct {
	fileInfo os.FileInfo
	path     string
}

var wg sync.WaitGroup

// fileCollector is the function that collects all of the found file
// information. File information comes from channel c, and is written
// to the files map
func fileCollector(files *map[string]int64, c chan FileWithContext) {
	for f := range c {
		analysis.Analize(f.fileInfo, f.path)
		wg.Done()
	}
}

func isHidden(n string) bool {
	return strings.HasPrefix(n, ".")
}

// alalizeDir takes a directory, feeds all non-directory filetypes
// into chanel c, and makes a go routine for any found directories
func analizeDir(c chan FileWithContext, filepath string) {
	files, _ := ioutil.ReadDir(filepath)
	for _, f := range files {
		if f.IsDir() {
			if isHidden(f.Name()) {
				continue
			}
			wg.Add(1)
			go analizeDir(c, filepath+"/"+f.Name())
		} else {
			wg.Add(1)

			c <- FileWithContext{fileInfo: f, path: filepath}
		}
	}
	wg.Done()
}

func main() {
	var filepath string

	if len(os.Args) == 1 {
		//		fmt.Println("Usage: libra <filepath>")
		filepath = "." // For debugging purposes only
		// return
	} else {
		filepath = os.Args[1]
	}

	c := make(chan FileWithContext, 10)
	fileStats := make(map[string]int64)

	wg.Add(1)
	go fileCollector(&fileStats, c)
	go analizeDir(c, filepath)

	wg.Wait()
	close(c)

	analysis.DisplayReport()
}
