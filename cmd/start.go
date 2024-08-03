/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/ospfarea0/quiz/game"
	"github.com/spf13/cobra"
)

// playCmd represents the play command
var startCmd = &cobra.Command{
	Use:   "start &",
	Short: "Start starts the quiz game",
	Long:  `Start is used to start the quiz game service in the background`,
	Run: func(cmd *cobra.Command, args []string) {
		game.StartService()
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
