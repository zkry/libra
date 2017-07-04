package analysis

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/fatih/color"
)

type GoStat struct {
	funcCt      int
	interfaceCt int
	fileCt      int
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

	goStatMutex.Lock()
	goStat.funcCt = goStat.funcCt + funcCt
	goStat.interfaceCt = goStat.interfaceCt + interfaceCt
	goStat.fileCt += 1
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
}

type PythonStat struct {
	fileCt int
}

var pythonStat PythonStat
var pythonStatMutex sync.Mutex

func updatePythonExtStat(f os.FileInfo, path string) error {

	pythonStatMutex.Lock()
	pythonStat.fileCt += 1
	pythonStatMutex.Unlock()

	return nil
}

func displayPythonExtStat() {
	fmt.Printf("\n============= Python =============\n")
	fmt.Printf("Number of Files   : %d\n", pythonStat.fileCt)
}

var sizeStat = map[string]int64{}
var lineStat = map[string]int{}
var sizeStatMutex sync.Mutex

// updateSizeStat takes a file and updates the size map for that file
func updateSizeStat(f os.FileInfo, path string) {
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

		// Get the line count
		data, err := ioutil.ReadFile(path + "/" + f.Name())
		if err == nil {
			lineCt := bytes.Count(data, []byte("\n"))
			if val, ok := lineStat[ext]; ok {
				lineStat[ext] = val + lineCt
			} else {
				lineStat[ext] = lineCt
			}
		}
	}
	sizeStatMutex.Unlock()
}

// dispFiletypeStats takes the map of extensions to total size and
// displays a graph showing how much of each file the directory is
// composed of.
func dispSizeStat() {
	barWidth := 50 // in character width

	// Generate statistics for display
	sum := int64(0)      // in bytes
	maxSize := int64(-1) // in bytes
	maxLabelWidth := 0   // in character width
	for key, val := range sizeStat {
		sum += val
		if val > maxSize {
			maxSize = val
		}
		if len(ValidExts[key]) > maxLabelWidth {
			maxLabelWidth = len(ValidExts[key])
		}
	}

	keys := keysBySortedVal(sizeStat)

	// Display the file types
	for _, key := range keys {
		val := sizeStat[key]
		blocks := int(float64(val) / float64(maxSize) * float64(barWidth))
		printBlocks(ValidExts[key], blocks, barWidth, float64(val)/float64(sum)*100.0, maxLabelWidth, lineStat[key])
	}
}

// sizeStatMapKeysByVal returns the list of keys for sizesStat
// sorted in the order of its keys.
func keysBySortedVal(m map[string]int64) []string {
	reverseMap := map[int]string{}
	reverseKeys := make([]int, 0, 10)

	for key, val := range m {
		ival := int(val)
		if _, exists := reverseMap[ival]; !exists {
			reverseMap[ival] = key
			reverseKeys = append(reverseKeys, ival)
		} else {
			for {
				ival++
				if _, taken := reverseMap[ival]; !taken {
					reverseMap[ival] = key
					reverseKeys = append(reverseKeys, ival)
				}
			}
		}
	}

	sort.Ints(reverseKeys)
	sortedKeys := make([]string, 0, 10)

	for _, revKey := range reverseKeys {
		val, ok := reverseMap[revKey]
		if !ok {
			panic("Fatal error in sortedSizeMapByVal(): should always contain key")
		}
		sortedKeys = append(sortedKeys, val)
	}

	if len(sortedKeys) != len(m) {
		panic("Fatal error in sortedSizeMapByVal(): lists should be same size")
	}

	// Sort to make Highest to lowest
	for left, right := 0, len(sortedKeys)-1; left < right; left, right = left+1, right-1 {
		sortedKeys[left], sortedKeys[right] = sortedKeys[right], sortedKeys[left]
	}
	return sortedKeys
}

// printBlocks displays a label with a bar to the right of it.
// Looks like the following:
//
// Go: [#######  ]
//
// With ext being the lefthand label, and bar of size n / max.
func printBlocks(name string, n int, max int, percent float64, labelWidth int, lineCt int) {
	// Generate padding to left of label
	padding := labelWidth - len(name)
	for i := 0; i < padding; i++ {
		fmt.Printf(" ")
	}
	fmt.Printf("%s: [", name)
	// Print bar filling
	for i := 0; i < max; i++ {
		if i < n {
			//fmt.Print("#")
			color.New(color.FgGreen).Fprintf(os.Stdout, "#")
		} else {
			fmt.Print(" ")
		}
	}
	// Ending and percentage
	fmt.Printf("] %6.2f%% (%6d)\n", percent, lineCt)
}

// Analize takes a file and performs various tests for file
func Analize(f os.FileInfo, path string, wg *sync.WaitGroup) {
	// Filesize Statistic
	updateSizeStat(f, path)

	// Filetype Specific Analysis
	switch ext, _ := GetExtension(f.Name()); ext {
	case "go":
		updateGoExtStat(f, path)
	case "py":
		updatePythonExtStat(f, path)
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
	if _, ok := sizeStat["python"]; ok {
		displayPythonExtStat()
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
