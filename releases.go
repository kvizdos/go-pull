package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func DownloadRelease(githubAssetURL string, version string, token string, outBinary string) {
	fmt.Printf("Downloading new release %s...\n", version)
	// Create the output file
	outputFile, err := os.Create(outBinary)
	outputFile.Chmod(0700)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Send an HTTP GET request to the file URL

	dlReq, err := http.NewRequest("GET", githubAssetURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	dlReq.Header.Set("Accept", "application/octet-stream")
	dlReq.Header.Set("Authorization", "Bearer "+token)
	dlReq.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(dlReq)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// Check if the response status code is successful (200)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error response:", resp.Status)
		return
	}

	// Copy the response body to the output file
	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		fmt.Println("Error copying response body to file:", err)
		return
	}

	err = setCachedVersion(version)
	if err != nil {
		fmt.Println("Error adding metadata to file:", err)
		return
	}
}

type Releases []Release

func (rs Releases) GetLatestRelease() Release {
	var latest int64
	latest = 0
	var lr Release

	for _, release := range rs {
		layout := "2006-01-02T15:04:05Z"

		t, err := time.Parse(layout, release.CreatedAt)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			continue
		}

		if t.Unix() > latest {
			latest = t.Unix()
			lr = release
		}
	}

	return lr
}

func GetLatestReleases(url string) (Release, error) {
	cachedVersion, cachedErr := getCachedVersion()

	if cachedErr != nil {
		if cachedErr.Error() != "version file does not exist" {
			fmt.Println(cachedErr)
			return Release{}, cachedErr
		}
	}

	token := os.Getenv("PIPELINE_GITHUB_TOKEN")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return Release{}, cachedErr
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return Release{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return Release{}, err
	}

	var releases Releases
	err = json.Unmarshal(body, &releases)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return Release{}, err
	}

	latestRelease := releases.GetLatestRelease()

	if latestRelease.Tag == cachedVersion {
		return Release{}, fmt.Errorf("no new version")
	}

	return latestRelease, nil
}
