package main

import (
	"errors"
	"fmt"
	"log/slog"
)

type Game struct {
	Quizzes []Quiz
	logger  *slog.Logger
}

func NewGame(numberOfQuizzes int) (Game, error) {
	if numberOfQuizzes <= 0 {
		return Game{}, errors.New("number of quizzes must be greater than 0")
	}
	if numberOfQuizzes > 10 {
		return Game{}, errors.New("number of quizzes must be less than 10")
	}
	games := make([]Quiz, numberOfQuizzes)
	for i := 0; i < numberOfQuizzes; i++ {
		quiz := NewQuiz()
		games[i] = quiz
		for k, _ := range games {
			for games[k].Id == quiz.Id && games[k].Id != 0 {
				quiz = NewQuiz()
			}
		}
	}
	logger := slog.Default()
	return Game{games, logger}, nil
}

func (g *Game) addPlayer(name string, quizID int64) error {
	IsFound, _ := g.getPlayer(name)
	if IsFound {
		g.logger.Error("player name already exists")
		return errors.New("player name already exists")
	}
	player := NewPlayer(name)
	for i := range g.Quizzes {
		if g.Quizzes[i].Id == quizID {
			player.QuizId = quizID
			g.Quizzes[i].Players = append(g.Quizzes[i].Players, player)
			return nil
		}
	}
	g.logger.Error("quiz id not found while added a player to the quiz")
	return errors.New("quiz id not found")
}

func (g *Game) submitAnswers(player Player) error {
	err, quiz := g.getQuizByID(player.QuizId)
	if err != nil {
		g.logger.Error(err.Error())
		return err
	}

	quiz.Players = append(quiz.Players, player)
	return nil
}

func (g *Game) getPlayer(name string) (bool, Player) {
	for i := range g.Quizzes {
		for j := range g.Quizzes[i].Players {
			if g.Quizzes[i].Players[j].Name == name && g.Quizzes[i].Players[j].Name != "" {
				return true, g.Quizzes[i].Players[j]
			}
		}
	}
	return false, Player{}
}

func (g *Game) getQuizIDs() []int64 {
	if len(g.Quizzes) == 0 {
	}
	ids := make([]int64, len(g.Quizzes))
	for i := range g.Quizzes {
		ids = append(ids, g.Quizzes[i].Id)
	}
	return ids
}

func (g *Game) getQuestionsByQuizID(quizID int64) ([]Question, error) {
	for i := range g.Quizzes {
		if g.Quizzes[i].Id == quizID {
			return g.Quizzes[i].Questions, nil
		}
	}
	g.logger.Info("failed to find quiz by quiz id %d", "id", quizID)
	return []Question{}, errors.New("quiz not found")
}

func (g *Game) getQuizByID(id int64) (error, *Quiz) {
	for i := range g.Quizzes {
		if g.Quizzes[i].Id == id {
			return nil, &g.Quizzes[i]
		}
	}
	g.logger.Error("failed to find quiz by id %d", "id", fmt.Sprint(id))
	return errors.New("quiz not found"), &Quiz{}
}
