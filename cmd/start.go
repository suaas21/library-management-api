package cmd

import (
	"github.com/spf13/cobra"
	"github.com/suaas21/library-management-api/server"
)

var (
	serverPort string
	dbPort     string
	dbPassword string
	dbName     string
	dbUser     string
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringVarP(&serverPort, "server-port", "", "4000", "This flag sets the serverPort of the server")
	startCmd.PersistentFlags().StringVarP(&dbPort, "db-port", "", "5432", "This flag sets the Database port of the database server")
	startCmd.PersistentFlags().StringVarP(&dbPassword, "db-password", "", "pass", "This flag sets for the Database password")
	startCmd.PersistentFlags().StringVarP(&dbName, "db-name", "", "library_management", "This flag sets for the Database name")
	startCmd.PersistentFlags().StringVarP(&dbUser, "db-user", "", "postgres", "This flag sets for the Database user")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "This command will start the api server",
	Long:  "This command will start the go-rest-api server",
	Run: func(cmd *cobra.Command, args []string) {
		svr := server.Server{
			ServerPort: serverPort,
			DBPort:     dbPort,
			DBPassword: dbPassword,
			DBName:     dbName,
			DBUser:     dbUser,
		}
		svr.StartAPIServer()
	},
}