package dumper

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/JamesStewy/go-mysqldump"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labinc-nico/labinc-backuper/internal/utils"
	"github.com/spf13/viper"
)

// BackupDB backups the given database to a file
//
// The function assumes that the environment variables
// DB_BACKUP_USER, DB_BACKUP_PWD, DB_BACKUP_HOSTNAME, DB_BACKUP_PORT and DB_BACKUP_DIR
// are defined. If any of these variables are not defined, the function will panic.
//
// The backups are stored in a directory structure like this:
// dumpDir/dbname-YYYYMMDDTHHMMSS.sql
//
// If the dumpDir directory does not exist, it will be created.
func Dump(dbname string) (string, error) {
	dumpDir := filepath.Join(viper.GetString("backup.directory"), viper.GetString("backup.dump_folder_name"))

	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", dbname)

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			viper.GetString("database.user"),
			viper.GetString("database.password"),
			viper.GetString("localhost"),
			"3306",
			dbname))
	if err != nil {
		return "", fmt.Errorf("error opening database: %e", err)
	}

	dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
	if err != nil {
		return "", fmt.Errorf("error registering database: %e", err)
	}

	resultFilename, err := dumper.Dump()
	if err != nil {
		os.Remove(resultFilename)
		return "", fmt.Errorf("error dumping: %e", err)
	}

	dumper.Close()
	utils.Clean(dumpDir, dbname, viper.GetInt("backup.max_files_to_keep"))
	return resultFilename, nil
}
