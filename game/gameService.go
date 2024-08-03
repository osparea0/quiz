package game

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type GameService struct {
	game   Game
	logger *slog.Logger
}

func NewGameService() *GameService {
	logger := slog.Default()
	game, err := NewGame(5)
	if err != nil {
		logger.Error("failed to create game", "error", err)
	}
	return &GameService{game: game, logger: logger}
}

func (gs *GameService) RegisterPlayer(w http.ResponseWriter, req *http.Request) {
	newPlayer := struct {
		Name   string `json:"name"`
		QuizID int64  `json:"quiz_id"`
	}{}
	err := json.NewDecoder(req.Body).Decode(&newPlayer)
	if err != nil {
		gs.logger.Error("failed to decode register player request from http request", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = gs.game.addPlayer(newPlayer.Name, newPlayer.QuizID)
	if err != nil {
		gs.logger.Error("failed to add player to game", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (gs *GameService) Play(w http.ResponseWriter, req *http.Request) {
	player := struct {
		Name   string `json:"name"`
		QuizID int64  `json:"quiz_id"`
	}{}
	err := json.NewDecoder(req.Body).Decode(&player)
	if err != nil {
		gs.logger.Error("failed to decode player request from http request", "error", err)
		gs.logger.Info("the payload was", "json", player)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	questions, err := gs.game.getQuestionsByQuizID(player.QuizID)
	if err != nil {
		gs.logger.Error("failed to get questions from game", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json, err := json.Marshal(questions)
	if err != nil {
		gs.logger.Error("failed to marshal questions to json", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		gs.logger.Error("failed to write response", "error", err)
	}
}

func (gs *GameService) Submit(w http.ResponseWriter, req *http.Request) {
	var player Player
	err := json.NewDecoder(req.Body).Decode(&player)
	if err != nil {
		gs.logger.Error("failed to decode submit request from http request", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = gs.game.submitAnswers(player)
	if err != nil {
		gs.logger.Error("failed to submit player to game", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (gs *GameService) GetGrade(w http.ResponseWriter, req *http.Request) {
	var player Player
	err := json.NewDecoder(req.Body).Decode(&player)
	if err != nil {
		gs.logger.Error("failed to decode player in getgrade  from http request", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err, quiz := gs.game.getQuizByID(player.QuizId)
	if err != nil {
		gs.logger.Error("failed to get quiz", "error", err)
	}
	err = quiz.GradeAll()
	if err != nil {
		gs.logger.Error("failed to grade all quiz", "error", err)
		return
	}
	score, err := quiz.getGradeByPlayerName(player.Name)
	if err != nil {
		gs.logger.Error("failed to get grade by player name", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s := struct {
		Score float32 `json:"score"`
	}{}
	s.Score = score
	j, err := json.Marshal(&s)
	if err != nil {
		gs.logger.Error("failed to marshal score to json", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)
	if err != nil {
		gs.logger.Error("failed to write response", "error", err)
		return
	}
}

func (gs *GameService) GetPercentile(w http.ResponseWriter, req *http.Request) {
	var player Player
	err := json.NewDecoder(req.Body).Decode(&player)
	if err != nil {
		gs.logger.Error("failed to decode player in getgrade  from http request", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err, quiz := gs.game.getQuizByID(player.QuizId)
	if err != nil {
		gs.logger.Error("failed to get quiz", "error", err)
		return
	}
	percentile, err := quiz.PercentageOverall(player.Id)
	if err != nil {
		gs.logger.Error("failed to get quiz percentile", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	p := struct {
		Percentile float32 `json:"percentile"`
	}{}
	p.Percentile = percentile
	j, err := json.Marshal(&p)
	if err != nil {
		gs.logger.Error("failed to marshal percentile to json", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(j)
	if err != nil {
		gs.logger.Error("failed to write response", "error", err)
		return
	}

}

func (gs *GameService) GetGameIDs(w http.ResponseWriter, req *http.Request) {
	IDs := gs.game.getQuizIDs()
	w.WriteHeader(http.StatusOK)
	j, err := json.Marshal(IDs)
	if err != nil {
		gs.logger.Error("failed to marshal Ids into json", "error", err)
		return
	}
	w.Write(j)
}

func StartService() {
	gameSvc := NewGameService()
	hostname, err := os.Hostname()
	if err != nil {
		gameSvc.logger.Error("failed to get hostname from os", "error", err)
		return
	}

	http.HandleFunc("/registerplayer/", gameSvc.RegisterPlayer)
	http.HandleFunc("/play/", gameSvc.Play)
	http.HandleFunc("/submitanswers/", gameSvc.Submit)
	http.HandleFunc("/getgameids", gameSvc.GetGameIDs)

	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			gameSvc.logger.Error("http server error %v", err)
		}

		gameSvc.logger.Info("stopped serving new connections", "hostname", hostname)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		gameSvc.logger.Error("http server shutdown error", "error", err.Error())
	}

	gameSvc.logger.Info("graceful shutdown complete")
}
