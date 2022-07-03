package config

import (
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func GetWeaselDir() (string, error) {
	root := os.Getenv("WEASEL_HOME")
	if root == "" {
		var err error
		root, err = os.UserCacheDir()
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(root, "weasel"), nil
}

func GetImageCacheRoot() (string, error) {
	weasel, err := GetWeaselDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(weasel, "images"), nil
}

func GetImageCachePath(ref name.Reference, digest v1.Hash) (string, error) {
	root, err := GetImageCacheRoot()
	if err != nil {
		return "", err
	}

	return filepath.Join(root, ref.Context().Name(), digest.Hex[:12]+".tar.gz"), nil
}

func GetWorkspaceCacheRoot() (string, error) {
	weasel, err := GetWeaselDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(weasel, "workspaces"), nil
}

func GetWorkspaceCachePath(ref name.Reference, digest v1.Hash) (string, error) {
	root, err := GetWorkspaceCacheRoot()
	if err != nil {
		return "", err
	}

	return filepath.Join(root, ref.Context().Name(), digest.Hex[:12]), nil
}
