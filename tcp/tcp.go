package tcp

import (
	"bufio"
	"log"
	"net"
	"strings"
	"fmt"
)

func StartTCP(){
	address := "localhost:6000"

	l, err := net.Listen("tcp",address) 
	if err != nil {
        log.Fatalf("error creating listener: %v", err)
    }

    defer l.Close()
	log.Printf("server listening on %s", address)
	
	for{
		conn, err := l.Accept();
		if err!=nil{
			log.Printf("Error accepting Connections");
			continue;
		}
		go handleConnection(conn);
	}
}
func handleConnection(conn net.Conn){
	defer conn.Close();
	reader := bufio.NewReader(conn)
    message, err := reader.ReadString('\n')
    if err != nil {
        log.Printf("Read error: %v", err)
        return
    }
	ackMsg := strings.ToUpper(strings.TrimSpace(message))
    response := fmt.Sprintf("ACK: %s\n", ackMsg)
    _, err = conn.Write([]byte(response))
    if err != nil {
        log.Printf("Server write error: %v", err)
    }
}
