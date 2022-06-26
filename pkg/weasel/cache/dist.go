package cache

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetDistPath(dist string) (string, error) {
	root, err := GetDistRoot()
	if err != nil {
		return "", fmt.Errorf("GetDistPath: GetDistRoot(): %v", err)
	}

	path := filepath.Join(root, dist+".tgz")
	return path, nil
}

func CleanDistCache() error {
	lock, err := GetLock()
	if err != nil {
		return err
	}
	defer lock.Unlock()

	root, err := GetDistRoot()
	if err != nil {
		return fmt.Errorf("CleanDistCache: GetDistCache(): %v", err)
	}

	err = os.RemoveAll(root)
	if err != nil {
		return fmt.Errorf("CleanDistCache: RemoveAll(): %v", err)
	}

	return nil
}
