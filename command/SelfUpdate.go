package command

import (
	"errors"
	"context"
	"runtime"
	"strings"
	"net/http"
	"github.com/google/go-github/github"
	"github.com/inconshreveable/go-update"
)

type SelfUpdate struct {
	CurrentVersion      string
	GithubOrganization  string
	GithubRepository    string
	GithubAssetTemplate string
	Force  bool   `long:"force"  description:"force update"`
}

var (
	selfUpdateOsTranslationMap = map[string]string{
		"darwin": "osx",
	}
	selfUpdateArchTranslationMap = map[string]string{
		"amd64": "x64",
		"386":   "x32",
	}
)

func (conf *SelfUpdate) Execute(args []string) error {
	Logger.Main("Starting self update")

	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), conf.GithubOrganization, conf.GithubRepository)

	if _, ok := err.(*github.RateLimitError); ok {
		Logger.Fatalln("GitHub rate limit, please try again later")
	}

	Logger.Step("latest version is %s", release.GetName())

	// check if latest version is current version
	if !conf.Force && release.GetName() == conf.CurrentVersion {
		Logger.Step("already using the latest version")
		return nil
	}

	// translate OS names
	os := runtime.GOOS
	if val, ok := selfUpdateOsTranslationMap[os]; ok {
		os = val
	}

	// translate arch names
	arch := runtime.GOARCH
	if val, ok := selfUpdateArchTranslationMap[arch]; ok {
		arch = val
	}

	// build asset name
	assetName := conf.GithubAssetTemplate
	assetName = strings.Replace(assetName, "%OS%", os, -1)
	assetName = strings.Replace(assetName, "%ARCH%", arch, -1)

	// search assets in release for the desired filename
	Logger.Step("searching for asset \"%s\"", assetName)
	for _, asset := range release.Assets {
		if asset.GetName() == assetName {
			downloadUrl := asset.GetBrowserDownloadURL()
			Logger.Step("found new update url \"%s\"", downloadUrl)
			conf.runUpdate(downloadUrl)
			Logger.Step("finished update to version %s", release.GetName())
			return nil
		}
	}

	return errors.New("unable to find asset, please contact maintainer")
}

func (conf *SelfUpdate) runUpdate(url string) error {
	Logger.Step("downloading update")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	Logger.Step("applying update")
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		// error handling
		Logger.Step("updating application failed: %s", err)
	}
	return err
}
