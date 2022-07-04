package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/codetent/weasel/pkg/weasel/utils"
	"github.com/codetent/weasel/pkg/weasel/wsl"
	log "github.com/sirupsen/logrus"
)

type GitPlugin struct {
}

func getCredManagerPath() (string, error) {
	gitPath, err := utils.Which("git")
	if err != nil {
		return "", err
	}

	gitRoot := filepath.Dir(filepath.Dir(gitPath))
	credManagerPaths := []string{
		gitRoot + "\\mingw64\\bin\\git-credential-manager-core.exe",
		gitRoot + "\\mingw64\\libexec\\git-core\\git-credential-manager-core.exe",
	}

	for _, path := range credManagerPaths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return path, nil
		}
	}

	return "", fmt.Errorf("git credential manager not found")
}

func setupCredManager(name string) error {
	credManagerPath, err := getCredManagerPath()
	if err != nil {
		log.Debug("No git credential manager found. Skipping")
		return nil
	}

	transPath, err := wsl.ExecuteSilently(name, "wslpath -u \""+credManagerPath+"\"")
	if err != nil {
		return err
	}

	log.Info("Configuring git credential manager")

	_, err = wsl.ExecuteSilently(name, "git config --global credential.helper \""+transPath+"\"")
	return err
}

func setupUserName(name string) error {
	cmd := exec.Command("git", "config", "--global", "--get", "user.name")
	cmdBytes, err := cmd.Output()
	if err != nil {
		return err
	}

	value := strings.TrimSpace(string(cmdBytes))
	_, err = wsl.ExecuteSilently(name, "git config --global user.name "+value)
	return err
}

func setupUserEmail(name string) error {
	cmd := exec.Command("git", "config", "--global", "--get", "user.email")
	cmdBytes, err := cmd.Output()
	if err != nil {
		return err
	}

	value := strings.TrimSpace(string(cmdBytes))
	_, err = wsl.ExecuteSilently(name, "git config --global user.email "+value)
	return err
}

func (plugin *GitPlugin) Enter(name string) error {
	_, err := wsl.ExecuteSilently(name, "which git")
	if err != nil {
		log.Debug("No git in distribution found. Skipping")
		return nil
	}

	err = setupCredManager(name)
	if err != nil {
		return fmt.Errorf("git credential manager could not be enabled: %v", err)
	}

	err = setupUserName(name)
	if err != nil {
		return fmt.Errorf("git user name could not be set: %v", err)
	}

	err = setupUserEmail(name)
	if err != nil {
		return fmt.Errorf("git user email could not be set: %v", err)
	}

	return err
}
