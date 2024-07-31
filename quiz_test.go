package main

import (
	"reflect"
	"testing"
)

func TestAnswers_HasOnlyOneTrue(t *testing.T) {
	q := NewQuiz()
	_, questions := q.Generate()

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
	_, questions := q.Generate()
	players := make([]Player, 2)
	players[0] = Player{
		Name:    "Tes Player",
		Id:      1,
		QuizId:  q.Id,
		Answers: q.Questions,
		Score:   0,
	}
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
		{name: "one hundred percent correct",
			fields: fields{Id: 1, Players: players, Questions: questions},
			args:   args{id: 1},
			want:   nil,
			want1:  1.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := Quiz{
				Id:        tt.fields.Id,
				Players:   tt.fields.Players,
				Questions: tt.fields.Questions,
			}
			got, got1 := q.Grade(tt.args.id)
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := Quiz{
				Id:        tt.fields.Id,
				Players:   tt.fields.Players,
				Questions: tt.fields.Questions,
			}
			got, got1 := q.PercentageOverall(tt.args.playerId)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PercentageOverall() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PercentageOverall() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_computeGrade(t *testing.T) {
	type args struct {
		submittedAnswers []Question
		correctAnswers   []Question
	}
	tests := []struct {
		name  string
		args  args
		want  error
		want1 float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := computeGrade(tt.args.submittedAnswers, tt.args.correctAnswers)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("computeGrade() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("computeGrade() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_createAnswers(t *testing.T) {
	type args struct {
		ans1     string
		ans1bool bool
		ans2     string
		ans2bool bool
		ans3     string
		ans3bool bool
		ans4     string
		ans4bool bool
	}
	tests := []struct {
		name string
		args args
		want Answers
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createAnswers(tt.args.ans1, tt.args.ans1bool, tt.args.ans2, tt.args.ans2bool, tt.args.ans3, tt.args.ans3bool, tt.args.ans4, tt.args.ans4bool); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createAnswers() = %v, want %v", got, tt.want)
			}
		})
	}
}
