package game

import (
	"reflect"
	"testing"
)

func TestAnswers_HasOnlyOneTrue(t *testing.T) {
	q := NewQuiz()
	questions, _ := q.Generate()

	qs := make([]Question, 5)
	qs[1].Question = "What color is the sun?"
	qs[1].Answers = createAnswers("Blue", true, "green", false, "yellow", true, "black", false)

	tests := []struct {
		name  string
		input Question
		want  bool
	}{

		{name: "Has only one true answer to question", input: questions[0], want: true},
		{name: "Has two falses and should fail", input: qs[1], want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasOnlyOneTrue(tt.input); got != tt.want {
				t.Errorf("HasOnlyOneTrue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuiz_Grade(t *testing.T) {
	q := NewQuiz()
	questions, _ := q.Generate()
	submittedAnswersWithOneWrong := make([]Question, 5)
	submittedAnswersWithOneWrong[0].Question = "What color is the sun?"
	submittedAnswersWithOneWrong[0].IsRight = true
	submittedAnswersWithOneWrong[0].Answers = createAnswers("Blue", false, "green", false, "yellow", true, "black", false)
	submittedAnswersWithOneWrong[1].Question = "Amazon does not have which of these named services?"
	submittedAnswersWithOneWrong[1].IsRight = true
	submittedAnswersWithOneWrong[1].Answers = createAnswers("Route 53", false, "Elastic Container Registry", false, "Elastic Beanstalk", false, "Elastic Monkey", true)
	submittedAnswersWithOneWrong[2].Question = "Which of these are not google cloud services?"
	submittedAnswersWithOneWrong[2].IsRight = true
	submittedAnswersWithOneWrong[2].Answers = createAnswers("Cloud Run", false, "Cloud SQL", false, "GKE", false, "Cloud Slide", true)
	submittedAnswersWithOneWrong[3].Question = "What color is the sea?"
	submittedAnswersWithOneWrong[3].IsRight = true
	submittedAnswersWithOneWrong[3].Answers = createAnswers("Yellow", false, "Purple", false, "Black", false, "Blue", true)
	submittedAnswersWithOneWrong[4].IsRight = false
	submittedAnswersWithOneWrong[4].Question = "Which subnet mask is the largest (provides the most IP addresses?"
	submittedAnswersWithOneWrong[4].Answers = createAnswers("/32", true, "/29", false, "/27", false, "/16", false)
	players := make([]Player, 2)
	players[0] = Player{
		Name:    "Test Player 1",
		Id:      0,
		QuizId:  q.Id,
		Answers: q.Questions,
		Score:   0,
	}
	players[1] = Player{
		Name:    "Test Player 2",
		Id:      1,
		QuizId:  q.Id,
		Answers: submittedAnswersWithOneWrong,
		Score:   0,
	}
	q.Players = players
	type fields struct {
		Id        int64
		Players   []Player
		Questions []Question
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   error
		want1  float32
	}{
		{name: "one hundred percent incorrect",
			fields: fields{Id: 0, Players: players, Questions: questions},
			args:   args{id: 0},
			want:   nil,
			want1:  0},
		{name: "80 percent correct",
			fields: fields{Id: 1, Players: players, Questions: questions},
			args:   args{id: 1},
			want:   nil,
			want1:  .8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := Quiz{
				Id:        tt.fields.Id,
				Players:   tt.fields.Players,
				Questions: tt.fields.Questions,
			}
			got1, got := q.Grade(tt.args.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Grade() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Grade() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestQuiz_PercentageOverall(t *testing.T) {
	q := NewQuiz()
	submittedAnswersWithOneWrong := make([]Question, 5)
	submittedAnswersWithOneWrong[0].Question = "What color is the sun?"
	submittedAnswersWithOneWrong[0].IsRight = true
	submittedAnswersWithOneWrong[0].Answers = createAnswers("Blue", false, "green", false, "yellow", true, "black", false)
	submittedAnswersWithOneWrong[1].Question = "Amazon does not have which of these named services?"
	submittedAnswersWithOneWrong[1].IsRight = true
	submittedAnswersWithOneWrong[1].Answers = createAnswers("Route 53", false, "Elastic Container Registry", false, "Elastic Beanstalk", false, "Elastic Monkey", true)
	submittedAnswersWithOneWrong[2].Question = "Which of these are not google cloud services?"
	submittedAnswersWithOneWrong[2].IsRight = true
	submittedAnswersWithOneWrong[2].Answers = createAnswers("Cloud Run", false, "Cloud SQL", false, "GKE", false, "Cloud Slide", true)
	submittedAnswersWithOneWrong[3].Question = "What color is the sea?"
	submittedAnswersWithOneWrong[3].IsRight = true
	submittedAnswersWithOneWrong[3].Answers = createAnswers("Yellow", false, "Purple", false, "Black", false, "Blue", true)
	submittedAnswersWithOneWrong[4].IsRight = false
	submittedAnswersWithOneWrong[4].Question = "Which subnet mask is the largest (provides the most IP addresses?"
	submittedAnswersWithOneWrong[4].Answers = createAnswers("/32", true, "/29", false, "/27", false, "/16", false)

	questionsAllCorrect := make([]Question, 5)
	questionsAllCorrect[0].Question = "What color is the sun?"
	questionsAllCorrect[0].IsRight = true
	questionsAllCorrect[0].Answers = createAnswers("Blue", false, "green", false, "yellow", true, "black", false)
	questionsAllCorrect[1].Question = "Amazon does not have which of these named services?"
	questionsAllCorrect[1].IsRight = true
	questionsAllCorrect[1].Answers = createAnswers("Route 53", false, "Elastic Container Registry", false, "Elastic Beanstalk", false, "Elastic Monkey", true)
	questionsAllCorrect[2].Question = "Which of these are not google cloud services?"
	questionsAllCorrect[2].IsRight = true
	questionsAllCorrect[2].Answers = createAnswers("Cloud Run", false, "Cloud SQL", false, "GKE", false, "Cloud Slide", true)
	questionsAllCorrect[3].Question = "What color is the sea?"
	questionsAllCorrect[3].IsRight = true
	questionsAllCorrect[3].Answers = createAnswers("Yellow", false, "Purple", false, "Black", false, "Blue", true)
	questionsAllCorrect[4].IsRight = true
	questionsAllCorrect[4].Question = "Which subnet mask is the largest (provides the most IP addresses?"
	questionsAllCorrect[4].Answers = createAnswers("/32", true, "/29", false, "/27", false, "/16", false)
	players := make([]Player, 2)
	players[0] = Player{
		Name:    "Test Player 1",
		Id:      0,
		QuizId:  q.Id,
		Answers: questionsAllCorrect,
		Score:   0,
	}
	players[1] = Player{
		Name:    "Test Player 2",
		Id:      1,
		QuizId:  q.Id,
		Answers: submittedAnswersWithOneWrong,
		Score:   0,
	}

	q.Players = players

	type fields struct {
		Id        int64
		Players   []Player
		Questions []Question
	}
	type args struct {
		playerId int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   error
		want1  float32
	}{
		{name: "top 50 percentile", fields: fields{
			Id:        1,
			Players:   q.Players,
			Questions: q.Questions,
		}, args: args{playerId: 0},
			want: nil, want1: 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := Quiz{
				Id:        tt.fields.Id,
				Players:   tt.fields.Players,
				Questions: tt.fields.Questions,
			}
			got1, got := q.PercentageOverall(tt.args.playerId)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PercentageOverall() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PercentageOverall() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
