/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	dumper "github.com/labinc-nico/labinc-backuper/internal/dump"
	"github.com/spf13/cobra"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump databases",
	Long:  `A tool to dump databases.`,
	Run: func(cmd *cobra.Command, args []string) {
		dumpFile, err := dumper.Dump(dbName)
		if err != nil {
			log.Fatalw("Error dumping database", err)
		}

		log.Infof("Dumped database %s to %s", dbName, dumpFile)
	},
}
var dbName string

func init() {
	rootCmd.AddCommand(dumpCmd)
	dumpCmd.PersistentFlags().StringVarP(&dbName, "database", "n", "", "name of database to dump")
	dumpCmd.MarkPersistentFlagRequired("database")
}
