package main

import (
	"flag"
	"fmt"
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
// TODO: Auto-ignore vendor file for go dirs
// Flags:
// TODO: Non-programming language flag
// TODO: Ignore file flag
// TODO: Add help flag
// TODO: Hidden files

var flagHelp bool

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

func helpMessage() {
	fmt.Print(`Libra is a tool for analizing directories and their code composition.

Usages:

	libra [flags] [directory]

Flags:

	-h   Show the help dialoge 
	`)
}

func main() {
	flag.BoolVar(&flagHelp, "h", false, "Show the help dialoge")
	flag.Parse()

	if flagHelp {
		helpMessage()
		return
	}

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
