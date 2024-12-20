package cmd

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func init() {
	rootCmd.AddCommand(add)
	rootCmd.AddCommand(write)
	rootCmd.AddCommand(del)
	rootCmd.AddCommand(copy)
	rootCmd.AddCommand(edit)

	copy.Flags().BoolP("append", "a", false, "appends into the file instead of truncating the file")
}

var copy = &cobra.Command{
	Use:   "copy",
	Short: "Copy files",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Not enough arguments")
			return
		}
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		appen, _ := cmd.Flags().GetBool("append")
		if appen {
			file2, err := os.OpenFile(args[1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file2.Close()
			_, err = file2.WriteString(string(content))
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			file2, err := os.Create(args[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file2.Close()
			_, err = file2.WriteString(string(content))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		fmt.Println("Copied the file", args[0], "successfully into", args[1])
	},
}

var write = &cobra.Command{
	Use:   "write",
	Short: "write to file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("filename is required")
			return
		}
		file, err := os.OpenFile(args[0], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		if args[1] != "" {
			_, err := file.WriteString(args[1] + "\n")
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		fmt.Println("writen to file successfully")

	},
}

var del = &cobra.Command{
	Use:   "delete",
	Short: "delete specified file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("filename is required")
			return
		}
		_, err := os.Stat(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = os.Remove(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

var add = &cobra.Command{
	Use:   "add",
	Short: "add a file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please enter a file name")
			return
		}
		file, err := os.OpenFile(args[0], os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error adding the file: ", err)
			return
		}
		defer file.Close()
		fmt.Println("Added file:", file.Name())
	},
}

var edit = &cobra.Command{
	Use:   "edit",
	Short: "opens editor with file name",
	Run:   runner,
}

func runner(cmd *cobra.Command, args []string) {
	app := tview.NewApplication()
	var err error
	var blob = InitBlob(args[0])
	var ed = InitDisplay(app, blob)
	ed.Render()
	err = app.SetRoot(ed.layout, true).Run()
	if err != nil {
		fmt.Println("Error setting root app:", err)
		return
	}

}
