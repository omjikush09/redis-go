package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"github.com/omjikush09/redis-go/app"
	"github.com/omjikush09/redis-go/app/data_structure"
)

var _ = net.Listen
var _ = os.Exit

type ConnectionConfig struct{
	conn net.Conn
	inTransaction bool
	commands [][]string
}


func handleConnection(config *ConnectionConfig) {
	conn := config.conn
	defer conn.Close()

	scanner := bufio.NewReader(conn)

	for {
		command, err := app.ParseResp(scanner)
		if err != nil {
			fmt.Printf("Failed to parse command %v", err)
			return
		}
		response:=app.HandleCommands(command, config)
		conn.Write([]byte(response))
	}

}

func main() {
	
	
	data_structure.InitializeStore()

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		var conn net.Conn
		conn, err = l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		config:= ConnectionConfig{
			conn: conn,
			inTransaction: false,
			commands: make([][]string,0),
		}
		go handleConnection(&config)
	}

}
