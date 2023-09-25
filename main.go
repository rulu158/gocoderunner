package main

import "github.com/rulu158/gocoderunner/runner"

func main() {
	runner := *runner.NewRunner(&runner.RunnerOptions{})
	runner.ExecCode()
}
