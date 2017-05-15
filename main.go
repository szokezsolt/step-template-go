package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"io/ioutil"

	"github.com/bitrise-io/go-utils/log"
)

// ConfigsModel ...
type ConfigsModel struct {
	ExampleInput string
	DownloadURL  string
	DownloadPth  string
}

// JSONResultModel ...
type JSONResultModel struct {
	android Android
	ios     IOS
}

// Android ...
type Android struct {
	release Release
}

// IOS ...
type IOS struct {
	debug Debug
}

// Release ...
type Release struct {
	keystore      string
	storePassword string
	alias         string
	password      string
	isPresent     bool
}

// Debug ...
type Debug struct {
	UID              int
	codeSignIdentity string
	developmentTeam  string
	packageType      string
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		ExampleInput: os.Getenv("example_step_input"),
		DownloadURL:  os.Getenv("download_url"),
		DownloadPth:  os.Getenv("download_path"),
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

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
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

	// STEP 3
	if err := downloadFile(configs.DownloadURL, configs.DownloadPth); err != nil {
		failf("Failed to download json file, error: %s", err)
	}

	// STEP 4
	jsonLine, err := ioutil.ReadFile(configs.DownloadPth)
	if err != nil {
		log.Errorf("File read error")
	}
	var result JSONResultModel

	if err := json.Unmarshal([]byte(jsonLine), &result); err != nil {
		log.Errorf("Failed to unmarshal result, error: %s", err)
	}

	log.Infof("%#v", result)
}
