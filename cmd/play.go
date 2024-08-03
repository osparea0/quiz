/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/ospfarea0/quiz/game"
	"github.com/spf13/cobra"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Play starts the quiz game",
	Long:  `Play is used to start the quiz game and presents the player with questions.`,
	Run: func(cmd *cobra.Command, args []string) {
		game.StartService()
		func() {
			client := &http.Client{}
			resp, err := client.Get("http://localhost:8080/getgameids")
			if err != nil {
				slog.Error("failed to call getgameids in play command", "error", err)
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}

			var data []map[string]interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				fmt.Println("Error unmarshaling JSON:", err)
				return
			}
			reader := bufio.NewReader(os.Stdin)
			for _, obj := range data {
				fmt.Println(obj)
				fmt.Println("Press 'n' to continue to the next item...")
				for {
					input, _ := reader.ReadString('\n')
					if input != "\n" {
						break
					}
					fmt.Println("Invalid input. Press 'n' to continue.")
				}
			}

		}()
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
