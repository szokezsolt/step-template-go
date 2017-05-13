package main

import (
	"errors"
	"fmt"
	"os"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

// ConfigsModel ...
type ConfigsModel struct {
	ExampleInput string
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		ExampleInput: os.Getenv("example_step_input"),
	}
}

func downloadFile(downloadURL, targetPath string) error {
	outFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create (%s), error: %s", targetPath, err)
	}
	defer func() {
		if err := outFile.Close(); err != nil {
			log.Warnf("Failed to close (%s)", targetPath)
		}
	}()

	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download from (%s), error: %s", downloadURL, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Warnf("failed to close (%s) body", downloadURL)
		}
	}()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to download from (%s), error: %s", downloadURL, err)
	}

	return nil
}

func (configs ConfigsModel) print() {
	log.Infof("Configs:")
	log.Printf("- ExampleInput: %s", configs.ExampleInput)
}

func (configs ConfigsModel) validate() error {
	if configs.ExampleInput == "" {
		return errors.New("no ExampleInput parameter specified")
	}

	return nil
}

func main() {
	// Input validation
	configs := createConfigsModelFromEnvs()

	fmt.Println()
	configs.print()

	if err := configs.validate(); err != nil {
		log.Errorf("Issue with input: %s", err)
		os.Exit(1)
	}

	fmt.Println()

	// Main
	
	// downloadUrl := DOWNLOAD URL TO BE IMPORTED SOMEHOW FROM STEP.YML
	// downloadPth := TEMPORARY DIRECTORY PATH TO BE DOWNLOADED TO
	if err := downloadFile(downloadUrl, downloadPth); err != nil {
		failf("Failed to download json file, error: %s", err)
	}
}
