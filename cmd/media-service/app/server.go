package app

import (
	"os"

	"github.com/rishabhkailey/media-service/internal/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type authServerCmdOptions struct {
	config string
}

var authServerOptions authServerCmdOptions

// authServerCmd represents the base command when called without any subcommands
var authServerCmd = &cobra.Command{
	Use: "media-service",
	Long: `The Auth Server can be used by clients as authentication Server.
It supports Authorization Code Grant flow.
	`,
	RunE: startServer,
}

func startServer(cmd *cobra.Command, args []string) error {
	return api.Start()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the authServerCmd.
func Execute() {
	err := authServerCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	authServerCmd.PersistentFlags().StringVarP(&authServerOptions.config, "config", "", "configs/authservice.yaml", "path of the config file")
	viper.BindPFlag("config", authServerCmd.PersistentFlags().Lookup("config"))
}
