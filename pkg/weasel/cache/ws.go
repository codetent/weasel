package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/yuk7/wsllib-go"
)

func GetWorkspacePath(dist string) (string, error) {
	root, err := GetWorkspaceRoot()
	if err != nil {
		return "", fmt.Errorf("GetWorkspacePath: GetWorkspaceRoot(): %v", err)
	}

	dir := filepath.Join(root, dist)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("GetWorkspacePath: MkdirAll(): %v", err)
	}

	return dir, nil
}

func CleanWorkspaceCache() error {
	lock, err := GetLock()
	if err != nil {
		return err
	}
	defer lock.Unlock()

	root, err := GetWorkspaceRoot()
	if err != nil {
		return fmt.Errorf("CleanWorkspaceCache: GetWorkspaceRoot(): %v", err)
	}

	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return fmt.Errorf("CleanWorkspaceCache: ReadDir(): %v", err)
	}

	for _, dir := range fileInfo {
		if dir.IsDir() && !wsllib.WslIsDistributionRegistered(dir.Name()) {
			path := filepath.Join(root, dir.Name())
			err = os.RemoveAll(path)
			if err != nil {
				return fmt.Errorf("CleanWorkspaceCache: RemoveAll(): %v", err)
			}
		}
	}

	return nil
}
