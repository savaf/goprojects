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
	listCmd.Flags().BoolVarP(&hardDelete, "hard", "f", false, "Show all tasks including completed")
	rootCmd.AddCommand(deleteCmd)
}

var hardDelete bool

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a task",
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

		// var myTask tasks.Task
		taskId, _ := strconv.ParseInt(args[0], 10, 0)
		var myTask *tasks.Task
		if hardDelete {
			myTask, _ = tasks.Delete(taskId)
		} else {
			myTask, _ = tasks.SoftDelete(taskId)
		}

		fmt.Printf("Task '%s' deleted!\n", myTask.Title)

		// Close connection
		db.CloseDB()

	},
}
