package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "Server_TCP/messages"

	"google.golang.org/grpc"
)

type gameServer struct {
	pb.UnimplementedGameServiceServer
	mu      sync.Mutex
	players map[int32]*pb.PlayerPosition
}

func (s *gameServer) UpdatePosition(ctx context.Context, pos *pb.PlayerPosition) (*pb.GameState, error) {

	// 스레드를 안전하게 락을 건다.
	s.mu.Lock()

	// 함수가 끝날 때 락을 푼다.
	defer s.mu.Unlock()

	// Server의 players의 pos를 업데이트
	s.players[pos.PlayerId] = pos

	// 저장할 State를 생성
	gameState := &pb.GameState{
		Players: make([]*pb.PlayerPosition, 0, len(s.players)),
	}

	// 플레이어들의 목록을 append
	for _, player := range s.players {
		gameState.Players = append(gameState.Players, player)
	}

	log.Printf("Updated position for player %d: (%f, %f, %f)", pos.PlayerId, pos.X, pos.Y, pos.Z)

	// gameState를 반환하고 error는 nil이다.
	return gameState, nil
}

func main() {
	// net 패키지 안에 listen을 호출하면 자동으로 tcp 서버가 생성
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// grpc에 있는 API인(자동 생성) NewServer을 만든다.
	s := grpc.NewServer()

	// 서버로 등록
	pb.RegisterGameServiceServer(s, &gameServer{
		players: make(map[int32]*pb.PlayerPosition),
	})

	fmt.Println("Game server is running on :50051")

	// 네트워크를 유지
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
