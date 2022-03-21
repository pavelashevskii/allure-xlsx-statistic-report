package pkg

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ParseFiles(path string) []TestObject {
	var testObject []TestObject
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("Reading %s is failed: %s", path, err)
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Printf("%s is not file\n", file.Name())
		} else {
			testObject = append(testObject, parseFile(path, file))
		}
	}
	return testObject
}

func parseFile(path string, file fs.FileInfo) TestObject {
	fileData, err := os.ReadFile(filepath.Join(path, file.Name()))
	if err != nil {
		log.Fatal("Reading file is failed", file.Name(), err)
	}

	var testObject TestObject
	err = json.Unmarshal(fileData, &testObject)
	if err != nil {
		log.Fatal("Parsing file is failed", file.Name(), err)
	}

	return testObject
}
