package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"stats/pkg"

	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Path     string `cli:"p,path" usage:"path to allure root directory or allure.zip"`
	From     string `cli:"from" usage:"time filter (2006-01-02 15:04:05 or 2006-01-02 format)"`
	Till     string `cli:"till" usage:"time filter (2006-01-02 15:04:05 or 2006-01-02 format)"`
	Pattern  string `cli:"pattern" usage:"regex for tests names"`
	Services string `cli:"s,services" usage:"list of servicess (, - delimeter)"`
	Output   string `cli:"o,output" usage:"save report by special path" dft:"Report.xlsx"`
	Envs     string `cli:"env" usage:"list of environmnents" dft:"il,ru"`
}

func main() {
	os.Exit(cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		var pathToMeta string

		path, err := os.Open(argv.Path)
		if err != nil {
			log.Fatalf("Reading %s is failed: %s", argv.Path, err)
		}

		stat, err := path.Stat()
		if err != nil {
			log.Fatalf("Getting stat %s is failed: %s", argv.Path, err)
		}

		if stat.IsDir() {
			pathToMeta = filepath.Join(argv.Path, "data", "test-cases")
		} else {
			tmpDir, err := os.MkdirTemp("", "allure-stats")
			if err != nil {
				log.Fatalf("Creating tmp dir %s is failed: %s", tmpDir, err)
			}
			defer os.RemoveAll(tmpDir)
			reader, err := zip.OpenReader(argv.Path)
			if err != nil {
				return err
			}
			defer reader.Close()

			pkg.Unzip(argv.Path, tmpDir)
			if err != nil {
				fmt.Println("Unzip", err)
			}
			pathToMeta = filepath.Join(tmpDir, "allure-report", "data", "test-cases")
		}

		testObjects := pkg.ParseFiles(pathToMeta)
		testObjects = pkg.FilterTestObjects(testObjects, pkg.ParseStringAsTimestamp(argv.From),
			pkg.ParseStringAsTimestamp(argv.Till), argv.Services, argv.Services, argv.Envs)
		pkg.PrepareReport(testObjects, argv.Output)

		fmt.Println("Report is done!")
		return nil
	}))
}
