/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/manifoldco/promptui"
	"github.com/ospfarea0/quiz/game"
	"github.com/spf13/cobra"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "play starts a user for a user",
	Long:  `play starts a quiz for a user`,
	Run: func(cmd *cobra.Command, args []string) {

		prompt := promptui.Prompt{
			Label:     "Enter your preferred username",
			Default:   "yourUserName",
			AllowEdit: true,
		}

		result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		client := http.Client{}
		resp, err := client.Get("http://localhost:8080/getgameids")
		if err != nil {
			slog.Error("failed to get game ids", "error", err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("failed to get game ids", "error", err)
		}

		ids := make([]int64, 1)
		err = json.Unmarshal(body, &ids)
		if err != nil {
			slog.Error("error unmarshaling json", "error", err)
		}

		player := struct {
			Name   string `json:"name"`
			QuizID int64  `json:"quiz_id"`
		}{}

		player.Name = result
		player.QuizID = ids[0]
		data, err := json.Marshal(player)
		if err != nil {
			slog.Error("failed to marshal player into json", "error", err)
		}

		_, err = client.Post("http://localhost:8080/registerplayer", "application/json", bytes.NewReader(data))
		if err != nil {
			slog.Error("failed to post player", "error", err)
		}

		resp, err = client.Post("http://localhost:8080/play", "application/json", bytes.NewReader(data))
		if err != nil {
			slog.Error("failed to post to play", "error", err, "postData", data)
		}

		body, err = io.ReadAll(resp.Body)
		defer resp.Body.Close()

		if err != nil {
			slog.Error("failed to read response body in play", "error", err)
		}

		questions := make([]game.Question, 5)
		err = json.Unmarshal(body, &questions)
		if err != nil {
			slog.Error("failed to unmarshal questions", "error", err)
		}
		fmt.Printf("Here are the questions %v", questions)
		//answers := make([]game.Question, 5)
		result = ""
		for i := range questions {
			prompt := promptui.Select{
				Label: questions[i].Question,
				Items: []string{questions[i].Answers.Answer1.Answer, questions[i].Answers.Answer2.Answer,
					questions[i].Answers.Answer3.Answer, questions[i].Answers.Answer4.Answer},
			}

			_, result, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
			// answer := game.Question{Question: questions[i].Question,
			// 	Answers: game.Answers{Answer1: questions[i].Answers.Answer1, Answer2: questions[i].Answers.Answer2,
			// 		Answer3: questions[i].Answers.Answer3, Answer4: questions[i].Answers.Answer4}}

			//answers := append(answers, answer)
			buildAnswers(questions, result)
			fmt.Printf("You choose %q\n", result)

		}
		fmt.Printf("these are the final answers %v", questions)
	},
}

func init() {
	rootCmd.AddCommand(playCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func buildAnswers(questions []game.Question, result string) []game.Question {
	for i := range questions {
		ans := questions[i].GetCorrectAnswer()
		if ans == result {
			questions[i].IsRight = true
		}
	}
	return questions
}
