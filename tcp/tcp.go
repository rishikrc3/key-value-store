package tcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
	"github.com/joho/godotenv"
)
var connLimiter = make(chan struct{}, 100)

func StartTCP(){
	err:= godotenv.Load()
	if err!=nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	l, err := net.Listen("tcp",port) 
	if err != nil {
        log.Fatalf("error creating listener: %v", err)
    }

    defer l.Close()
	log.Printf("server listening on %s", port)
	
	for{
		conn, err := l.Accept();
		if err!=nil{
			log.Printf("accept error: %v", err)
			continue;
		}
		connLimiter <- struct{}{}
		go handleConnection(conn);
	}
}
func handleConnection(conn net.Conn){
	defer func(){
		conn.Close();
		<- connLimiter
	} ()
	const maxMessageSize = 4 * 1024
	scanner := bufio.NewScanner(conn)
	buff := make([]byte, maxMessageSize)
	scanner.Buffer(buff, maxMessageSize)
	
    
	for{
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				log.Printf("Read error: %v", err)
			} else {
				log.Printf("Client %s disconnected", conn.RemoteAddr())
			}
			return
		}
		message := scanner.Text()
		ackMsg := strings.ToUpper(strings.TrimSpace(message))
		response := fmt.Sprintf("ACK: %s\n", ackMsg)

		_, err := conn.Write([]byte(response))
		if err != nil {
			log.Printf("Server write error: %v", err)
		}
	}
}
