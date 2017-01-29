//this file contains functions to handle requests from table group
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
)

func handleTableGet(conn net.Conn) {
	rows, err := dbConn.Query("SELECT * FROM tables ORDER BY number ASC")
	if err != nil {
		log.Fatal("Error on getting table numbers: ", err)
	}

	encoder := json.NewEncoder(conn)
	for rows.Next() {
		var number int
		var currentOrder int

		err = rows.Scan(&number, &currentOrder)
		if err != nil {
			log.Fatal("Error on handling sql response")
		}

		responseMap := make(map[string]string)
		responseMap["number"] = fmt.Sprintf("%d", number)
		responseMap["current_order"] = fmt.Sprintf("%d", currentOrder)
		err = encoder.Encode(responseMap)
		if err != nil {
			log.Fatal("Error on encode request map: ", err)
		}
	}
}

func handleTableAdd(requestMap map[string]string, conn net.Conn) {
	rows, err := dbConn.Query(fmt.Sprintf("SELECT * FROM tables WHERE number=%s",
		requestMap["number"]))
	if err != nil {
		log.Fatal("Error on getting tables: ", err)
	}

	responseMap := make(map[string]string)
	if rows.Next() {
		responseMap["result"] = "ERR"
		responseMap["error"] = "Столик с таким номером уже есть"
		encoder := json.NewEncoder(conn)
		err := encoder.Encode(responseMap)
		if err != nil {
			log.Fatal("Error on encoding response json: ", err)
		}

		return
	}

	_, err = dbConn.Exec(fmt.Sprintf("INSERT INTO tables VALUES(%s, -1)", requestMap["number"]))
	if err != nil {
		log.Fatal("Error on inserting new table: ", err)
	}

	responseMap["result"] = "OK"
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encode response map: ", err)
	}
}

func handleTableDelete(requestMap map[string]string, conn net.Conn) {
	_, err := dbConn.Exec(fmt.Sprintf("DELETE FROM tables WHERE number=%s;",
		requestMap["number"]))
	if err != nil {
		log.Fatal("Error on deleting table: ", err)
	}

	responseMap := make(map[string]string)
	responseMap["result"] = fmt.Sprintf("OK")
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encode resoponse map: ", err)
	}
}
