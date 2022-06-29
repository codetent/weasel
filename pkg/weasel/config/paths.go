package config

import (
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func GetConfigRoot(file *ConfigFile) string {
	return filepath.Dir(file.Path)
}

func GetWeaselDir(file *ConfigFile) string {
	root := GetConfigRoot(file)
	return filepath.Join(root, ".weasel")
}

func GetArchiveCacheRoot(file *ConfigFile) string {
	weasel := GetWeaselDir(file)
	return filepath.Join(weasel, "cache")
}

func GetArchiveCachePath(file *ConfigFile, ref name.Reference, digest v1.Hash) string {
	root := GetArchiveCacheRoot(file)
	return filepath.Join(root, ref.Context().Name(), digest.Hex[:12]+".tar.gz")
}

func GetWorkspaceCacheRoot(file *ConfigFile) string {
	weasel := GetWeaselDir(file)
	return filepath.Join(weasel, "workspaces")
}

func GetWorkspaceCachePath(file *ConfigFile, ref name.Reference, digest v1.Hash) string {
	root := GetWorkspaceCacheRoot(file)
	return filepath.Join(root, ref.Context().Name(), digest.Hex[:12])
}
