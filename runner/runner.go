package runner

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
)

type Runner struct {
	ID             string
	Context        context.Context
	Client         *client.Client
	DockerfilePath string
	Options        RunnerOptions
}

func NewRunner(opts *RunnerOptions) *Runner {
	r := &Runner{
		ID:      getImageName(),
		Context: context.Background(),
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	r.Client = cli

	r.Options = getOptions(opts)

	r.DockerfilePath = getDockerfilePath(r.ID, r.Options.DockerfilesFolder)

	return r
}

type RunnerOptions struct {
	DockerfilesFolder  string
	DockerfileBasePath string
	Stdin              io.Reader
	Stdout             io.Writer
	Stderr             io.Writer
}

var DefaultRunnerOptions = RunnerOptions{
	DockerfilesFolder:  filepath.Join(".", "dockerfiles"),
	DockerfileBasePath: filepath.Join(".", "dockerfile_bases", "Dockerfile_base"),
	Stdin:              os.Stdin,
	Stdout:             os.Stdout,
	Stderr:             os.Stderr,
}

func getOptions(opts *RunnerOptions) RunnerOptions {
	options := DefaultRunnerOptions
	if opts == nil {
		return options
	}
	if opts.DockerfilesFolder != "" {
		options.DockerfilesFolder = opts.DockerfilesFolder
	}
	if opts.DockerfileBasePath != "" {
		options.DockerfileBasePath = opts.DockerfileBasePath
	}
	if opts.Stdin != nil {
		options.Stdin = opts.Stdin
	}
	if opts.Stdout != nil {
		options.Stdout = opts.Stdout
	}
	if opts.Stderr != nil {
		options.Stderr = opts.Stderr
	}
	return options
}
