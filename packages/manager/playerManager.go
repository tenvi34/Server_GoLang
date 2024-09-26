package manager

import (
	"encoding/binary"
	"errors"
	"log"

	pb "Server_TCP/messages"

	"net"

	"google.golang.org/protobuf/proto"
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

	myPlayerSapwn := &pb.GameMessage{
		Message: &pb.GameMessage_SpawnMyPlayer{
			SpawnMyPlayer: &pb.SpawnMyPlayer{
				X: 0,
				Y: 0,
				Z: 0,
			},
		},
	}

	response := GetNetManager().MakePacket(myPlayerSapwn)
	(*player.Conn).Write(response)

	otherPlayerSpawnPacket := &pb.GameMessage{
		Message: &pb.GameMessage_SpawnOtherPlayer{
			SpawnOtherPlayer: &pb.SpawnOtherPlayer{
				PlayerId: name,
				X:        0,
				Y:        0,
				Z:        0,
			},
		},
	}

	response = GetNetManager().MakePacket(otherPlayerSpawnPacket)

	for _, p := range pm.players {
		if p.Name == name {
			continue
		}

		(*p.Conn).Write(response)
	}

	for _, p := range pm.players {
		if p.Name == name {
			continue
		}

		otherPlayerSpawnPacket := &pb.GameMessage{
			Message: &pb.GameMessage_SpawnOtherPlayer{
				SpawnOtherPlayer: &pb.SpawnOtherPlayer{
					PlayerId: p.Name,
					X:        0,
					Y:        0,
					Z:        0,
				},
			},
		}

		response = GetNetManager().MakePacket(otherPlayerSpawnPacket)

		(*player.Conn).Write(response)
	}

	return player
}

func (pm *PlayerManager) MovePlayer(name string, x float32, y float32, z float32) {
	gameMessage := &pb.GameMessage{
		Message: &pb.GameMessage_PlayerPosition{
			PlayerPosition: &pb.PlayerPosition{
				PlayerId: name,
				X:        x,
				Y:        y,
				Z:        z,
			},
		},
	}

	response, err := proto.Marshal(gameMessage)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}

	for _, player := range pm.players {
		if player.Name == name {
			continue
		}

		lengthBuf := make([]byte, 4)
		binary.LittleEndian.PutUint32(lengthBuf, uint32(len(response)))
		lengthBuf = append(lengthBuf, response...)
		(*player.Conn).Write(lengthBuf)
	}
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
