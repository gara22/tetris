package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/gara22/tetris/game"
)

type JsonRepository struct {
	file string
	mu   sync.Mutex
}

var _ Repository = (*JsonRepository)(nil)

func New(path string) JsonRepository {
	return JsonRepository{
		file: path,
	}
}

func (r *JsonRepository) SaveGame(game game.PersistedGame) error {
	r.mu.Lock()
	fmt.Println("Locked for saving")
	defer func() {
		fmt.Println("Unlocked after saving")
		r.mu.Unlock()
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

	err = os.WriteFile(r.file, bytes, 0644)
	if err != nil {
		return fmt.Errorf("failed write to file %w", err)
	}

	spew.Dump(bytes)
	return nil
}

func (r *JsonRepository) GetGames() ([]game.PersistedGame, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.loadGames()
}

func (r *JsonRepository) GetByID(id string) (game.PersistedGame, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

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

// LoadGames loads all games from file
func (r *JsonRepository) loadGames() ([]game.PersistedGame, error) {

	file, err := os.ReadFile(r.file)
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
