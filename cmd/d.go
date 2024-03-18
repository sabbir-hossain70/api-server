/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dCmd represents the d command
var (
	Val  string
	dCmd = &cobra.Command{
		Use:   "d",
		Short: "A brief description of your command",
		Long:  `long desc`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("d called")
			println(Val)
		},
	}
)

func init() {
	aCmd.AddCommand(dCmd)
	dCmd.PersistentFlags().StringVarP(&Val, "value", "x", "SabbirHossain", "This is to print")
}
