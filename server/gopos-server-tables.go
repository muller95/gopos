//this file contains functions to handle requests from tables group
package main

import (
	"encoding/json"
	"fmt"
	_"github.com/go-sql-driver/mysql"
	"log"
	"net"
)

func handleTableGet(conn net.Conn) {
	rows, err := dbConn.Query("SELECT * from tables ORDER BY number ASC")
	if err != nil {
		log.Fatal("Error on getting table numbers: ", err)
	}

	encoder := json.NewEncoder(conn)
	for rows.Next() {
		var number int

		err = rows.Scan(&number)
		if err != nil {
			log.Fatal("Error on handling sql response")
		}

		responseMap := make(map[string]string)
		responseMap["number"] = fmt.Sprintf("%d", number)
		err = encoder.Encode(responseMap)
		if err != nil {
			log.Fatal("Error on encode request map: ", err)
		}
	}

	conn.Close()
}

func handleTableAdd(requestMap map[string]string, conn net.Conn) {
	rows, err := dbConn.Query(fmt.Sprintf("SELECT * from tables where number=%s",
		requestMap["number"]))
	responseMap := make(map[string]string)
	if err != nil {
		log.Fatal("Error on getting tables: ", err)
	}

	if rows.Next() {
		responseMap := make(map[string]string)
		responseMap["result"] = "ERR"
		responseMap["error"] = "Столик с таким номером уже есть"
		encoder := json.NewEncoder(conn)
		err := encoder.Encode(responseMap)
		if err != nil {
			log.Fatal("Error on encoding response json: ", err)
		}

		return
	}


	_, err = dbConn.Exec(fmt.Sprintf("INSERT INTO tables VALUES(%s)", requestMap["number"]))
	if err != nil {
		log.Fatal("Error on inserting new table: ", err)
	}

	responseMap["result"] = "OK"
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", err)
	}

	conn.Close()
}

func handleTableDelete(requestMap map[string]string, conn net.Conn) {
	_, err := dbConn.Exec(fmt.Sprintf("DELETE FROM tables where number=%s;",
		requestMap["number"]))
	if err != nil {
		log.Fatal("Error on deleting table: ", err)
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
