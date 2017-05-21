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

var wg sync.WaitGroup

// isHidden returns wheather or not the file is a hidden file by checking
// if its first character is a '.'
func isHidden(n string) bool {
	return strings.HasPrefix(n, ".")
}

// alalizeDir takes a directory, feeds all non-directory filetypes
// into chanel c, and makes a go routine for any found directories
func analizeDir(filepath string) {
	files, _ := ioutil.ReadDir(filepath)
	for _, f := range files {
		if f.IsDir() {
			if isHidden(f.Name()) {
				continue
			}
			wg.Add(1)
			go analizeDir(filepath + "/" + f.Name())
		} else {
			wg.Add(1)
			go analysis.Analize(f, filepath, &wg)
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

	wg.Add(1)
	go analizeDir(filepath)

	wg.Wait()

	analysis.DisplayReport()
}
