package beholder

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skratchdot/open-golang/open"
)

const updatesURL = "https://api.github.com/repos/dhleong/beholder/releases/latest"
const downloadURL = "https://github.com/dhleong/beholder/releases/latest"

type release struct {
	Name string `json:"name"`
}

// CheckForUpdates returns the version string of the latest update,
// or an empty string if there is no update
func CheckForUpdates() string {

	// request the remote file
	client := &http.Client{}

	req, err := http.NewRequest("GET", updatesURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Add("User-Agent", fmt.Sprintf("beholder %s", Version))

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	result := release{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return ""
	}

	if result.Name == Version {
		return ""
	}
	return result.Name
}

// LaunchUpdateDownload .
func LaunchUpdateDownload() {
	if err := open.Run(downloadURL); err != nil {
		panic(err)
	}
}
