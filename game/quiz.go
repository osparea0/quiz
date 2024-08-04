package game

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"sort"
)

type Quizzer interface {
	Generate() ([]Question, error)
	Grade(ID int64) (float32, error)
	PercentageOverall(playerId int64) (float32, error)
}

// Quiz struct holds the fields for a single quiz and implements the Quizzer interface
type Quiz struct {
	Id        int64      `json:"id"`
	Players   []Player   `json:"players"`
	Questions []Question `json:"questions"`
}

func NewQuiz() Quiz {
	q := Quiz{}
	questions, err := q.Generate()
	if err != nil {
		slog.Error(err.Error())
		return q
	}
	id := rand.Int64()
	for id == 0 {
		id = rand.Int64()
	}
	return Quiz{Id: id, Questions: questions}
}

// Grade takes the player id and computes their grade
func (q *Quiz) Grade(ID int64) (float32, error) {
	for i := 0; i < len(q.Players); i++ {
		if q.Players[i].Id == ID {
			err, score := computeGrade(q.Players[i].Answers, q.Questions)
			if err != nil {
				slog.Error("failed to compute grade", "error", err)
				return 0, err
			}

			q.Players[i].Score = score
			return score, nil
		}
	}
	return 0, errors.New("failed to find player's game record in history")
}

func (q *Quiz) GradeAll() error {
	for i := range q.Players {
		grade, err := q.Grade(q.Players[i].Id)
		if err != nil {
			slog.Error("failed to grade all in quiz", "error", err)
			return err
		}
		slog.Info("gradeAll", "grade", grade)
	}
	return nil
}
func (q *Quiz) PercentageOverall(playerId int64) (float32, error) {
	err := q.GradeAll()
	if err != nil {
		return 0, err
	}
	var idIndex int = -1
	sort.Slice(q.Players, func(i, j int) bool {
		return q.Players[i].Score < q.Players[j].Score
	})

	for i := range q.Players {
		if q.Players[i].Id == playerId {
			idIndex = i
			break
		}
	}

	if len(q.Players) == 1 {
		return 1.0, nil
	}

	if idIndex == -1 {
		return 0, fmt.Errorf("player ID %d not found", playerId)
	}

	percentile := float32(idIndex) / float32(len(q.Players)-1) * 100
	slog.Info("logging percentile", "percentile", percentile)
	return percentile, nil
}

func (q *Quiz) Generate() ([]Question, error) {
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

	return questions, nil
}

// Player is the struct to hold the history of a single players past quizzes
type Player struct {
	Name    string     `json:"name"`
	Id      int64      `json:"id"`
	QuizId  int64      `json:"quiz_id"`
	Answers []Question `json:"answers"`
	Score   float32    `json:"score"`
}

func NewPlayer(name string, quizID int64) Player {
	id := rand.Int64()
	for id == 0 {
		id = rand.Int64()
	}
	p := Player{
		Name:    name,
		Id:      id,
		QuizId:  quizID,
		Answers: nil,
		Score:   0,
	}
	return p
}

// Question struct holds a single question and all of it's answers
type Question struct {
	Question string  `json:"question"`
	Answers  Answers `json:"answers"`
	IsRight  bool    `json:"is_right"`
}

func (q *Question) GetCorrectAnswer() string {
	if q.Answers.Answer1.IsTrue {
		return q.Answers.Answer1.Answer
	}
	if q.Answers.Answer2.IsTrue {
		return q.Answers.Answer2.Answer
	}
	if q.Answers.Answer3.IsTrue {
		return q.Answers.Answer3.Answer
	}
	if q.Answers.Answer4.IsTrue {
		return q.Answers.Answer4.Answer
	}
	return ""
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
	for i := 0; i < len(submittedAnswers); i++ {
		if submittedAnswers[i].IsRight {
			count++
		}
	}

	if count == 0 {
		return nil, 0
	}
	score := float32(count) / float32(len(submittedAnswers))

	return nil, score
}

// HasOnlyOneTrue ensures there is only one true answer
func HasOnlyOneTrue(q Question) bool {
	counter := 0

	if q.Answers.Answer1.IsTrue {
		counter++
	}
	if q.Answers.Answer2.IsTrue {
		counter++
	}
	if q.Answers.Answer3.IsTrue {
		counter++
	}
	if q.Answers.Answer4.IsTrue {
		counter++
	}
	return counter == 1
}

func (q *Quiz) getGradeByPlayerID(ID int64) (float32, error) {
	for i := range q.Players {
		if q.Players[i].Id == ID {
			return q.Players[i].Score, nil
		}
	}
	return 0, errors.New("failed to get grade by player name")
}
