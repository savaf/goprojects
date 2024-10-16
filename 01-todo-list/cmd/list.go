package cmd

import (
	"fmt"
	"log"
	"os"
	"savaf/todo-list/db"
	"savaf/todo-list/models/tasks"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all tasks including completed")
	rootCmd.AddCommand(listCmd)
}

var showAll bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
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

		var myTasks tasks.Tasks
		if showAll {
			myTasks, _ = tasks.ShowAll()
		} else {
			myTasks, _ = tasks.ShowPending()
		}

		db.CloseDB()

		printTasksTable(myTasks, showAll)
	},
}

// Function that prints the tasks table
func printTasksTable(tasks tasks.Tasks, showAll bool) {
	// Create a writer to handle the formatted output with tabs
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.Debug)

	// Print the table header
	headers := "ID\tDescription\tCreated At"
	if showAll {
		headers = headers + "\tDone"
	}
	fmt.Fprintln(w, headers)

	// Iterate over the tasks and display them
	for _, task := range tasks {
		row := task.ToRow()

		if showAll {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", task.Id, task.Title, row[len(row)-2], row[len(row)-1])
		} else {
			fmt.Fprintf(w, "%d\t%s\t%s\n", task.Id, task.Title, row[len(row)-2])
		}
	}

	// Write the formatted content to the terminal
	w.Flush()
}
