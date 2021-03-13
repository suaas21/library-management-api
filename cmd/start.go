package cmd

import (
	"github.com/spf13/cobra"
	"github.com/suaas21/library-management-api/pkg/db"
	"github.com/suaas21/library-management-api/pkg/server"
)

var (
	ServerPort string
	dbPort     string
	dbPassword string
	dbName     string
	dbUser     string
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringVarP(&ServerPort, "Server-port", "", "4000", "This flag sets the ServerPort of the server")
	startCmd.PersistentFlags().StringVarP(&dbPort, "db-port", "", "5432", "This flag sets the Database port of the database server")
	startCmd.PersistentFlags().StringVarP(&dbPassword, "db-password", "", "pass", "This flag sets for the Database password")
	startCmd.PersistentFlags().StringVarP(&dbName, "db-name", "", "library_management", "This flag sets for the Database name")
	startCmd.PersistentFlags().StringVarP(&dbUser, "db-user", "", "postgres", "This flag sets for the Database user")

	db.InitializeDB(dbUser, dbPassword, dbName, dbPort)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "This command will start the api server",
	Long:  "This command will start the go-rest-api server",
	Run: func(cmd *cobra.Command, args []string) {
		server.StartAPIServer(ServerPort)
	},
}