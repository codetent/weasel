/*
Copyright © 2021 Christoph Swoboda

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/codetent/weasel/pkg/weasel"
	"github.com/codetent/weasel/pkg/weasel/docker"
	"github.com/codetent/weasel/pkg/weasel/store"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var (
	PullParent   bool
	DisableCache bool
	Tags         []string
	DockerFile   string
	BuildArgs    map[string]string
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build distribution",
	Long: `Build distribution archive that can later be loaded using WSL.
	
This command supports multiple input sources.
They are differented by specifying a prefix before the argument.
- "context:<value>" - Docker image using Dockerfile. The value is a path to an existing directory (context).
- "hub:<value>" - Docker image from docker hub. The value is the tag of a public available image.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		code := func() int {
			var ref string
			var stream io.ReadCloser
			var err error

			if strings.HasPrefix(args[0], "context:") {
				contextPath := strings.TrimPrefix(args[0], "context:")

				ref = weasel.GenerateId()
				stream, err = buildDockerImage(contextPath, ref)
				if err != nil {
					log.Errorf("Error building image: %v", err)
					return 1
				}
			} else if strings.HasPrefix(args[0], "hub:") {
				ref = strings.TrimPrefix(args[0], "hub:")
				stream, err = docker.ImagePull(ref)
				if err != nil {
					log.Errorf("Error pulling image: %v", err)
					return 1
				}
			} else {
				log.Errorf("Unknown specifier %s", args[0])
				return 1
			}
			defer stream.Close()

			// Show live build output
			outStream := os.Stdout
			if ErrorsOnly {
				outStream, _ = os.Open(os.DevNull)
				defer outStream.Close()
			}
			termFd, isTerm := term.GetFdInfo(outStream)
			jsonmessage.DisplayJSONMessagesStream(stream, outStream, termFd, isTerm, nil)

			// Get image id
			imageId, err := docker.ImageIdByTag(ref)
			if err != nil {
				log.Errorf("Error getting image: %v", err)
				return 1
			}

			// Export image to archive
			cacheDir, err := store.GetDistCache()
			if err != nil {
				log.Errorf("Error getting distribution cache: %v", err)
				return 1
			}

			archivePath := filepath.Join(cacheDir, imageId+".tgz")
			_, err = os.Stat(archivePath)
			if errors.Is(err, os.ErrNotExist) {
				err = docker.ImageExport(imageId, archivePath)
				if err != nil {
					log.Errorf("Error exporting image: %v", err)
					return 1
				}
			}

			// Register tags for distribution
			for _, tag := range append(Tags, imageId) {
				err = store.RegisterDistribution(tag, archivePath)
				if err != nil {
					log.Errorf("Error registering distribution: %v", err)
					return 1
				}

				log.Infof("Registered tag %s for distribution", tag)
			}

			return 0
		}()
		os.Exit(code)
	},
}

func buildDockerImage(contextPath string, tag string) (io.ReadCloser, error) {
	dockerFile := DockerFile
	if DockerFile == "" {
		dockerFile = filepath.Join(contextPath, "Dockerfile")
	}

	buildArgs := map[string]*string{}
	for name, value := range BuildArgs {
		buildArgs[name] = &value
	}

	return docker.ImageBuild(contextPath, types.ImageBuildOptions{
		Dockerfile: filepath.ToSlash(dockerFile),
		PullParent: PullParent,
		NoCache:    DisableCache,
		BuildArgs:  buildArgs,
		Tags:       []string{tag},
	})
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVarP(&DockerFile, "file", "f", "", "Name of the Dockerfile (Default is 'PATH/Dockerfile').")
	buildCmd.Flags().StringArrayVarP(&Tags, "tag", "t", []string{}, "Name and optionally a tag in the 'name:tag' format.")
	buildCmd.Flags().StringToStringVar(&BuildArgs, "build-arg", map[string]string{}, "Set build-time variables.")
	buildCmd.Flags().BoolVar(&PullParent, "pull", true, "Always attempt to pull a newer version of the image.")
	buildCmd.Flags().BoolVar(&DisableCache, "no-cache", false, "Do not use cache when building the image.")
}
