package main 

import (
    "fmt"
    "key-value-store/tcp" 
)

func main(){
    fmt.Println("Hello World")
    tcp.StartTCP()
}