package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var goposServerPassword, goposServerPort, goposSQLUser, goposSQLPassword string
var goposCheckPath, goposOrderPath, goposPrintserviceFont string
var goposCheckWidth, goposCheckHeight, goposOrderWidth, goposOrderHeight string
var dbConn *sql.DB

func handleRequestGroup(requestMap map[string]string, conn net.Conn) {
	switch requestMap["group"] {
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
	case "ORDER":
		switch requestMap["action"] {
		case "CREATE":
			handleOrderCreate(requestMap, conn)
		case "GET":
			handleOrderGet(requestMap, conn)
		case "ADD DISCOUNT":
			handleAddDiscount(requestMap, conn)
		case "DELETE DISCOUNT":
			handleDeleteDiscount(requestMap, conn)
		case "CLOSE":
			handleOrderClose(requestMap, conn)
		case "UPDATE":
			handleOrderUpdate(requestMap, conn)
		}
	}
}

func handleConnection(conn net.Conn) {
	requestMap := make(map[string]string)
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&requestMap)
	if err != nil {
		log.Println("Handle connection: Error on decoding json: ", err)
		conn.Close()
		return
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

	goposPrintserviceFont = os.Getenv("GOPOS_PRINTSERVICE_FONT")
	if goposPrintserviceFont == "" {
		log.Fatal("GOPOS_PRINTSERVICE_FONT is not set")
	}

	goposCheckPath = os.Getenv("GOPOS_CHECK_PATH")
	if goposCheckPath == "" {
		log.Fatal("GOPOS_CHECK_PATH is not set")
	}

	goposOrderPath = os.Getenv("GOPOS_ORDER_PATH")
	if goposOrderPath == "" {
		log.Fatal("GOPOS_ORDER_PATH is not set")
	}

	goposCheckWidth = os.Getenv("GOPOS_CHECK_WIDTH")
	if goposCheckWidth == "" {
		log.Fatal("GOPOS_CHECK_WIDTH is not set")
	}

	goposCheckHeight = os.Getenv("GOPOS_CHECK_HEIGHT")
	if goposCheckHeight == "" {
		log.Fatal("GOPOS_CHECK_HEIGHT is not set")
	}

	goposOrderWidth = os.Getenv("GOPOS_ORDER_WIDTH")
	if goposOrderWidth == "" {
		log.Fatal("GOPOS_ORDER_WIDTH is not set")
	}

	goposOrderHeight = os.Getenv("GOPOS_ORDER_HEIGHT")
	if goposOrderHeight == "" {
		log.Fatal("GOPOS_ORDER_HEIGHT is not set")
	}

	dbConn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/gopos?parseTime=true", goposSQLUser,
		goposSQLPassword))
	_, err = dbConn.Exec("SET CHARSET utf8")
	if err != nil {
		log.Fatal("Error on setting charset: ", err)
	}

	if err != nil {
		log.Fatal("Error on opening database", err)
	}

	listener, err := net.Listen("tcp", ":"+goposServerPort)
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
