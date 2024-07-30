package main

type Quizzer interface {
	Generate() (error, questions []Question)
	Grade(submittedQuestions []Question) (error, grade float64)
}

type Quiz struct {
	Id        int64      `json:"id"`
	Questions []Question `json:"questions"`
}

type Player struct {
	Name    string `json:"name"`
	Id      int64  `json:"id"`
	Quizzes []Quiz `json:"quizzes"`
}

type Question struct {
	Question string    `json:"question"`
	Answers  []Answers `json:"answers"`
}

type Answers struct {
	Answer1 Answer `json:"answer1,omitempty"`
	Answer2 Answer `json:"answer2,omitempty"`
	Answer3 Answer `json:"answer3,omitempty"`
	Answer4 Answer `json:"answer4,omitempty"`
}

type Answer struct {
	Answer string `json:"answer"`
	IsTrue bool   `json:"is_true"`
}

func (a *Answers) HasOnlyOneTrue() bool {
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
