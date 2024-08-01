package main

import "errors"

type Game struct {
	Quizzes []Quiz
}

func NewGame(numberOfQuizzes int) (error, Game) {
	if numberOfQuizzes <= 0 {
		return errors.New("number of quizzes must be greater than 0"), Game{}
	}
	if numberOfQuizzes > 10 {
		return errors.New("number of quizzes must be less than 10"), Game{}
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
	return nil, Game{games}
}

func (g *Game) addPlayer(name string, quizID int64) error {
	IsFound, _ := g.getPlayer(name)
	if IsFound {
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
	return errors.New("quiz id not found")
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
