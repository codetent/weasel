package helper

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codetent/weasel/pkg/weasel/utils"
)

func CreateEmptyWorkspace() (string, error) {
	return ioutil.TempDir("", "weasel-test")
}

func CreateConfigWorkspace(config string) (string, error) {
	dir, err := CreateEmptyWorkspace()
	if err != nil {
		return "", err
	}

	cfgFilePath := filepath.Join(dir, "weasel.yml")
	err = utils.CopyFile(config, cfgFilePath)
	if err != nil {
		os.RemoveAll(dir)
		return "", err
	}

	return dir, err
}
