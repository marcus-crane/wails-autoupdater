package autoupdater

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

const (
	UpdateSource   = "marcus-crane/wails-autoupdater"
	ExecutableName = "wails-autoupdater"
)

func downloadLatestVersion() (string, error) {
	latest, found, err := selfupdate.DetectLatest(UpdateSource)
	if err != nil {
		return "", err
	}
	if !found {
		return "", err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Print(err)
		return "", err
	}

	downloadPath := filepath.Join(homeDir, "Downloads", fmt.Sprintf("%s.zip", ExecutableName))

	out, err := os.Create(downloadPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	resp, err := http.Get(latest.AssetURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return downloadPath, err
	}

	log.Printf("Downloaded %s", downloadPath)

	return downloadPath, nil
}

func unzipLatestVersion(downloadPath string) error {
	destination := filepath.Dir(downloadPath)
	archive, err := zip.OpenReader(downloadPath)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(destination, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			continue
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return nil
}

func PerformUpdateWindows() (bool, error) {
	downloadPath, err := downloadLatestVersion()
	if err != nil {
		return false, err
	}
	if err = unzipLatestVersion(downloadPath); err != nil {
		return false, err
	}
	return true, nil
}

func CleanupOldDarwinApp() error {
	oldAppLocation := filepath.Join("/tmp", fmt.Sprintf("%s.app", ExecutableName))
	if _, err := os.Stat(oldAppLocation); !os.IsNotExist(err) {
		if err = exec.Command("rm", "-rf", oldAppLocation).Run(); err != nil {
			return err
		}
	}
	return nil
}

func RestartDarwinApp() error {
	pid := strconv.Itoa(os.Getpid())
	installPath := filepath.Join("/Applications/", fmt.Sprintf("%s.app", ExecutableName))

	if err := exec.Command("open", installPath).Run(); err != nil {
		return err
	}

	if err := exec.Command("kill", "-3", pid).Run(); err != nil {
		return err
	}

	return nil
}

func PerformUpdateDarwin() (bool, error) {
	f, err := os.OpenFile("/Users/marcus/Desktop/test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("failed to create file")
	}
	defer f.Close()
	log.SetOutput(f)
	downloadPath, err := downloadLatestVersion()
	if err != nil {
		return false, err
	}
	if err = unzipLatestVersion(downloadPath); err != nil {
		return false, err
	}
	var installPath string
	cmdPath, err := os.Executable()
	installPath = strings.TrimSuffix(cmdPath, fmt.Sprintf("%s.app/Contents/MacOS/%s", ExecutableName, ExecutableName))
	if err != nil {
		log.Print(err)
		installPath = "/Applications/"
	}
	log.Printf("Going to install to %s", installPath)

	appLocation := filepath.Join(installPath, fmt.Sprintf("%s.app", ExecutableName))
	dlAppLocation := strings.Replace(downloadPath, ".zip", ".app", -1)

	if err := CleanupOldDarwinApp(); err != nil {
		log.Printf("Failed to clean up old darwin app")
		return false, err
	}

	if err := exec.Command("mv", appLocation, fmt.Sprintf("/tmp/%s.app", ExecutableName)).Run(); err != nil {
		log.Printf("Failed to mv %s to /tmp: %+v", appLocation, err)
		return false, err
	}

	if err := exec.Command("mv", dlAppLocation, appLocation).Run(); err != nil {
		log.Printf("Failed to run mv %s %s: %+v", dlAppLocation, appLocation, err)
		return false, err
	}

	if err := exec.Command("rm", downloadPath).Run(); err != nil {
		// LOG: Failed to cleanup tmp folder
		log.Printf("Failed to rm %s: %+v", downloadPath, err)
		return false, err
	}

	log.Println("Successfully updated")

	return true, nil
}

func CheckForNewerVersion(currentVersion string) (bool, string) {
	latest, found, err := selfupdate.DetectLatest(UpdateSource)
	if err != nil {
		return false, ""
	}

	if !found {
		// LOG: Update manifest not found
		return false, ""
	}

	v, err := semver.Parse(currentVersion)
	if err != nil {
		return false, ""
	}

	return compareVersions(v, latest.Version)
}

func compareVersions(currentVersion semver.Version, latest semver.Version) (bool, string) {
	if latest.LTE(currentVersion) {
		return false, ""
	}
	return true, latest.String()
}
