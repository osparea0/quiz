package main

type Quizzer interface {
	Generate() (error, []Question)
	Grade(submittedQuestions []Question) (error, float32)
	PercentageOverall() (error, float32)
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
	Question string    `json:"question"`
	Answers  []Answers `json:"answers"`
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
