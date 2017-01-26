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

func handleRequestGroup(requestMap map[string]string, conn net.Conn) {
	switch requestMap["group"]{
		case "WORKER":
			switch requestMap["action"] {
				case "ADD":
					handleWorkerAdd(requestMap, conn)
				case "GET":
					handleWorkerGet(conn)
				case "DELETE":
					handleWorkerDelete(requestMap, conn)
			}
		case "TABLE":
			switch requestMap["action"] {
				case "ADD":
					handleTableAdd(requestMap, conn)
				case "GET":
					handleTableGet(conn)
				case "DELETE":
					handleTableDelete(requestMap, conn)
			}
		case "CATEGORY":
			switch requestMap["action"] {
				case "ADD":
					handleCategoryAdd(requestMap, conn)
				case "GET":
					handleCategoryGet(conn)
				case "DELETE":
					handleCategoryDelete(requestMap, conn)
			}
		case "DISH":
			switch requestMap["action"] {
				case "ADD":
					handleDishAdd(requestMap, conn)
				case "GET":
					handleDishGet(requestMap, conn)
				case "DELETE":
					handleDishDelete(requestMap, conn)
			}
		case "CARD":
			switch requestMap["action"] {
				case "ADD":
					handleCardAdd(requestMap, conn)
				case "GET":
					handleCardGet(conn)
				case "DELETE":
					handleCardDelete(requestMap, conn)
			}
	}
}

func handleConnection(conn net.Conn) {
	requestMap := make(map[string]string)
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&requestMap)
	if err != nil {
		log.Fatal("Error on decoding json: ", err)
	}


	if requestMap["password"] != goposServerPassword {
		responseMap := make(map[string]string)
		responseMap["result"] = "ERR"
		responseMap["error"] = "Неправильный пароль"
		encoder := json.NewEncoder(conn)
		err := encoder.Encode(responseMap)
		if err != nil {
			log.Fatal("Error on encoding response json: ", err)
		}

		return
	}

	handleRequestGroup(requestMap, conn)
	conn.Close()
}

func main() {
	var err error

	goposServerPassword = os.Getenv("GOPOS_SERVER_PASSWORD")
	if goposServerPassword == "" {
		log.Fatal("GOPOS_SERVER_PASSWORD is not set")
	}

	goposServerPort = os.Getenv("GOPOS_SERVER_PORT")
	if goposServerPort == "" {
		log.Fatal("GOPOS_SERVER_PORT is not set")
	}

	goposSQLUser = os.Getenv("GOPOS_SQL_USER")
	if goposSQLUser == "" {
		log.Fatal("GOPOS_SQL_USER is not set")
	}

	goposSQLPassword = os.Getenv("GOPOS_SQL_PASSWORD")
	if goposSQLPassword == "" {
		log.Fatal("GOPOS_SQL_PASSWORD is not set")
	}

	dbConn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/gopos?parseTime=true", goposSQLUser,
		goposSQLPassword))
	_, err = dbConn.Exec("SET CHARSET utf8")
	if err != nil {
		log.Fatalf("Error on setting charset: ", err)
	}

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
