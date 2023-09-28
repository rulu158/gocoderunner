package runner

import (
	"context"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func (r *Runner) FreeResources() {
	imagesToRemove := []string{r.ID}
	containers, _ := r.Client.ContainerList(r.Context, types.ContainerListOptions{})
	for _, container := range containers {
		if container.Command == "/gocoderunner/"+r.ID {
			imagesToRemove = append(imagesToRemove, container.ImageID)
			if container.State == "running" {
				timeout := 0
				err := r.Client.ContainerStop(r.Context, container.ID, containertypes.StopOptions{Timeout: &timeout})
				if err != nil {
					log.Println(errors.Join(errors.New("Error stopping container "+container.ID+":"), err))
				}
			}
		}
	}

	for _, imageID := range imagesToRemove {
		if _, err := r.Client.ImageRemove(r.Context, imageID, types.ImageRemoveOptions{PruneChildren: true, Force: true}); err != nil {
			log.Println(errors.Join(errors.New("Error removing image "+imageID+":"), err))
		}
	}

	err := os.Remove(r.DockerfilePath)
	if err != nil {
		log.Println(errors.Join(errors.New("Error removing dockerfile "+r.DockerfilePath+":"), err))
	}

	err = os.Remove(r.CodePath)
	if err != nil {
		log.Println(errors.Join(errors.New("Error removing code file "+r.CodePath+":"), err))
	}
}

func FreeAllResources() {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	imagesToRemove := []string{}
	containers, _ := cli.ContainerList(ctx, types.ContainerListOptions{})
	for _, container := range containers {
		if strings.HasPrefix(container.Command, "/gocoderunner/") {
			imagesToRemove = append(imagesToRemove, container.ImageID)
			if container.State == "running" {
				timeout := 0
				err := cli.ContainerStop(ctx, container.ID, containertypes.StopOptions{Timeout: &timeout})
				if err != nil {
					log.Println(errors.Join(errors.New("Error stopping container "+container.ID+":"), err))
				}
			}
		}
	}

	for _, imageID := range imagesToRemove {
		if _, err := cli.ImageRemove(ctx, imageID, types.ImageRemoveOptions{PruneChildren: true, Force: true}); err != nil {
			log.Println(errors.Join(errors.New("Error removing image "+imageID+":"), err))
		}
	}

	cmd := exec.Command("./run_docker_prune.sh")
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}

	items, err := os.ReadDir(DockerfilesFolder)
	for _, item := range items {
		if !item.IsDir() && strings.HasPrefix(item.Name(), dockerfilePrefix) {
			err := os.Remove(filepath.Join(DockerfilesFolder, item.Name()))
			if err != nil {
				log.Println(errors.Join(errors.New("Error removing dockerfile "+filepath.Join(DockerfilesFolder, item.Name())+":"), err))
			}
		}
	}

	items, err = os.ReadDir(CodeFolder)
	for _, item := range items {
		if !item.IsDir() && strings.HasPrefix(item.Name(), imageBaseID) {
			err := os.Remove(filepath.Join(CodeFolder, item.Name()))
			if err != nil {
				log.Println(errors.Join(errors.New("Error removing code file "+filepath.Join(CodeFolder, item.Name())+":"), err))
			}
		}
	}
}
