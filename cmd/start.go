package cmd

import (
	"github.com/spf13/cobra"
	"github.com/suaas21/library-management-api/database"
	"github.com/suaas21/library-management-api/server"
)

var svrPort string

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringVarP(&svrPort, "server-port", "", "4000", "This flag sets the serverPort of the server")
	startCmd.PersistentFlags().StringVarP(&database.Cfg.DBPort, "db-port", "", "5432", "This flag sets the Database port of the database server")
	startCmd.PersistentFlags().StringVarP(&database.Cfg.DBPassword, "db-password", "", "pass", "This flag sets for the Database password")
	startCmd.PersistentFlags().StringVarP(&database.Cfg.DBName, "db-name", "", "library_management", "This flag sets for the Database name")
	startCmd.PersistentFlags().StringVarP(&database.Cfg.BDUser, "db-user", "", "postgres", "This flag sets for the Database user")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "This command will start the api server",
	Long:  "This command will start the go-rest-api server",
	Run: func(cmd *cobra.Command, args []string) {
		server.StartAPIServer(svrPort)
	},
}