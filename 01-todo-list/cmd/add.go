package cmd

import (
	"fmt"
	"log"
	"savaf/todo-list/db"
	"savaf/todo-list/models/tasks"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task",
	Run: func(cmd *cobra.Command, args []string) {
		// Create an exported global variable to hold the database connection pool.
		myDb, err := db.ConnectToDB("./tasks.db")
		if err != nil {
			log.Fatal("Error connecting database:", err)
		}

		// Initialize the SQLite database
		_, err = tasks.InitializeDB(myDb)
		if err != nil {
			log.Fatal("Error initializing database:", err)
		}

		title := args[0]
		tasks.Add(title)
		fmt.Printf("Task '%s' added!\n", title)

		// Close connection
		db.CloseDB()

	},
}
