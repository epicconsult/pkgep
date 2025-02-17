package pkgep

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func GodotEnv(key string) string {

	if os.Getenv("GO_ENV") != "production" {
		godotenv.Load(filepath.Join(".env"))
		return os.Getenv(key)
	}
	return os.Getenv(key)
}
