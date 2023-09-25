package runner

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
)

func (r *Runner) CreateDockerfile() error {
	b, err := os.ReadFile(r.Options.DockerfileBasePath)
	if err != nil {
		return errors.New("Couldn't open base file")
	}

	dockerfileContents := bytes.ReplaceAll(
		b,
		[]byte("{executable}"),
		[]byte(r.ID),
	)

	err = os.WriteFile(r.DockerfilePath, dockerfileContents, 0644)
	if err != nil {
		return errors.New("Error writing to file")
	}

	return nil
}

func (r *Runner) BuildImage() error {
	tar, err := archive.TarWithOptions(".", &archive.TarOptions{})
	if err != nil {
		return err
	}

	opts := types.ImageBuildOptions{
		Dockerfile: r.DockerfilePath,
		Tags:       []string{r.ID},
		Remove:     true,
	}
	res, err := r.Client.ImageBuild(r.Context, tar, opts)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = r.printSyntaxErrors(res.Body, debug)
	if err != nil {
		return err
	}

	return nil
}

func (r *Runner) printSyntaxErrors(rd io.Reader, showLogs bool) error {
	var lastLine string

	var lines []string
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		lastLine = scanner.Text()
		if showLogs {
			fmt.Println(scanner.Text())
		}
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		if len(lines) <= 1 {
			return errors.New(errLine.Error)
		}
		syntaxErrors := strings.Split(lines[len(lines)-2], `\n`)
		syntaxErrors = syntaxErrors[1 : len(syntaxErrors)-1]
		for _, syntaxError := range syntaxErrors {
			fmt.Fprintln(r.Options.Stdout, syntaxError)
		}
		return errors.New("Code could not compile")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
