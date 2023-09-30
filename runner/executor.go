package runner

import (
	"context"
	"errors"
	"log"
	"os/exec"
)

const (
	debug = false
)

var (
	UnrecoverableError = errors.New("Server error")
)

func (r *Runner) ExecCode(code []byte) error {
	defer r.FreeResources()

	err := r.CreateCodeFileFromBytes(code)
	if err != nil {
		log.Println("CREATE_CODE_FILE: " + err.Error())
		return UnrecoverableError
	}

	err = r.CreateDockerfile()
	if err != nil {
		log.Println("CREATE_DOCKERFILE: " + err.Error())
		return UnrecoverableError
	}

	err = r.BuildImage()
	if err != nil {
		log.Println("BUILD_IMAGE: " + err.Error())
		return err
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
		log.Println("INITIALIZE_CONTAINER: " + err.Error())
		return UnrecoverableError
	}

	return nil
}

func (r *Runner) InitializeContainer() error {
	command := ""
	if r.Options.Interactive {
		command = "./run_docker_it.sh"
	} else {
		command = "./run_docker.sh"
	}

	var cmd *exec.Cmd
	if r.Options.Timeout != 0 {
		ctx, cancel := context.WithTimeout(context.Background(), r.Options.Timeout)
		defer cancel()
		cmd = exec.CommandContext(ctx, command, r.ID)
	} else {
		cmd = exec.Command(command, r.ID)
	}

	cmd.Stdin = r.Options.Stdin
	cmd.Stdout = r.Options.Stdout
	cmd.Stderr = r.Options.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
