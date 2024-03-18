/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// aCmd represents the a command
var aCmd = &cobra.Command{
	Use:   "a",
	Short: "A brief description of your command",
	Long:  `A longer description of A`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("a called")
	},
}

func init() {
	rootCmd.AddCommand(aCmd)
	/*err := aCmd.Execute()
	if err != nil {
		panic(err)
	}*/

}
