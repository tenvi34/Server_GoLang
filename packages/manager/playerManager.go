package manager

import (
	"errors"
	"net"
)

var playerManager *PlayerManager

// Player represents a single player with some attributes
type Player struct {
	ID   int
	Name string
	Age  int
	Conn *net.Conn
}

// PlayerManager manages a list of players
type PlayerManager struct {
	players map[int]Player
	nextID  int
}

// NewPlayerManager creates a new PlayerManager
func GetPlayerManager() *PlayerManager {
	if playerManager == nil {
		playerManager = &PlayerManager{
			players: make(map[int]Player),
			nextID:  1,
		}
	}

	return playerManager
}

// AddPlayer adds a new player to the manager
func (pm *PlayerManager) AddPlayer(name string, age int, conn *net.Conn) Player {
	player := Player{
		ID:   pm.nextID,
		Name: name,
		Age:  age,
		Conn: conn,
	}

	pm.players[pm.nextID] = player
	pm.nextID++
	return player
}

// GetPlayer retrieves a player by ID
func (pm *PlayerManager) GetPlayer(id int) (Player, error) {
	player, exists := pm.players[id]
	if !exists {
		return Player{}, errors.New("player not found")
	}
	return player, nil
}

// RemovePlayer removes a player by ID
func (pm *PlayerManager) RemovePlayer(id int) error {
	if _, exists := pm.players[id]; !exists {
		return errors.New("player not found")
	}
	delete(pm.players, id)
	return nil
}

// ListPlayers returns all players in the manager
func (pm *PlayerManager) ListPlayers() []Player {
	playerList := []Player{}
	for _, player := range pm.players {
		playerList = append(playerList, player)
	}
	return playerList
}
