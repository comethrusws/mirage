package updater

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"mirage/internal/logger"
)

const (
	githubAPI = "https://api.github.com/repos/comethrusws/mirage/releases/latest"
	repoOwner = "comethrusws"
	repoName  = "mirage"
)

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func Update(currentVersion string) error {
	logger.LogInfo("Checking for updates...")
	
	release, err := fetchLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to fetch latest release: %w", err)
	}
	
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion = strings.TrimPrefix(currentVersion, "v")
	
	if latestVersion == currentVersion {
		logger.LogSuccess("Already on latest version: " + currentVersion)
		return nil
	}
	
	logger.LogInfo(fmt.Sprintf("New version available: %s (current: %s)", latestVersion, currentVersion))
	
	assetName := getAssetName()
	var downloadURL string
	
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}
	
	if downloadURL == "" {
		return fmt.Errorf("no release found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	
	logger.LogInfo("Downloading update...")
	
	tmpFile, err := downloadFile(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer os.Remove(tmpFile)
	
	logger.LogInfo("Extracting binary...")
	
	binaryPath, err := extractBinary(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to extract: %w", err)
	}
	defer os.Remove(binaryPath)
	
	logger.LogInfo("Replacing current binary...")
	
	if err := replaceBinary(binaryPath); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}
	
	logger.LogSuccess(fmt.Sprintf("Successfully updated to version %s", latestVersion))
	logger.LogInfo("Please restart mirage to use the new version")
	
	return nil
}

func fetchLatestRelease() (*Release, error) {
	resp, err := http.Get(githubAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}
	
	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	
	return &release, nil
}

func getAssetName() string {
	osName := runtime.GOOS
	archName := runtime.GOARCH
	
	if osName == "darwin" {
		osName = "Darwin"
	} else if osName == "linux" {
		osName = "Linux"
	}
	
	if archName == "amd64" {
		archName = "x86_64"
	} else if archName == "arm64" {
		archName = "arm64"
	}
	
	return fmt.Sprintf("mirage_%s_%s.tar.gz", osName, archName)
}

func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	tmpFile, err := os.CreateTemp("", "mirage-update-*.tar.gz")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()
	
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}
	
	return tmpFile.Name(), nil
}

func extractBinary(tarPath string) (string, error) {
	file, err := os.Open(tarPath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzr.Close()
	
	tr := tar.NewReader(gzr)
	
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		
		if header.Name == "mirage" {
			tmpFile, err := os.CreateTemp("", "mirage-binary-*")
			if err != nil {
				return "", err
			}
			defer tmpFile.Close()
			
			if _, err := io.Copy(tmpFile, tr); err != nil {
				os.Remove(tmpFile.Name())
				return "", err
			}
			
			if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
				os.Remove(tmpFile.Name())
				return "", err
			}
			
			return tmpFile.Name(), nil
		}
	}
	
	return "", fmt.Errorf("mirage binary not found in archive")
}

func replaceBinary(newBinaryPath string) error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return err
	}
	
	backupPath := execPath + ".backup"
	if err := os.Rename(execPath, backupPath); err != nil {
		return err
	}
	
	if err := copyFile(newBinaryPath, execPath); err != nil {
		os.Rename(backupPath, execPath)
		return err
	}
	
	os.Remove(backupPath)
	
	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()
	
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}
	
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	
	return os.Chmod(dst, sourceInfo.Mode())
}
