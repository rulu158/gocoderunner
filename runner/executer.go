package runner

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
)

const (
	timeout = 15 * time.Second
	debug   = false
)

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func handleSignals(fn func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for sig := range c {
			if sig == os.Interrupt || sig == syscall.SIGTERM {
				fn()
				os.Exit(1)
			}
		}
	}()
}

func (r *Runner) ExecCode() {
	handleSignals(func() { r.freeResources() })
	defer r.freeResources()

	err := r.CreateDockerfile()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = r.BuildImage()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	/*
		imagesList, _ := cli.ImageList(ctx, types.ImageListOptions{})
		var image types.ImageSummary
		for _, img := range imagesList {
			if len(image.RepoTags) > 0 && strings.Contains(image.RepoTags[0], imageName) {
				image = img
			}
		}
	*/
	err = r.InitializeContainer()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (r *Runner) InitializeContainer() error {
	cmd := exec.Command("./run_docker_it.sh", r.ID)

	cmd.Stdin = r.Options.Stdin
	cmd.Stdout = r.Options.Stdout
	cmd.Stderr = r.Options.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (r *Runner) freeResources() {
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

	/*
		cmd := exec.Command("./run_docker_prune.sh")
		if err := cmd.Run(); err != nil {
			log.Println(err)
		}
	*/

	err := os.Remove(r.DockerfilePath)
	if err != nil {
		log.Println(errors.Join(errors.New("Error removing dockerfile "+r.DockerfilePath+":"), err))
	}
}
