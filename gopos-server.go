package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_"github.com/go-sql-driver/mysql"
	"log"
	"net"
	"os"
	"time"
)

var goposServerPassword, goposServerPort, goposSQLUser, goposSQLPassword string
var dbConn *sql.DB

func handleWorkerGet(conn net.Conn) {
	rows, err := dbConn.Query("SELECT * from workers ORDER BY id ASC")
	if err != nil {
		log.Fatal("Error on getting workers ids: ", err)
	}

	encoder := json.NewEncoder(conn)
	for rows.Next() {
		var id int
		var name string
		var date time.Time

		err = rows.Scan(&id, &name, &date)
		if err != nil {
			log.Fatal("Error on handling sql response")
		}

		responseMap := make(map[string]string)
		responseMap["id"] = fmt.Sprintf("%d", id)
		responseMap["name"] = fmt.Sprintf("%s", name)
		responseMap["date"] = fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(),
			date.Day())
		err = encoder.Encode(responseMap)
		if err != nil {
			log.Fatal("Error on encode request map: ", err)
		}
//		conn.Close()
	}

	conn.Close()
}

func handleWorkerAdd(requestMap map[string]string, conn net.Conn) {
	id := 0
	rows, err := dbConn.Query("SELECT id from workers ORDER BY id ASC")
	responseMap := make(map[string]string)
	if err != nil {
		log.Fatal("Error on getting workers ids: ", err)
	}

	for ;rows.Next(); id++ {
		var currId int
		err := rows.Scan(&currId)
		if err != nil {
			log.Fatal("Error on handling sql response")
		}

		if id != currId {
			break
		}
	}

	_, err = dbConn.Exec(fmt.Sprintf("INSERT INTO workers VALUES(%d, '%s', '%s')", id,
		requestMap["name"], requestMap["date"]))
	if err != nil {
		log.Fatal("Error on inserting new worker: ", err)
	}

	responseMap["id"] = fmt.Sprintf("%d", id)
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", err)
	}

	conn.Close()
}

func handleWorkerDelete(requestMap map[string]string, conn net.Conn) {
	_, err := dbConn.Exec(fmt.Sprintf("DELETE FROM workers where id=%s;", requestMap["id"]))
	if err != nil {
		log.Fatal("Error on deleting worker: ", err)
	}

	responseMap := make(map[string]string)
	responseMap["result"] = fmt.Sprintf("OK")
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", err)
	}

	conn.Close()
}

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
		responseMap["id"] = "-1"
		responseMap["error"] = "Неправильный пароль"
		encoder := json.NewEncoder(conn)
		err := encoder.Encode(responseMap)
		if err != nil {
			log.Fatal("Error on encoding response json: ", err)
		}

		return
	}

	handleRequestGroup(requestMap, conn)
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

	dbConn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/gopos?parseTime=true", goposSQLUser,
		goposSQLPassword))

	if err != nil {
		log.Fatal("Error on opening database", err)
	}

//	handleWorkerGet()

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
