/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/spf13/cobra"
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "A dump of the current quiz state",
	Long:  `A dump for troubleshooting`,
	Run: func(cmd *cobra.Command, args []string) {
		client := http.Client{}
		resp, err := client.Get("http://localhost:8080/printquiz")
		if err != nil {
			slog.Error("failed to get game ids", "error", err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("failed to read resp body", "error", err)
			return
		}
		defer resp.Body.Close()

		fmt.Print(string(body))

	},
}

func init() {
	rootCmd.AddCommand(printCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// printCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// printCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
