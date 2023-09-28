package runner

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/rulu158/gocoderunner/runner/languages"
)

const (
	imageBaseID      = "gocoderunner"
	dockerfilePrefix = "dockerfile_"
)

func getImageName() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	year, month, day := time.Now().Date()
	return fmt.Sprintf("%s-%v%v%v%v", imageBaseID, year, int(month), day, r.Uint32())
}

func getCodePath(id, folder string, lang languages.Language) string {
	return filepath.Join(folder, fmt.Sprintf("%s.%s", id, languages.LanguageExtensions[lang]))
}

func getDockerfilePath(id, folder string) string {
	return filepath.Join(folder, fmt.Sprintf("%s%s", dockerfilePrefix, id))
}
