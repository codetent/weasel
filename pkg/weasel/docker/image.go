package docker

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/mholt/archiver/v3"
)

func ImageBuild(contextPath string, opts types.ImageBuildOptions) (io.ReadCloser, error) {
	// Create temporary directory for context archive
	tmpDir, err := ioutil.TempDir("", "context")
	if err != nil {
		return nil, fmt.Errorf("ImageBuild: TempDir(): %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Archive context directory
	contextFilePath := path.Join(tmpDir, "context.tar")
	err = archiver.Archive([]string{contextPath}, contextFilePath)
	if err != nil {
		return nil, fmt.Errorf("ImageBuild: Archive(): %v", err)
	}

	// Create docker client
	ctx := context.Background()
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("ImageBuild: NewClient(): %v", err)
	}

	// Prepare args
	opts.Dockerfile = filepath.ToSlash(opts.Dockerfile)

	contextFile, err := os.Open(contextFilePath)
	if err != nil {
		return nil, fmt.Errorf("ImageBuild: contextFile.Open(): %v", err)
	}
	defer contextFile.Close()

	// Build image
	resp, err := docker.ImageBuild(ctx, contextFile, opts)
	if err != nil {
		return nil, fmt.Errorf("ImageBuild: ImageBuild(): %v", err)
	}

	return resp.Body, nil
}

func ImagePull(tag string) (io.ReadCloser, error) {
	// Create docker client
	ctx := context.Background()
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("ImagePull: NewClient(): %v", err)
	}

	// Pull image from registry
	stream, err := docker.ImagePull(ctx, tag, types.ImagePullOptions{})
	if err != nil {
		return nil, fmt.Errorf("ImagePull: ImagePull(): %v", err)
	}

	return stream, nil
}

func ImageExport(tag string, targetPath string) error {
	// Create docker client
	ctx := context.Background()
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("ImageExport: NewClient(): %v", err)
	}

	// Create container running the image with the given tag
	container, err := docker.ContainerCreate(ctx, &container.Config{
		Image: tag,
		Cmd:   []string{""}, // This is require for creating a container if no command is specified
	}, nil, nil, nil, "")
	if err != nil {
		return fmt.Errorf("ImageExport: ContainerCreate(): %v", err)
	}

	// Read container filesystem & remove container
	reader, exportErr := docker.ContainerExport(ctx, container.ID)
	removeErr := docker.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
	if exportErr != nil {
		return fmt.Errorf("ImageExport: ContainerExport(): %v", exportErr)
	}
	defer reader.Close()
	if removeErr != nil {
		return fmt.Errorf("ImageExport: ContainerRemove(): %v", removeErr)
	}

	// Write filesystem to target archive
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("ImageExport: targetFile.Create(): %v", err)
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, reader)
	if err != nil {
		return fmt.Errorf("ImageExport: targetFile.Copy(): %v", err)
	}

	return nil
}

func ImageIdByTag(tag string) (string, error) {
	// Create docker client
	ctx := context.Background()
	docker, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", fmt.Errorf("ImagePull: NewClient(): %v", err)
	}

	// Get id of tagged image
	filters := filters.NewArgs()
	filters.Add("reference", tag)
	images, err := docker.ImageList(ctx, types.ImageListOptions{
		Filters: filters,
	})
	if err != nil {
		return "", fmt.Errorf("ImagePull: ImageList(): %v", err)
	}
	if len(images) == 0 {
		return "", fmt.Errorf("ImagePull: tag not found")
	}

	idValue := strings.SplitN(images[0].ID, ":", 2)
	return idValue[1][:12], nil
}
