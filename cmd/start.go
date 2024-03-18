/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/sabbir-hossain70/api-server/apiHandler"
	"github.com/spf13/cobra"
	"strconv"
)

// startCmd represents the start command
var (
	// Port stores port number for starting a connection
	Port     int
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "start cmd starts the server on a port",
		Long: `It starts the server on a given port number, 
				Port number will be given in the cmd`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(args)
			fmt.Println("Port: ", Port)
			apiHandler.RunServer(strconv.Itoa(Port))
		},
	}
)

func init() {
	//fmt.Println("Here")
	startCmd.PersistentFlags().IntVarP(&Port, "port", "p", 8080, "Port number for starting server")
	rootCmd.AddCommand(startCmd)
}
