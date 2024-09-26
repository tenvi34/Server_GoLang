package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	pb "Server_TCP/messages"

	"google.golang.org/protobuf/proto"

	mg "Server_TCP/packages/manager"
)

func main() {

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()
	fmt.Println("Server is listening on :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		// 메시지 길이를 먼저 읽습니다 (4바이트)
		lengthBuf := make([]byte, 4)
		_, err := conn.Read(lengthBuf)
		if err != nil {
			log.Printf("Failed to read message length: %v", err)
			return
		}
		length := binary.LittleEndian.Uint32(lengthBuf)

		// 메시지 본문을 읽습니다
		messageBuf := make([]byte, length)
		_, err = conn.Read(messageBuf)
		if err != nil {
			log.Printf("Failed to read message body: %v", err)
			return
		}

		// Protocol Buffers 메시지를 파싱합니다
		message := &pb.GameMessage{}
		err = proto.Unmarshal(messageBuf, message)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// 메시지 처리
		processMessage(message, &conn)
	}
}

func processMessage(message *pb.GameMessage, conn *net.Conn) {
	switch msg := message.Message.(type) {
	case *pb.GameMessage_PlayerPosition:
		pos := msg.PlayerPosition
		fmt.Println("Position : ", pos.X, pos.Y, pos.Z)
		mg.GetPlayerManager().MovePlayer(pos.PlayerId, pos.X, pos.Y, pos.Z)
	case *pb.GameMessage_Chat:
		chat := msg.Chat
		mg.GetChatManager().Broadcast(chat.Sender, chat.Content)
	case *pb.GameMessage_Login:
		playerId := msg.Login.PlayerId
		fmt.Println(playerId)
		playerManager := mg.GetPlayerManager()
		playerManager.AddPlayer(playerId, 0, conn)
	default:
		panic(fmt.Sprintf("unexpected messages.isGameMessage_Message: %#v", msg))
	}
}
