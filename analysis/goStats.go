package analysis

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
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
