package main

import (
	"errors"
	"log/slog"
	"sort"
)

type Quizzer interface {
	Generate() (error, []Question)
	Grade(player Player, submittedQuestions []Question) (error, float32)
	PercentageOverall(playerId int64) (error, float32)
}

type Games struct {
	Quizzes []Quiz
}

// Quiz struct holds the fields for a single quiz and implements the Quizzer interface
type Quiz struct {
	Id        int64      `json:"id"`
	Players   []Player   `json:"players"`
	Questions []Question `json:"questions"`
}

// Grade takes the player id and computes their grade
func (q Quiz) Grade(id int64) (error, float32) {
	for i := 0; i < len(q.Players); i++ {
		if q.Players[i].Id == id {
			err, score := computeGrade(q.Players[i].Answers, q.Questions)
			if err != nil {
				slog.Error("failed to compute grade", "error", err)
				return err, 0
			}
			return nil, score
		}
	}
	return nil, 0
}

func (q Quiz) PercentageOverall(playerId int64) (error, float32) {
	var idIndex int
	var total float32 = 0
	sort.Slice(q.Players, func(i, j int) bool {
		return q.Players[i].Score < q.Players[j].Score
	})
	for i := range q.Players {
		if q.Players[i].Id == playerId {
			idIndex = i
		}
	}
	if idIndex == 0 {
		return nil, 0
	}
	percentile := len(q.Players) / idIndex
	return nil, float32(percentile)
}

func (q Quiz) Generate() (error, []Question) {
	questions := make([]Question, 5)
	questions[0].Question = "What color is the sun?"
	questions[0].Answers = createAnswers("Blue", false, "green", false, "yellow", true, "black", false)
	questions[1].Question = "Amazon does not have which of these named services?"
	questions[1].Answers = createAnswers("Route 53", false, "Elastic Container Registry", false, "Elastic Beanstalk", false, "Elastic Monkey", true)
	questions[2].Question = "Which of these are not google cloud services?"
	questions[2].Answers = createAnswers("Cloud Run", false, "Cloud SQL", false, "GKE", false, "Cloud Slide", true)
	questions[3].Question = "What color is the sea?"
	questions[3].Answers = createAnswers("Yellow", false, "Purple", false, "Black", false, "Blue", true)
	questions[4].Question = "Which subnet mask is the largest (provides the most IP addresses?"
	questions[4].Answers = createAnswers("/32", false, "/29", false, "/27", false, "/16", true)

	return nil, questions
}

// Player is the struct to hold the history of a single players past quizzes
type Player struct {
	Name    string     `json:"name"`
	Id      int64      `json:"id"`
	QuizId  int64      `json:"quiz_id"`
	Answers []Question `json:"answers"`
	Score   float32    `json:"score"`
}

// Question struct holds a single question and all of it's answers
type Question struct {
	Question string  `json:"question"`
	Answers  Answers `json:"answers"`
}

// Answer holds all possible answers for a single question and is comprised of Answer structs
type Answers struct {
	Answer1 Answer `json:"answer1,omitempty"`
	Answer2 Answer `json:"answer2,omitempty"`
	Answer3 Answer `json:"answer3,omitempty"`
	Answer4 Answer `json:"answer4,omitempty"`
}

// Answer is a single answer and whether it is true or not
type Answer struct {
	Answer string `json:"answer"`
	IsTrue bool   `json:"is_true"`
}

func createAnswers(ans1 string, ans1bool bool, ans2 string, ans2bool bool, ans3 string, ans3bool bool, ans4 string, ans4bool bool) Answers {
	Answer1 := Answer{Answer: ans1, IsTrue: ans1bool}
	Answer2 := Answer{Answer: ans2, IsTrue: ans2bool}
	Answer3 := Answer{Answer: ans3, IsTrue: ans3bool}
	Answer4 := Answer{Answer: ans4, IsTrue: ans4bool}
	return Answers{Answer1: Answer1, Answer2: Answer2, Answer3: Answer3, Answer4: Answer4}
}

func computeGrade(submittedAnswers []Question, correctAnswers []Question) (error, float32) {
	count := 0
	if len(submittedAnswers) != len(correctAnswers) {
		return errors.New("number of submitted answers do not match number of questions"), 0
	}
	for i := 0; i < len(submittedAnswers[i].Question); i++ {
		if submittedAnswers[i].Answers == correctAnswers[i].Answers {
			count++
		}
	}
	score := float32(len(submittedAnswers) / count)

	return nil, score
}

// HasOnlyOneTrue ensures there is only one true answer
func (a Answers) HasOnlyOneTrue() bool {
	counter := 0

	if a.Answer1.IsTrue {
		counter++
	}
	if a.Answer2.IsTrue {
		counter++
	}
	if a.Answer3.IsTrue {
		counter++
	}
	if a.Answer4.IsTrue {
		counter++
	}
	return counter == 1
}
