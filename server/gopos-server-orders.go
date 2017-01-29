package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var mutex = &sync.Mutex{}

func handleOrderCreate(requestMap map[string]string, conn net.Conn) {
	mutex.Lock()

	rows, err := dbConn.Query(fmt.Sprintf("SELECT * FROM tables WHERE number=%s and "+
		"current_order=-1", requestMap["table_number"]))
	if err != nil {
		log.Fatal("Error on getting tables: ", err)
	}

	responseMap := make(map[string]string)
	if !rows.Next() {
		responseMap["result"] = "ERR"
		responseMap["error"] = "Столик уже занят."
		encoder := json.NewEncoder(conn)
		err := encoder.Encode(responseMap)
		if err != nil {
			log.Fatal("Error on encoding response json: ", err)
		}

		return
	}

	id := 0
	rows, err = dbConn.Query("SELECT id from orders ORDER BY id ASC")
	for ; rows.Next(); id++ {
		var currId int
		err := rows.Scan(&currId)
		if err != nil {
			log.Fatal("Error on handling sql response")
		}

		if id != currId {
			break
		}
	}

	tx, err := dbConn.Begin()
	fmt.Printf("INSERT INTO orders VALUES(%d, '%s', %s, 0.0;)\n", id,
		requestMap["order_string"], requestMap["price"])
	_, err = tx.Exec(fmt.Sprintf("INSERT INTO orders VALUES(%d, '%s', %s, 0.0);", id,
		requestMap["order_string"], requestMap["price"]))
	if err != nil {
		log.Fatal("Error on inserting order: ", err)
	}

	_, err = tx.Exec(fmt.Sprintf("UPDATE tables SET current_order=%d WHERE number=%s", id,
		requestMap["table_number"]))
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			log.Fatal("Error on rollback order transaction: ", err)
		}
		log.Fatal("Error on updating table current order: ", err)
	}

	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			log.Fatal("Error on rollback order transaction: ", err)
		}
		log.Fatal("Error on commiting order transaction: ", err)
	}

	responseMap["result"] = fmt.Sprintf("OK")
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encode resoponse map: ", err)
	}

	mutex.Unlock()
}

/*func handleOrdersGet(requestMap map[string]string, conn net.Conn) {
	_, err := dbConn.Query(fmt.Sprntf("SELECT orders ")
}*/
