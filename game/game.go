package game

import (
	"errors"
	"fmt"
	"log/slog"
)

type Game struct {
	Quizzes []Quiz
	logger  *slog.Logger
}

type Gamer interface {
	addPlayerToQuiz(name string, quizID int64) error
	submitAnswers(player Player) error
	getPlayer(name string) (bool, Player)
	getQuizIDs() []int64
	getQuestionsByQuizID(quizID int64) ([]Question, error)
	getQuizByID(id int64) (*Quiz, error)
}

func NewGame(numberOfQuizzes int) (*Game, error) {
	if numberOfQuizzes <= 0 {
		return &Game{}, errors.New("number of quizzes must be greater than 0")
	}
	if numberOfQuizzes > 10 {
		return &Game{}, errors.New("number of quizzes must be less than 10")
	}
	games := make([]Quiz, numberOfQuizzes)
	for i := 0; i < numberOfQuizzes; i++ {
		quiz := NewQuiz()
		games[i] = quiz
	}
	logger := slog.Default()
	return &Game{games, logger}, nil
}

func (g *Game) addPlayerToQuiz(name string, quizID int64) error {
	IsFound, _ := g.getPlayer(name)
	if IsFound {
		g.logger.Error("player name already exists")
		return errors.New("player name already exists")
	}
	player := NewPlayer(name, quizID)
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
	q, err := g.getQuizByID(player.QuizId)
	if err != nil {
		g.logger.Error(err.Error())
		return err
	}

	for i := range q.Players {
		if q.Players[i].Id == player.Id {
			q.Players[i] = player
		}
	}
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
		return []int64{}
	}
	ids := make([]int64, len(g.Quizzes))
	for i := range g.Quizzes {
		ids[i] = g.Quizzes[i].Id
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

func (g *Game) getQuizByID(id int64) (*Quiz, error) {
	for i := range g.Quizzes {
		if g.Quizzes[i].Id == id {
			return &g.Quizzes[i], nil
		}
	}
	g.logger.Error("failed to find quiz by id %d", "id", fmt.Sprint(id))
	return &Quiz{}, errors.New("quiz not found")
}
