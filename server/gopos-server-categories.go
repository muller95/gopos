//this file contains functions to handle requests from category groups
package main

import (
	"encoding/json"
	"fmt"
	_"github.com/go-sql-driver/mysql"
	"log"
	"net"
)

func handleCategoryGet(conn net.Conn) {
	rows, err := dbConn.Query("SELECT * FROM categories ORDER BY id ASC")
	if err != nil {
		log.Fatal("Error on getting categories ids: ", err)
	}

	encoder := json.NewEncoder(conn)
	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal("Error on handling sql response")
		}

		responseMap := make(map[string]string)
		responseMap["id"] = fmt.Sprintf("%d", id)
		responseMap["name"] = fmt.Sprintf("%s", name)
		err = encoder.Encode(responseMap)
		if err != nil {
			log.Fatal("Error on encode request map: ", err)
		}
	}
}

func handleCategoryAdd(requestMap map[string]string, conn net.Conn) {
	rows, err := dbConn.Query(fmt.Sprintf("SELECT * FROM categories WHERE name='%s'",
		requestMap["name"]))
	if err != nil {
		log.Fatal("Error on getting category names: ", err)
	}

	responseMap := make(map[string]string)
	if rows.Next() {
		responseMap["result"] = "ERR"
		responseMap["error"] = "Категория с таким названием уже есть"
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
	rows, err = dbConn.Query("SELECT id FROM categories ORDER BY id ASC")
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

	_, err = dbConn.Exec(fmt.Sprintf("INSERT INTO categories VALUES(%d, '%s')", id,
		requestMap["name"]))
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

func handleCategoryDelete(requestMap map[string]string, conn net.Conn) {
	tx, err := dbConn.Begin()
	if err != nil {
		log.Fatalf("Error on begining category delete transaction: ", err)
	}

	_, err = tx.Exec(fmt.Sprintf("DELETE FROM categories WHERE id=%s;",
		requestMap["id"]))
	if err != nil {
		log.Fatal("Error on deleting category: ", err)
	}

	_, err = tx.Exec(fmt.Sprintf("DELETE FROM dishes WHERE category_id=%s", requestMap["id"]))
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			log.Fatalf("Error on tollback category transaction: ", err)
		}
		log.Fatalf("Error on deleting dishes in category: ", err)
	}

	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			log.Fatalf("Error on tollback category transaction: ", err)
		}
		log.Fatalf("Error on commiting category transaction: ", err)
	}

	responseMap := make(map[string]string)
	responseMap["result"] = fmt.Sprintf("OK")
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", err)
	}


}
