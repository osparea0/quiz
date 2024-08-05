package game

import (
	"context"
	"encoding/json"
	"errors"
	"io"
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

func NewGameService(numberOfGames int) *GameService {
	logger := slog.Default()
	game, err := NewGame(numberOfGames)
	if err != nil {
		logger.Error("failed to create game", "error", err)
	}
	return &GameService{game: game, logger: logger}
}

func (gs *GameService) RegisterPlayer(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		gs.logger.Info("Invalid request method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		gs.logger.Error("failed to read req body", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newPlayer := Player{}
	err = json.Unmarshal(body, &newPlayer)
	if err != nil {
		gs.logger.Error("failed to decode register player request from http request", "error", err)
		gs.logger.Info("the player payload is", "payload", newPlayer)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = gs.game.addPlayerToQuiz(newPlayer.Name, newPlayer.QuizId)
	if err != nil {
		gs.logger.Error("failed to add player to game", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ok, player := gs.game.getPlayer(newPlayer.Name)
	if !ok {
		gs.logger.Error("failed to get player in register player", "error", err)
		return
	}
	j, err := json.Marshal(player)
	if err != nil {
		gs.logger.Error("failed to marshal player in registerplayer", "error", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(j)
	if err != nil {
		gs.logger.Error("failed to write in registerplayer", "error", err)
		return
	}
}

func (gs *GameService) Play(w http.ResponseWriter, req *http.Request) {
	player := Player{}
	err := json.NewDecoder(req.Body).Decode(&player)
	if err != nil {
		gs.logger.Error("failed to decode player request from http request", "error", err)
		gs.logger.Info("the payload was", "json", player)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	questions, err := gs.game.getQuestionsByQuizID(player.QuizId)
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
	score, err := quiz.getGradeByPlayerID(player.Id)
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

func (gs *GameService) PrintQuiz(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	j, err := json.Marshal(gs.game)
	if err != nil {
		gs.logger.Error("failed to marshal quiz into json", "error", err)
		return
	}
	w.Write(j)
}

func StartService() {
	gameSvc := NewGameService(5)
	hostname, err := os.Hostname()
	if err != nil {
		gameSvc.logger.Error("failed to get hostname from os", "error", err)
		return
	}

	http.HandleFunc("/registerplayer", gameSvc.RegisterPlayer)
	http.HandleFunc("/play", gameSvc.Play)
	http.HandleFunc("/submitanswers", gameSvc.Submit)
	http.HandleFunc("/getgameids", gameSvc.GetGameIDs)
	http.HandleFunc("/getgrade", gameSvc.GetGrade)
	http.HandleFunc("/printquiz", gameSvc.PrintQuiz)
	http.HandleFunc("/getpercentile", gameSvc.GetPercentile)

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
