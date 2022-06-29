package oci

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codetent/weasel/pkg/weasel/utils"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

func ExportRootFs(img v1.Image, ref name.Reference, path string) error {
	pullPath, err := ioutil.TempDir("", "weasel")
	if err != nil {
		return err
	}
	defer os.RemoveAll(pullPath)

	tarPath := filepath.Join(pullPath, "image.tar.gz")

	err = tarball.WriteToFile(tarPath, ref, img)
	if err != nil {
		return err
	}

	untaredPath := filepath.Join(pullPath, "content")
	err = utils.UntarPattern(tarPath, untaredPath)
	if err != nil {
		return err
	}

	archivePathCandidates, err := filepath.Glob(filepath.Join(untaredPath, "*.tar.gz"))
	if err != nil {
		return err
	} else if len(archivePathCandidates) == 0 {
		return fmt.Errorf("archive not found")
	}

	return utils.CopyFile(archivePathCandidates[0], path)
}
