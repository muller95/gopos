//this file contains functions to handle requests from dish group
package main

import (
	"encoding/json"
	"fmt"
	_"github.com/go-sql-driver/mysql"
	"log"
	"net"
)

func handleDishGet(requestMap map[string]string, conn net.Conn) {
	rows, err := dbConn.Query(fmt.Sprintf("SELECT id, name, price FROM dishes " +
	"WHERE category_id=%s ORDER BY id ASC", requestMap["category_id"]))
	if err != nil {
		log.Fatal("Error on getting categories ids: ", err)
	}

	encoder := json.NewEncoder(conn)
	for rows.Next() {
		var id int
		var name string
		var price float64

		err = rows.Scan(&id, &name, &price)
		if err != nil {
			log.Fatal("Error on handling sql response")
		}

		responseMap := make(map[string]string)
		responseMap["id"] = fmt.Sprintf("%d", id)
		responseMap["name"] = fmt.Sprintf("%s", name)
		responseMap["price"] = fmt.Sprintf("%f", price)
		err = encoder.Encode(responseMap)
		if err != nil {
			log.Fatal("Error on encode request map: ", err)
		}
	}
}

func handleDishAdd(requestMap map[string]string, conn net.Conn) {
	rows, err := dbConn.Query(fmt.Sprintf("SELECT * FROM dishes WHERE name='%s' AND " +
		"category_id=%s", requestMap["name"], requestMap["category_id"]))
	if err != nil {
		log.Fatal("Error on getting category names: ", err)
	}

	responseMap := make(map[string]string)
	if rows.Next() {
		responseMap["result"] = "ERR"
		responseMap["error"] = "В выбранной категории блюдо с таким названием уже есть"
		encoder := json.NewEncoder(conn)
		err := encoder.Encode(responseMap)
		if err != nil {
			log.Fatal("Error on encoding response json: ", err)
		}
		return
	}

	if err != nil {
		log.Fatal("Error on getting tables: ", err)
	}
	id := 0
	rows, err = dbConn.Query("SELECT id from dishes ORDER BY id ASC")
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

	_, err = dbConn.Exec(fmt.Sprintf("INSERT INTO dishes VALUES(%d, '%s', %s, %s)", id,
		requestMap["name"], requestMap["price"], requestMap["category_id"]))
	if err != nil {
		log.Fatal("Error on inserting new category: ", err)
	}

	responseMap["result"] = "OK"
	responseMap["id"] = fmt.Sprintf("%d", id)
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", err)
	}

	conn.Close()
}

func handleDishDelete(requestMap map[string]string, conn net.Conn) {
	_, err := dbConn.Exec(fmt.Sprintf("DELETE FROM dishes WHERE id=%s;",
		requestMap["id"]))
	if err != nil {
		log.Fatal("Error on deleting category: ", err)
	}

	responseMap := make(map[string]string)
	responseMap["result"] = fmt.Sprintf("OK")
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", err)
	}
}
