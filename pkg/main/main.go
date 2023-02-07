package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	calculator "check-files/pkg/calculator"
	"check-files/pkg/config"
	builder "check-files/pkg/repository"
	writer "check-files/pkg/writer"
)

type CheckResult struct {
	FirstPath  string
	SecondPath string
	Hash       string
	Size       int64
}

var repo = builder.NewRedisRepository("fileCheck")

var cfg = config.NewConfig()

func getSupportedFiles() []string {
	return cfg.GetListWithDefault("app.supportedfiles", "jpg,jpeg,png,gif,mp4,mpeg")
}

func isSupportedFile(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	for _, supportedExt := range getSupportedFiles() {
		if ext == "."+supportedExt {
			return true
		}
	}
	return false
}

// using a stack to traverse the directory tree
// making it tail recursive, so that we mitigate the risk of stack overflow
func readDir(path string, stack []string) {
	// Check if the stack is empty
	if len(stack) == 0 {
		return
	}

	// Pop the last item from the stack
	currentPath := stack[len(stack)-1]
	stack = stack[:len(stack)-1]

	files, err := ioutil.ReadDir(currentPath)

	if err != nil {
		fmt.Println("Error reading directory:", err)
	} else {
		for _, file := range files {

			fullPath := currentPath + "/" + file.Name()

			if file.IsDir() {
				// If the file is a directory, add it to the stack
				stack = append(stack, fullPath)
			} else {
				// check whether the file is supported
				checkFile(file, fullPath)
			}
		}
	}

	// Recursively call readDir with the updated stack
	readDir(path, stack)

}

func checkFile(file os.FileInfo, fullPath string) {
	fileExtension := strings.ToLower(filepath.Ext(fullPath))
	if isSupportedFile(fileExtension) {
		// read the file
		content, err := ioutil.ReadFile(fullPath)
		// if there is an error reading the file, skip it
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		// calculate the hash of the file
		hash := calculator.CalculatetHash(content)

		// check if the hash is already in the repository
		storedPath, err := repo.GetValue(hash)

		// if the hash is not in the repository, save it and return
		if err != nil {
			repo.SaveValue(hash, fullPath)
			return
		}

		// if the hash is in the repository, check if the file path is the same, if so, return
		if storedPath == fullPath {
			return
		}

		// if the hash is in the repository, but the file path is different, save the check result to the output file
		// and that means this is a duplicate file
		original := strings.Replace(storedPath, " ", "\\ ", -1)
		duplicated := strings.Replace(fullPath, " ", "\\ ", -1)

		// create the check result
		check := writer.CheckResult{
			FirstPath:  original,
			SecondPath: duplicated,
			Hash:       hash,
			Size:       file.Size(),
		}
		// write the check result to the output file
		errorWriting := writer.AddToJSONFile(check, cfg.GetWithDefault("app.outputfile", "output.json"))
		// if there is an error writing to the file, print it
		if errorWriting != nil {
			fmt.Println(err)
		}

	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a path")
		return
	}

	path := os.Args[1]
	readDir(path, []string{path})
}
