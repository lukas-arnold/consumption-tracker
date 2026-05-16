package configs

import (
	"io/fs"
	"os"

	"github.com/lukas-arnold/consumption-tracker/web"
)

const port = ":8080"
const storageFile = "data/consumption-tracker.json"
const language = "en"

func GetPort() string {
	portEnv := os.Getenv("PORT")
	if portEnv == "" {
		return port
	}
	return portEnv
}

func GetStorageFile() string {
	storageFileEnv := os.Getenv("STORAGE_FILE")
	if storageFileEnv == "" {
		return storageFile
	}
	return storageFileEnv
}

func GetLanguage() string {
	languageEnv := os.Getenv("LANGUAGE")
	if languageEnv == "" {
		return language
	}
	return languageEnv
}

func GetWebFiles() fs.FS {
	return web.WebFiles
}
