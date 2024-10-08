package physical

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/arpanrec/secretsquirrel/internal/appconfig"
)

var mutexPhysicalFile = &sync.Mutex{}

var oncePhysicalFile = &sync.Once{}

type FileStorageConfig struct {
	Path string `json:"path"`
}

var fileStorageConfigVar FileStorageConfig

func getFileStorageConfigVar() FileStorageConfig {
	oncePhysicalFile.Do(func() {
		storagePath := appconfig.GetConfig().Storage.Config["path"].(string)
		if storagePath == "" {
			log.Fatalln("Fatal Storage path not set")
		}
		fileStorageConfigVar = FileStorageConfig{
			Path: storagePath,
		}
		log.Printf("File storage path set to %v", fileStorageConfigVar)
	})
	return fileStorageConfigVar
}

func (fs FileStorageConfig) GetData(Location string) (string, error) {
	p := path.Join(getFileStorageConfigVar().Path, Location)
	d, err := os.ReadFile(p)
	return string(d), err
}

func (fs FileStorageConfig) PutData(Location string, Data string) (bool, error) {
	mutexPhysicalFile.Lock()
	defer mutexPhysicalFile.Unlock()
	p := path.Join(getFileStorageConfigVar().Path, Location)
	dir := filepath.Dir(p)
	errMakeDir := os.MkdirAll(dir, 0755)
	if errMakeDir != nil {
		log.Println("Error creating directory: ", errMakeDir)
		return false, errMakeDir
	}
	errWriteFile := os.WriteFile(p, []byte(Data), 0644)
	if errWriteFile != nil {
		log.Println("Error writing file: ", errWriteFile)
		return false, errWriteFile
	}
	return true, nil
}

func (fs FileStorageConfig) DeleteData(Location string) error {
	mutexPhysicalFile.Lock()
	defer mutexPhysicalFile.Unlock()
	p := path.Join(getFileStorageConfigVar().Path, Location)
	err := os.Remove(p)
	return err
}
