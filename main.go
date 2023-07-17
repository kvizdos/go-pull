package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
	"time"
)

const versionFilePath = ".pipeline-version"

var runningVersion *os.Process

type Asset struct {
	Url string `json:"url"`
}

type Release struct {
	Tag       string  `json:"tag_name"`
	Assets    []Asset `json:"assets"`
	CreatedAt string  `json:"published_at"`
}

func main() {
	err := godotenv.Load(".pipeline.env")
	if err != nil {
		panic(fmt.Sprintf("Could not find env file: %s", err))
	}

	if _, ok := os.LookupEnv("PIPELINE_GITHUB_TOKEN"); !ok {
		panic("Requires PIPELINE_GITHUB_TOKEN to be set.")
	}
	if _, ok := os.LookupEnv("PIPELINE_RELEASES_API"); !ok {
		panic("Requires PIPELINE_RELEASES_API to be set.")
	}
	if _, ok := os.LookupEnv("PIPELINE_BUILD_OUT"); !ok {
		panic("Requires PIPELINE_BUILD_OUT to be set.")
	}

	for {
		release, err := GetLatestReleases(os.Getenv("PIPELINE_RELEASES_API"))

		if err != nil {
			if err != nil && err.Error() == "no new version" && runningVersion != nil {
				fmt.Println("no new version found")
				time.Sleep(1 * time.Hour)
				continue
			} else if err.Error() == "no new version" && runningVersion == nil {
				fmt.Println("First run with pre-existing executable, no running executable.")
			} else {
				fmt.Println(err)
				time.Sleep(1 * time.Hour)
				continue
			}
		}

		if err == nil || (err != nil && err.Error() != "no new version") {
			token := os.Getenv("PIPELINE_GITHUB_TOKEN")

			latestAsset := release.Assets[0]

			DownloadRelease(latestAsset.Url, release.Tag, token, os.Getenv("PIPELINE_BUILD_OUT"))

			fmt.Println("File downloaded successfully.")
		}
		if runningVersion != nil {
			fmt.Println("Stopping previous process..")
			StopPreviousProcess()
		}

		fmt.Println("Starting new process..")
		StartNewProcess(os.Getenv("PIPELINE_BUILD_OUT"))

		time.Sleep(1 * time.Hour)
	}
}

func getCachedVersion() (string, error) {
	if _, err := os.Stat(versionFilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("version file does not exist")
	}

	content, err := ioutil.ReadFile(versionFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read version file: %v", err)
	}

	return string(content), nil
}

func setCachedVersion(version string) error {
	err := ioutil.WriteFile(versionFilePath, []byte(version), 0644)
	if err != nil {
		return fmt.Errorf("failed to write version file: %v", err)
	}

	return nil
}
