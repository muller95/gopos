package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_"github.com/go-sql-driver/mysql"
	"log"
	"net"
	"os"
)

var goposServerPassword, goposServerPort, goposSQLUser, goposSQLPassword string
var dbConn *sql.DB

func handleConnection(conn net.Conn) {
	requestMap := make(map[string]string)
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&requestMap)
	if err != nil {
		log.Fatal("Error on decoding json: ", err)
	}

	fmt.Println(requestMap)
}

func main() {
	var err error

	goposServerPassword = os.Getenv("GOPOS_SERVER_PASSWORD")
	if goposServerPassword == "" {
		log.Fatal("GOPOS_SERVER_PASSWORD is not set")
	}
	fmt.Println(goposServerPassword)

	goposServerPort = os.Getenv("GOPOS_SERVER_PORT")
	if goposServerPort == "" {
		log.Fatal("GOPOS_SERVER_PORT is not set")
	}
	fmt.Println(goposServerPort)

	goposSQLUser = os.Getenv("GOPOS_SQL_USER")
	if goposSQLUser == "" {
		log.Fatal("GOPOS_SQL_USER is not set")
	}
	fmt.Println(goposSQLUser)

	goposSQLPassword = os.Getenv("GOPOS_SQL_PASSWORD")
	if goposSQLPassword == "" {
		log.Fatal("GOPOS_SQL_PASSWORD is not set")
	}
	fmt.Println(goposSQLPassword)

	dbConn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/gopos", goposSQLUser, goposSQLPassword))

	if err != nil {
		log.Fatal("Error on opening database", err)
	}

	listener, err :=  net.Listen("tcp", ":" + goposServerPort)
	if err != nil {
		log.Fatal("Cannot start listening port: ", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Cannot accept connection: ", err)
		}

		go handleConnection(conn)
	}
}
