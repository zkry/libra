package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

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

// wg is the wait group for managing threads that read files. Upon finding a directory, a new go routine is added and the wait groups cout is increased.
var wg sync.WaitGroup

// shouldSkipDir returns wheather or not the file is a hidden file by checking
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
	rand.Seed(time.Now().UnixNano())
	flag.BoolVar(&flagHelp, "h", false, "Show the help dialoge")
	flag.BoolVar(&flagAll, "a", false, "Traverse through all of the directories, including hidden ones")

	flag.Parse()

	if flagHelp {
		helpMessage()
		return
	}

	var filepath string
	if flag.NArg() == 0 {
		fmt.Println("Usage: libra <filepath>")
		return
	} else {
		filepath = flag.Arg(0)
	}

	// Check for online github repo
	githubRegexp := regexp.MustCompile("^.*github\\.com/([[:alnum:]]+)/([[:alnum:]]+)$")
	// If the filepath is a valid github repo
	if githubRegexp.MatchString(filepath) {
		match := githubRegexp.FindStringSubmatch(filepath)
		username := match[1]
		project := match[2]
		file, err := getGithubInfo(username, project)
		if err != nil {
			fmt.Println("Could not read from github.com/" + username + "/" + project)
			return
		}
		filepath = file
		defer deleteFile(filepath)
	}

	// Start file search
	wg.Add(1)
	go analizeDir(filepath)

	wg.Wait()

	analysis.DisplayReport()
}

// deleteFile deletes the tempoarary repository downloaded from github
func deleteFile(filepath string) {
	cmd := exec.Command("rm", "-rf", filepath)
	err := cmd.Run()
	if err != nil {
		fmt.Println("There was an error in removing temporary file. You may need to manually remove " + filepath)
	}
}

// getGithubInfo clones the specified github repo in a temporary repository
func getGithubInfo(username, project string) (string, error) {
	fileID := strconv.Itoa(rand.Intn(999999))
	cmd := exec.Command("git", "clone", "https://github.com/"+username+"/"+project+".git", "/tmp/"+"libra-"+fileID) // BUG: Different OS
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return "/tmp/" + "libra-" + fileID, nil
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
