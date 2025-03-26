package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/gara22/tetris/game"
)

type JsonRepository struct {
	gameMutex sync.Mutex
}

var _ Repository = (*JsonRepository)(nil)

const (
	GAME_FILE           = "./games.json"
	HIGH_SCORE_FILE     = "./high_scores.json"
	MAXIMUM_HIGH_SCORES = 100
)

func New() JsonRepository {
	return JsonRepository{}
}

func (r *JsonRepository) SaveGame(game game.PersistedGame) error {
	r.gameMutex.Lock()
	fmt.Println("Locked for saving")
	defer func() {
		fmt.Println("Unlocked after saving")
		r.gameMutex.Unlock()
	}()

	existing, err := r.loadGames()
	if err != nil {
		return fmt.Errorf("failed loading games")
	}

	var found bool

	for i, g := range existing {
		if g.ID == game.ID {
			existing[i] = game
			found = true
			break
		}
	}

	if !found {
		existing = append(existing, game)
	}

	bytes, err := json.Marshal(existing)
	if err != nil {
		return fmt.Errorf("failed to save game %w", err)
	}

	err = os.WriteFile(GAME_FILE, bytes, 0644)
	if err != nil {
		return fmt.Errorf("failed write to file %w", err)
	}

	spew.Dump(bytes)
	return nil
}

func (r *JsonRepository) GetGames() ([]game.PersistedGame, error) {
	r.gameMutex.Lock()
	defer r.gameMutex.Unlock()
	return r.loadGames()
}

func (r *JsonRepository) GetByID(id string) (game.PersistedGame, error) {
	r.gameMutex.Lock()
	defer r.gameMutex.Unlock()

	games, err := r.loadGames()
	if err != nil {
		return game.PersistedGame{}, fmt.Errorf("failed loading games")
	}

	for _, g := range games {
		if g.ID == id {
			return g, nil
		}
	}

	return game.PersistedGame{}, fmt.Errorf("game not found")
}

func (r *JsonRepository) GetHighScores() ([]game.HighScore, error) {
	r.gameMutex.Lock()
	defer r.gameMutex.Unlock()

	return r.loadHighScore()
}

// SaveHighScore implements Repository.
func (r *JsonRepository) SaveHighScore(highScore game.HighScore) error {
	r.gameMutex.Lock()
	defer r.gameMutex.Unlock()

	highScores, err := r.loadHighScore()
	if err != nil {
		return fmt.Errorf("failed loading high scores")
	}

	// find the index to insert the new high score
	var index int
	for i, h := range highScores {
		if h.Score <= highScore.Score {
			index = i
			break
		}
	}

	spew.Dump(index)

	// insert the new high score
	highScores = slices.Insert(highScores, index, highScore)

	// remove the last high score if the list is too long
	if len(highScores) > MAXIMUM_HIGH_SCORES {
		highScores = highScores[:MAXIMUM_HIGH_SCORES]
	}

	bytes, err := json.Marshal(highScores)
	if err != nil {
		return fmt.Errorf("failed to save high score %w", err)
	}

	err = os.WriteFile(HIGH_SCORE_FILE, bytes, 0644)
	if err != nil {
		return fmt.Errorf("failed write to file %w", err)
	}

	return nil

}

func (r *JsonRepository) loadHighScore() ([]game.HighScore, error) {
	file, err := os.ReadFile(HIGH_SCORE_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			return []game.HighScore{}, nil
		}
		return nil, err
	}

	var highScores []game.HighScore
	if err := json.Unmarshal(file, &highScores); err != nil {
		return nil, err
	}
	return highScores, nil

}

// LoadGames loads all games from file
func (r *JsonRepository) loadGames() ([]game.PersistedGame, error) {

	file, err := os.ReadFile(GAME_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			return []game.PersistedGame{}, nil
		}
		return nil, err
	}

	var games []game.PersistedGame
	if err := json.Unmarshal(file, &games); err != nil {
		return nil, err
	}
	return games, nil
}
