package utils

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/labinc-nico/labinc-backuper/internal/logger"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func init() {

	log = logger.GetLogger()
}
func CreateDirectory(directories []string) error {
	for _, directory := range directories {
		err := os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

type FileInfo struct {
	os.FileInfo
	path string
}

func getFiles(dir string, prefix string) ([]FileInfo, error) {
	var files []FileInfo
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), prefix) {
			files = append(files, FileInfo{info, path})
		}
		return nil
	})
	return files, err
}
func Clean(folder string, prefix string, filesToKeep int) {
	files, err := getFiles(folder, prefix)
	if err != nil {
		log.Warnf("error while getting files : %s", err)
		return
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	for i := 0; i < len(files)-filesToKeep; i++ {
		log.Infof("Cleanup folder -- suppr fichier : %s", files[i].path)
		err := os.Remove(files[i].path)
		if err != nil {
			log.Errorf("error while removing file %s", err)
		}
	}
}
