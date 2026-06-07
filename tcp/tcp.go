package tcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
	"key-value-store/config"
)
var connLimiter = make(chan struct{}, 10)
var wg sync.WaitGroup

func StartTCP(){
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}
	port := cfg.Port
	l, err := net.Listen("tcp",port) 
	if err != nil {
        log.Printf("error creating listener: %v", err)
    }

	sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func(){
		 <-sigChan
        log.Println("shutting down")
        l.Close() 
	}()



    defer l.Close()
	log.Printf("server listening on %s", port)
	
	for{
		conn, err := l.Accept();
		if err!=nil{
			log.Printf("listener closed, waiting for connections to drain")
			break;
		}
		wg.Add(1)
		connLimiter <- struct{}{}
		go handleConnection(conn);
	}
	wg.Wait();
}
func handleConnection(conn net.Conn){
	defer func(){
		conn.Close();
		wg.Done()
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
				if err == bufio.ErrTooLong{
					log.Printf("message exceeds %d bytes, closing connection", maxMessageSize)
				}else{
					log.Printf("Read error: %v", err)
				}	
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
			return;
		}
	}
}
