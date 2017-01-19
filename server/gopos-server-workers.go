//this file contains functions to handle requests from worker groups
package main

import (
	"encoding/json"
	"fmt"
	_"github.com/go-sql-driver/mysql"
	"log"
	"net"
	"time"
)

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
	}
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

	responseMap["result"] = "OK"
	responseMap["id"] = fmt.Sprintf("%d", id)
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", err)
	}
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
}
