package runner

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/client"
	"github.com/rulu158/gocoderunner/runner/languages"
)

type Runner struct {
	ID             string
	Language       languages.Language
	Context        context.Context
	Client         *client.Client
	DockerfilePath string
	Code           []byte
	Options        RunnerOptions
}

func NewRunner(lang languages.Language, opts *RunnerOptions) *Runner {
	r := &Runner{
		ID:       getImageName(),
		Language: lang,
		Context:  context.Background(),
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	r.Client = cli

	r.Code = nil

	r.Options = getOptions(opts)

	r.DockerfilePath = getDockerfilePath(r.ID, DockerfilesFolder)

	return r
}

type RunnerOptions struct {
	DockerfileBasePath string
	Stdin              io.Reader
	Stdout             io.Writer
	Stderr             io.Writer
	Interactive        bool
	Timeout            time.Duration
}

var (
	DockerfilesFolder = filepath.Join(".", "dockerfiles")
)

func SetDockerfilesFolder(newDockerfilesFolder string) {
	DockerfilesFolder = newDockerfilesFolder
}

var DefaultRunnerOptions = RunnerOptions{
	DockerfileBasePath: filepath.Join(".", "dockerfile_bases", "Dockerfile_base"),
	Stdin:              os.Stdin,
	Stdout:             os.Stdout,
	Stderr:             os.Stderr,
	Interactive:        false,
	Timeout:            0,
}

func getOptions(opts *RunnerOptions) RunnerOptions {
	options := DefaultRunnerOptions
	if opts == nil {
		return options
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
	if opts.Interactive != false {
		options.Interactive = true
	}
	if opts.Timeout != 0 {
		options.Timeout = opts.Timeout
	}
	return options
}
