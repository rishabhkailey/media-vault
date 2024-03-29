package app

import (
	cmdPackage "github.com/rishabhkailey/media-vault/internal/cmd"
	"github.com/rishabhkailey/media-vault/internal/config"
	"github.com/rishabhkailey/media-vault/internal/db"
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
	return cmdPackage.MeiliSearchMigrate(gormDb, meiliesearch, meiliSearchmigrateOptions.batchSize)
}

func init() {
	// migrate cmd
	meiliSearchmigrateCmd.PersistentFlags().IntVarP(&meiliSearchmigrateOptions.batchSize, "batch", "", 1000, "batch size")

	// sub commands of mieliesearch
	meiliSearchBaseCmd.AddCommand(meiliSearchmigrateCmd)

	// root cmd
	authServerCmd.AddCommand(meiliSearchBaseCmd)
}
