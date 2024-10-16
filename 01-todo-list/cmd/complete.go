package cmd

import (
	"fmt"
	"log"
	"savaf/todo-list/db"
	"savaf/todo-list/models/tasks"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(completeCmd)
}

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Complete a task",
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

		taskId, _ := strconv.ParseInt(args[0], 10, 0)
		task, _ := tasks.Toggle(taskId)
		fmt.Printf("Task '%s' completed!\n", task.Title)

		// Close connection
		db.CloseDB()

	},
}
