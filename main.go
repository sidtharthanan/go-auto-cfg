package main

import (
	"fmt"
	"io/ioutil"
	"os"
	fp "path/filepath"

	"github.com/sidtharthanan/go-auto-cfg/transpiler"
	_ "github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	inputSchemaFilepath := kingpin.Arg("schema-file", "Definitions of the configs read from here.").Required().ExistingFile()
	outputFilePath := kingpin.Arg("output-file", "Config loader go file is created here.").Required().String()
	kingpin.Parse()
	kingpin.FatalIfError(mainAction(*inputSchemaFilepath, *outputFilePath), "AutoCfg generation failed")
}

func mainAction(filename string, outputFilePath string) error {
	packageName, err := ensureParentDir(outputFilePath)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	outputWriter, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	return transpiler.Transpile(filename, content, packageName, outputWriter)
}

func ensureParentDir(filePath string) (string, error) {
	parentDirPath, filename := fp.Split(filePath)
	parentDirName := fp.Base(parentDirPath)
	if filename == "" {
		return "", fmt.Errorf("filename missing in filpath '%v'", filePath)
	}
	if parentDirName == "." {
		parentDirPath, _ = fp.Abs(parentDirPath)
		parentDirName = fp.Base(parentDirPath)
	}

	if _, err := os.Stat(parentDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(parentDirPath, 0755); err != nil {
			return "", fmt.Errorf("unable to create directory '%v'\n%v", parentDirPath, err)
		}
	} else if err != nil {
		return "", fmt.Errorf("unable to create directory '%v'\n%v", parentDirPath, err)
	}
	return parentDirName, nil
}
