//this file contains functions to handle requests from card group
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"
)

func handleCardGet(conn net.Conn) {
	rows, err := dbConn.Query("SELECT * FROM cards")
	if err != nil {
		log.Fatal("Error on getting categories ids: ", err)
	}

	encoder := json.NewEncoder(conn)
	for rows.Next() {
		var cardNumber string
		var holderName string
		var discount float64

		err = rows.Scan(&cardNumber, &holderName, &discount)
		if err != nil {
			log.Fatal("Error on handling sql response")
		}

		responseMap := make(map[string]string)
		responseMap["number"] = cardNumber
		responseMap["holder_name"] = holderName
		responseMap["discount"] = fmt.Sprintf("%f", discount)
		err = encoder.Encode(responseMap)
		if err != nil {
			log.Fatal("Error on encode request map: ", err)
		}
	}
}

func handleCardAdd(requestMap map[string]string, conn net.Conn) {
	rows, err := dbConn.Query(fmt.Sprintf("SELECT * FROM cards WHERE number='%s'",
		requestMap["number"]))
	if err != nil {
		log.Fatal("Error on getting tables: ", err)
	}

	responseMap := make(map[string]string)
	if rows.Next() {
		responseMap["result"] = "ERR"
		responseMap["error"] = "Карта с таким номером уже есть"
		encoder := json.NewEncoder(conn)
		err := encoder.Encode(responseMap)
		if err != nil {
			log.Fatal("Error on encoding response json: ", err)
		}

		return
	}

	_, err = dbConn.Exec(fmt.Sprintf("INSERT INTO cards VALUES('%s', '%s', %s)",
		requestMap["number"], requestMap["holder_name"], requestMap["discount"]))
	if err != nil {
		log.Fatal("Error on inserting new table: ", err)
	}

	responseMap["result"] = "OK"
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", err)
	}
}

func handleCardDelete(requestMap map[string]string, conn net.Conn) {
	_, err := dbConn.Exec(fmt.Sprintf("DELETE FROM cards WHERE number='%s';",
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
}
