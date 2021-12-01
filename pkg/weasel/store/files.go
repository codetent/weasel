package store

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetCacheRoot() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("GetCacheRoot: UserHomeDir(): %v", err)
	}

	root := filepath.Join(userHome, ".weasel")
	err = os.MkdirAll(root, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("GetCacheRoot: MkdirAll(): %v", err)
	}

	return root, nil
}

func GetDistCache() (string, error) {
	root, err := GetCacheRoot()
	if err != nil {
		return "", fmt.Errorf("GetDistCache: GetCacheRoot(): %v", err)
	}

	path := filepath.Join(root, "dists")
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("GetDistCache: MkdirAll(): %v", err)
	}

	return path, nil
}

func GetWorkspaceRoot() (string, error) {
	root, err := GetCacheRoot()
	if err != nil {
		return "", fmt.Errorf("GetWorkspaceRoot: GetCacheRoot(): %v", err)
	}

	path := filepath.Join(root, "workspaces")
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("GetWorkspaceRoot: MkdirAll(): %v", err)
	}

	return path, nil
}

func GetBuilderRoot() (string, error) {
	root, err := GetCacheRoot()
	if err != nil {
		return "", fmt.Errorf("GetBuilderRoot: GetCacheRoot(): %v", err)
	}

	path := filepath.Join(root, "builder")
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("GetBuilderRoot: MkdirAll(): %v", err)
	}

	return path, nil
}
