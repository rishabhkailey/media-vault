package app

import (
	"github.com/rishabhkailey/media-service/internal/config"
	"github.com/rishabhkailey/media-service/internal/db"
	"github.com/rishabhkailey/media-service/internal/utils"
	"github.com/spf13/cobra"
)

type meiliSearchBaseCmdOptions struct {
}

var meiliSearchBaseOptions meiliSearchBaseCmdOptions

// meiliSearchBaseCmd represents the base command when called without any subcommands
var meiliSearchBaseCmd = &cobra.Command{
	Use:   "meiliesearch",
	Short: "manage meiliesearch data",
}

type meiliSearchmigrateCmdOptions struct {
	batchSize int
}

var meiliSearchmigrateOptions meiliSearchmigrateCmdOptions

// meiliSearchBaseCmd represents the base command when called without any subcommands
var meiliSearchmigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate data from SQL DB to meiliesearch",
	RunE:  executeMelieSearchMigrate,
}

func executeMelieSearchMigrate(cmd *cobra.Command, args []string) error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}
	gormDb, err := db.NewGoOrmConnection(config.Database)
	if err != nil {
		return err
	}
	meiliesearch, err := db.NewMeiliSearchClient(config.MeiliSearch)
	if err != nil {
		return err
	}
	return utils.MeiliSearchMigrate(gormDb, meiliesearch, meiliSearchmigrateOptions.batchSize)
}

func init() {
	// migrate cmd
	meiliSearchmigrateCmd.PersistentFlags().IntVarP(&meiliSearchmigrateOptions.batchSize, "batch", "", 1000, "batch size")

	// sub commands of mieliesearch
	meiliSearchBaseCmd.AddCommand(meiliSearchmigrateCmd)

	// root cmd
	authServerCmd.AddCommand(meiliSearchBaseCmd)
}
