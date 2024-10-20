/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/labinc-nico/labinc-backuper/internal/config"
	"github.com/labinc-nico/labinc-backuper/internal/logger"
	"github.com/labinc-nico/labinc-backuper/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func init() {

	log = logger.GetLogger()
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "labinc-backuper",
	Short: "Simple backup tool for files and databases",
	Long: `A tool to compress files and dump databases simply. 
	It will create a tarball of any folder/file and a mysqldump of any database. 
	It will also upload to a remote server via sftp and remove the local file.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(appVersion string) {
	rootCmd.Version = appVersion
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initialize)
}
func initialize() {
	err := config.GetAppConfig()
	if err != nil {
		panic(err)
	}
	initApplication()
}

func initApplication() {
	root_dir := viper.GetString("backup.directory")
	dump_dir := filepath.Join(root_dir, viper.GetString("backup.dump_folder_name"))
	archive_dir := filepath.Join(root_dir, viper.GetString("backup.archive_folder_name"))
	log_dir := viper.GetString("logger.path")
	var directories = []string{dump_dir, archive_dir, log_dir}

	err := utils.CreateDirectory(directories)
	if err != nil {
		log.Fatalw("Error creating directories", err)
	}
}
