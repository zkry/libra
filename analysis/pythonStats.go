package analysis

import (
	"fmt"
	"os"
	"sync"
)

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
