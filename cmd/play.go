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

		ids := make([]int64, 5)
		err = json.Unmarshal(body, &ids)
		if err != nil {
			slog.Error("error unmarshaling json", "error", err)
		}

		player := game.Player{}

		player.Name = result
		player.QuizId = ids[0]
		data, err := json.Marshal(player)
		if err != nil {
			slog.Error("failed to marshal player into json", "error", err)
		}

		resp, err = client.Post("http://localhost:8080/registerplayer", "application/json", bytes.NewReader(data))
		if err != nil {
			slog.Error("failed to post player", "error", err)
		}

		body, err = io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			slog.Error("failed to read response body after registering player", "error", err)
		}
		updatedPlayer := game.Player{}
		err = json.Unmarshal(body, &updatedPlayer)
		if err != nil {
			slog.Error("failed to unmarshal into updateplayer", "error", err)
		}

		resp, err = client.Post("http://localhost:8080/play", "application/json", bytes.NewReader(body))
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

			buildAnswers(questions, result)

			fmt.Printf("You choose %q\n", result)

		}
		updatedPlayer.Answers = questions
		data, err = json.Marshal(updatedPlayer)
		if err != nil {
			slog.Error("failed to marshal player", "error", err)
			return
		}

		resp, err = client.Post("http://localhost:8080/submitanswers", "application/json", bytes.NewReader(data))
		if err != nil {
			slog.Error("failed to post to submitanswers", "error", err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			slog.Error("failed to post to submitanswers", "error", err, "http code", resp.StatusCode)
			return
		}

		resp, err = client.Post("http://localhost:8080/getgrade", "application/json", bytes.NewReader(data))
		if err != nil {
			slog.Error("failed tro getgrade", "error", err)
			return
		}
		defer resp.Body.Close()
		gradeBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("failed to read response body while getting grade", "error", err)
		}

		s := struct {
			Score float32 `json:"score"`
		}{}

		err = json.Unmarshal(gradeBody, &s)
		if err != nil {
			slog.Error("failed to unmarshal into score", "error", err)
		}

		resp, err = client.Post("http://localhost:8080/getpercentile", "application/json", bytes.NewReader(data))
		if err != nil {
			slog.Error("failed to get percentile", "error", err)
			return
		}
		defer resp.Body.Close()
		percentileBody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("failed to read response body while getting grade", "error", err)
		}

		percentile := struct {
			Percentile float32 `json:"percentile"`
		}{}

		err = json.Unmarshal(percentileBody, &percentile)
		if err != nil {
			slog.Error("failed to unmarshal into percentile", "error", err)
		}

		fmt.Printf("Your score is: %2f\n", s.Score)
		fmt.Printf("Your percentile is: %2f\n", percentile.Percentile)
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
