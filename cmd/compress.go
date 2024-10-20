/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/labinc-nico/labinc-backuper/internal/archive"
	"github.com/spf13/cobra"
)

var folderName string

// archiveCmd represents the archive command
var compressCmd = &cobra.Command{
	Use:   "compress",
	Short: "Compress a folder to archive",
	Long:  `A tool to compress a folder to archive.`,
	Run: func(cmd *cobra.Command, args []string) {
		compressedFile, err := archive.BackupFiles(folderName)
		if err != nil {
			log.Fatalw("Error compressing folder", err)
		}
		log.Infof("Compress folder %s to %s", dbName, compressedFile)
	},
}

func init() {
	rootCmd.AddCommand(compressCmd)
	compressCmd.PersistentFlags().StringVarP(&folderName, "folder", "n", "", "name of folder to compress")
	compressCmd.MarkPersistentFlagRequired("folder")
}
