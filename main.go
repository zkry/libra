package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
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

var wg sync.WaitGroup

// getExtensions takes a filename string, n, and returns the last part
// of the file
func getExtension(n string) (string, bool) {
	parts := strings.Split(n, ".")
	if len(parts) == 1 {
		return n, false
	}

	return parts[len(parts)-1], true
}

// fileCollector is the function that collects all of the found file
// information. File information comes from channel c, and is written
// to the files map
func fileCollector(files *map[string]int64, c chan os.FileInfo) {
	for f := range c {
		// Check if file has extension
		if ext, ok := getExtension(f.Name()); ok {
			if _, ok := ValidExts[ext]; !ok {
				continue
			}
			if _, ok := (*files)[ext]; ok {
				(*files)[ext] = (*files)[ext] + f.Size()
				analysis.Analize(f.Name())
			} else {
				(*files)[ext] = f.Size()
			}
		}
	}
}

func isHidden(n string) bool {
	return strings.HasPrefix(n, ".")
}

// alalizeDir takes a directory, feeds all non-directory filetypes
// into chanel c, and makes a go routine for any found directories
func analizeDir(c chan os.FileInfo, filepath string) {
	files, _ := ioutil.ReadDir(filepath)
	for _, f := range files {
		if f.IsDir() {
			if isHidden(f.Name()) {
				continue
			}
			wg.Add(1)
			go analizeDir(c, filepath+"/"+f.Name())
		} else {
			c <- f
		}
	}
	wg.Done()
}

// printBlocks displays a label with a bar to the right of it.
// Looks like the following:
//
// Go: [#######  ]
//
// With ext being the lefthand label, and bar of size n / max.
func printBlocks(ext string, n int, max int) {
	fmt.Printf("%5s: [", ext)
	for i := 0; i < max; i++ {
		if i < n {
			//fmt.Print("#")
			color.New(color.FgGreen).Fprintf(os.Stdout, "#")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Print("]\n")
}

// dispFiletypeStatistics takes the map of extensions to total size and
// displays a graph showing how much of each file the directory is
// composed of.
func dispFiletypeStatistics(fileStats map[string]int64) {
	barWidth := 50
	// Get max size
	var maxSize int64 = -1
	for _, val := range fileStats {
		if val > maxSize {
			maxSize = val
		}
	}

	// Display the file types
	for key, val := range fileStats {
		blocks := int(float64(val) / float64(maxSize) * float64(barWidth))
		printBlocks(key, blocks, barWidth)
	}
}

func main() {
	//filepath := os.Args[1]
	filepath := "."

	c := make(chan os.FileInfo, 10)
	fileStats := make(map[string]int64)

	wg.Add(1)
	go fileCollector(&fileStats, c)
	go analizeDir(c, filepath)

	wg.Wait()
	close(c)

	dispFiletypeStatistics(fileStats)
}
