package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/zkry/libra/analysis"
)

// Features:
// TODO: Directory summary
// TODO: Line Count
// TODO: File Specific analysis
// Flags:
// TODO: Non-programming language flag
// TODO: Ignore file flag

// Flags
var flagHelp bool
var flagAll bool

var wg sync.WaitGroup

// isHidden returns wheather or not the file is a hidden file by checking
// if its first character is a '.' TODO: Other systems hidden files???
func shouldSkipDir(n string) bool {
	// All flag overwrites all directory decisions
	if flagAll {
		return false
	}

	if strings.HasPrefix(n, ".") {
		return true
	} else if strings.ToLower(n) == "vendor" {
		return true
	}
	return false
}

// alalizeDir takes a directory, feeds all non-directory filetypes
// into chanel c, and makes a go routine for any found directories
func analizeDir(filepath string) {
	files, _ := ioutil.ReadDir(filepath)
	for _, f := range files {
		if f.IsDir() {
			if shouldSkipDir(f.Name()) {
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
	flag.BoolVar(&flagHelp, "h", false, "Show the help dialoge")
	flag.BoolVar(&flagAll, "a", false, "Traverse through all of the directories, including hidden ones")

	flag.Parse()

	if flagHelp {
		helpMessage()
		return
	}

	var filepath string
	if flag.NArg() == 0 {
		//		fmt.Println("Usage: libra <filepath>")
		filepath = "." // For debugging purposes only
		// return
	} else {
		filepath = flag.Arg(0)
	}

	wg.Add(1)
	go analizeDir(filepath)

	wg.Wait()

	analysis.DisplayReport()
}

func helpMessage() {
	fmt.Print(`Libra is a tool for analizing directories and their code composition.

Usages:

	libra [flags] [directory]

Flags:

	-h   Show the help dialoge 
	-a   Search through all directories (including hidden ones)
	`)
}
