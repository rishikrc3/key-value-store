package tcp

import (
	"bufio"
	"net"
	"os"
	"testing"
	"time"
)

func TestEdgeCase_MessageTooBig(t *testing.T) {
    os.Setenv("PORT", ":9005")
    os.Setenv("REDIS_URL", "mock-url")
    os.Setenv("APP_NAME", "Test-Big-Msg")
    os.Setenv("DEBUG", "false")

    go StartTCP()

    var conn net.Conn
    var err error
    for i := 0; i < 20; i++ {
        conn, err = net.Dial("tcp", "localhost:9005")
        if err == nil {
            break
        }
        time.Sleep(10 * time.Millisecond)
    }
    if err != nil {
        t.Fatalf("failed to connect: %v", err)
    }
    defer conn.Close()

    giantMessage := make([]byte, 4500)
    for i := range giantMessage {
        giantMessage[i] = 'A'
    }
    giantMessage[len(giantMessage)-1] = '\n' // scanner needs a newline

    _, err = conn.Write(giantMessage)
    if err != nil {
        t.Fatalf("failed to write: %v", err)
    }

    conn.SetReadDeadline(time.Now().Add(2 * time.Second))
    reader := bufio.NewReader(conn)
    _, err = reader.ReadString('\n')
    if err == nil {
        t.Error("expected server to drop connection on oversized message, but got a response")
    }


}

func TestHappyPath_return(t *testing.T){
	os.Setenv("PORT", ":9006")
    os.Setenv("REDIS_URL", "mock-url")
    os.Setenv("APP_NAME", "Test-Big-Msg")
    os.Setenv("DEBUG", "false")

	go StartTCP()

	var conn net.Conn
	var err error
	for i:=0; i<20; i++{
		conn, err = net.Dial("tcp", "localhost:9006")
		if err==nil{
			break;
		}
		time.Sleep(10 *time.Millisecond)
	}

	if err!=nil{
		t.Fatalf("Failed to connect with server")
	}
	defer conn.Close()

	message :=make([]byte,10)

	for i:= range(message){
		message[i] = 'a'
	}
	message[len(message)-1]= '\n'
	_, err = conn.Write(message)
	  if err != nil {
        t.Fatalf("failed to write: %v", err)
    }
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	reader:=bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	if response != "ACK: AAAAAAAAA\n" {
		t.Errorf("expected 'ACK: AAAAAAAAA\n', got %q", response)
	}
}
func TestEdgecase_returnaubrupt(t *testing.T){
	os.Setenv("PORT", ":9007") 
    os.Setenv("REDIS_URL", "mock-url")
    os.Setenv("APP_NAME", "Test-Big-Msg")
    os.Setenv("DEBUG", "false")

	go StartTCP()

	var conn net.Conn
	var err error
	for i:=0; i<20; i++{
		conn, err = net.Dial("tcp", "localhost:9007")
		if err==nil{
			break;
		}
		time.Sleep(10 *time.Millisecond)
	}

	if err!=nil{
		t.Fatalf("Failed to connect with server")
	}
	defer conn.Close()

	message :=make([]byte,10)

	for i:= range(message){
		message[i] = 'a'
	}
	message[len(message)-1]= '\n'
	_, err = conn.Write(message)
	  if err != nil {
        t.Fatalf("failed to write: %v", err)
		return
    }
	conn.Close()
}

