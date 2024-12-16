package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func init() {
	rootCmd.AddCommand(add)
	rootCmd.AddCommand(open)
	rootCmd.AddCommand(write)
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

var open = &cobra.Command{
	Use:   "search",
	Short: "search a file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please enter a file name")
			return
		}
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Println("Error in the filepath: ", err)
			return
		}
		data, err := io.ReadAll(file)
		if err != nil {
			fmt.Println("Error reading the file: ", err)
			return
		}
		fmt.Println(string(data))
	},
}
