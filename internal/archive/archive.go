package archive

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/labinc-nico/labinc-backuper/internal/utils"
	"github.com/mholt/archiver/v4"
	"github.com/spf13/viper"
)

// BackupFiles backups the given folder to a file
//
// The function assumes that the environment variables DB_BACKUP_DIR
// is defined. If the variable is not defined, the function will panic.
//
// The backups are stored in a directory structure like this:
// dumpDir/foldername-YYYYMMDDTHHMMSS.tar.gz
//
// If the dumpDir directory does not exist, it will be created.
func BackupFiles(folder string) (string, error) {
	root_dir := viper.GetString("backup.directory")

	backupDir := filepath.Join(root_dir, viper.GetString("backup.archive_folder_name"))

	files, err := archiver.FilesFromDisk(nil, map[string]string{
		folder: "", // contents added recursively
	})
	if err != nil {
		return "", fmt.Errorf("error reading files from disk: %v", err)
	}
	t := time.Now()
	archiveFilenameFormat := filepath.Base(folder) + "-" + t.Format("20060102T150405") + ".tar.gz"
	compressedFile := filepath.Join(backupDir, archiveFilenameFormat)
	out, err := os.Create(compressedFile)
	if err != nil {
		return "", fmt.Errorf("errorrror creating output file: %v", err)
	}
	defer out.Close()

	format := archiver.CompressedArchive{
		Compression: archiver.Gz{},
		Archival:    archiver.Tar{},
	}
	// create the archive
	err = format.Archive(context.Background(), out, files)
	if err != nil {
		return "", fmt.Errorf("error creating archive: %v", err)
	}
	utils.Clean(backupDir, filepath.Base(folder), viper.GetInt("backup.max_files_to_keep"))
	return compressedFile, nil
}
