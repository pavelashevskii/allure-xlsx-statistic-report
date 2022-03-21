package pkg

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const layoutTime string = "2006-01-02 15:04:05"
const layoutDate string = "2006-01-02"

func GetServiceAndClass(fullName string) (string, string) {
	r := regexp.MustCompile(`com\.gett\.automation\.tests\.services\.(?P<service>[A-Za-z0-9]+)(\.[A-Za-z0-9]+)*\.(?P<testname>[A-Za-z0-9]+)\.([A-Za-z0-9]+)`)
	matches := r.FindStringSubmatch(fullName)
	serviceIndex := r.SubexpIndex("service")
	testnameIndex := r.SubexpIndex("testname")

	return matches[serviceIndex], matches[testnameIndex]
}

func GetUrlToGithub(fullName string) string {
	subUrl := strings.ReplaceAll(fullName[0:strings.LastIndex(fullName, ".")], ".", "/")
	return fmt.Sprintf("https://github.com/gtforge/automation_tests/blob/master/gett-tests/src/test/java/%s.java", subUrl)
}

func ParseStringAsTimestamp(timeStr string) uint64 {
	if timeStr == "" {
		return 0
	}
	result, err := time.Parse(layoutTime, timeStr)
	if err != nil {
		result, err = time.Parse(layoutDate, timeStr)
		if err != nil {
			log.Fatal(err)
		}
	}

	return uint64(result.UnixMilli())
}

func Unzip(source, destination string) error {
	archive, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer archive.Close()
	for _, file := range archive.Reader.File {
		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()
		path := filepath.Join(destination, file.Name)
		// Remove file if it already exists; no problem if it doesn't; other cases can error out below
		_ = os.Remove(path)
		// Create a directory at path, including parents
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		// If file is _supposed_ to be a directory, we're done
		if file.FileInfo().IsDir() {
			continue
		}
		// otherwise, remove that directory (_not_ including parents)
		err = os.Remove(path)
		if err != nil {
			return err
		}
		// and create the actual file.  This ensures that the parent directories exist!
		// An archive may have a single file with a nested path, rather than a file for each parent dir
		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer writer.Close()
		_, err = io.Copy(writer, reader)
		if err != nil {
			return err
		}
	}
	return nil
}
