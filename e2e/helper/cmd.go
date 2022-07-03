package helper

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/onsi/gomega/gexec"
)

func NewWeaselCommand(args ...string) (*exec.Cmd, error) {
	execPath, err := gexec.Build("github.com/codetent/weasel")
	if err != nil {
		return nil, err
	}

	homePath, err := ioutil.TempDir("", "weasel")
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(execPath, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "WEASEL_HOME="+homePath)

	return cmd, nil
}
