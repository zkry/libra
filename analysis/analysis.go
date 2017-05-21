package analysis

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
)

type GoStat struct {
	funcCt      int
	interfaceCt int
	fileCt      int
	lineCt      int
}

var goStat GoStat
var goStatMutex sync.Mutex

func updateGoExtStat(f os.FileInfo, path string) error {
	data, err := ioutil.ReadFile(path + "/" + f.Name())
	if err != nil {
		return err
	}

	funcCt := bytes.Count(data, []byte("func"))
	interfaceCt := bytes.Count(data, []byte("interface"))
	lineCt := bytes.Count(data, []byte("\n"))

	goStatMutex.Lock()
	goStat.funcCt = goStat.funcCt + funcCt
	goStat.interfaceCt = goStat.interfaceCt + interfaceCt
	goStat.fileCt += 1
	goStat.lineCt += lineCt
	goStatMutex.Unlock()

	return nil
}

// displayGoExtStat takes the goStat variable and prints it to the screen
// in a pretty way
func displayGoExtStat() {
	fmt.Printf("\n============= Go =============\n")
	fmt.Printf("Number of Files      : %d\n", goStat.fileCt)
	fmt.Printf("Number of Functions  : %d\n", goStat.funcCt)
	fmt.Printf("Number of Interfaces : %d\n", goStat.interfaceCt)
	fmt.Printf("Nubmer of Lines      : %d\n", goStat.lineCt)
}

var sizeStat = map[string]int64{}
var sizeStatMutex sync.Mutex

// updateSizeStat takes a file and updates the size map for that file
func updateSizeStat(f os.FileInfo) {
	sizeStatMutex.Lock()
	// Check if a . extension exists
	ext, extOK := GetExtension(f.Name())
	_, typeOK := ValidExts[ext]

	if extOK && typeOK {
		if val, ok := sizeStat[ext]; ok {
			sizeStat[ext] = val + f.Size()
		} else {
			sizeStat[ext] = f.Size()
		}
	}
	sizeStatMutex.Unlock()
}

// dispFiletypeStats takes the map of extensions to total size and
// displays a graph showing how much of each file the directory is
// composed of.
func dispSizeStat() {
	barWidth := 50
	sum := int64(0)
	// Get max size
	var maxSize int64 = -1
	for _, val := range sizeStat {
		fmt.Println("DEBUG: size - ", val)
		sum += val
		if val > maxSize {
			maxSize = val
		}
	}

	// Display the file types
	for key, val := range sizeStat {
		blocks := int(float64(val) / float64(maxSize) * float64(barWidth))
		printBlocks(ValidExts[key], blocks, barWidth, float64(val)/float64(sum)*100.0)
	}
}

// printBlocks displays a label with a bar to the right of it.
// Looks like the following:
//
// Go: [#######  ]
//
// With ext being the lefthand label, and bar of size n / max.
func printBlocks(ext string, n int, max int, percent float64) {
	fmt.Printf("%5s: [", ext)
	for i := 0; i < max; i++ {
		if i < n {
			//fmt.Print("#")
			color.New(color.FgGreen).Fprintf(os.Stdout, "#")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Printf("] %6.2f%%\n", percent)
}

// Analize takes a file and performs various tests for file
func Analize(f os.FileInfo, path string, wg *sync.WaitGroup) {
	// Filesize Statistic
	updateSizeStat(f)

	// Filetype Specific Analysis
	switch ext, _ := GetExtension(f.Name()); ext {
	case "go":
		updateGoExtStat(f, path)
	}

	// TODO: Git Analysis

	wg.Done()
}

// DisplayReport writes all of the collected data to the screen. Called at the
// end of the program when all files are analized.
func DisplayReport() {
	// Size summary report
	dispSizeStat()

	// File Specific Reports
	if _, ok := sizeStat["go"]; ok {
		displayGoExtStat()
	}
	// TODO: Git report
}

// GetExtensions takes a filename string, n, and returns the last part
// of the file
func GetExtension(n string) (string, bool) {
	parts := strings.Split(n, ".")
	if len(parts) == 1 {
		return n, false
	}

	return parts[len(parts)-1], true
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
