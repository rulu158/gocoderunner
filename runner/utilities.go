package runner

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"
)

const (
	imageBaseID = "gocoderunner"
)

func getImageName() string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	year, month, day := time.Now().Date()
	return fmt.Sprintf("%s-%v%v%v%v", imageBaseID, year, int(month), day, r.Uint32())
}

func getDockerfilePath(id, folder string) string {
	return filepath.Join(folder, fmt.Sprintf("dockerfile_%s", id))
}
