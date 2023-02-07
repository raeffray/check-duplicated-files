package writer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type CheckResult struct {
	FirstPath  string `json:"firstPath"`
	SecondPath string `json:"secondPath"`
	Hash       string `json:"hash"`
	Size       int64  `json:"size"`
}

func AddToJSONFile(check CheckResult, fileName string) error {
	var checks []CheckResult
	// read the file to get the current array of checks
	fileData, err := ioutil.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			// if the file does not exist, create an empty array of checks
			checks = []CheckResult{}
			err := os.WriteFile(fileName, []byte("[]"), 0644)

			if err != nil {
				fmt.Println(err)
			}

		} else {
			return err
		}
	} else {
		// if the file exists, unmarshal the JSON data into the checks array
		err = json.Unmarshal(fileData, &checks)
		if err != nil {
			return err
		}
	}
	// append the new check to the array of checks
	checks = append(checks, check)
	// write the updated array of checks to the file
	resultJson, err := json.Marshal(checks)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, resultJson, 0644)
}
