package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"strings"

	"strconv"

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

func handleOrdersGet(requestMap map[string]string, conn net.Conn) {
	var orderId, totalDishes int64
	var orderString string
	var price, discount float64
	var orderParts []string
	var splitedOrderParts [][]string

	dbConn.QueryRow(fmt.Sprintf("SELECT current_order FROM tables WHERE number=%s",
		requestMap["table_number"])).Scan(&orderId)
	dbConn.QueryRow(fmt.Sprintf("SELECT order_string, price, discount FROM orders "+
		"WHERE id=%d", orderId)).Scan(&orderString, &price, &discount)

	orderParts = strings.Split(orderString, " ")
	splitedOrderParts = make([][]string, len(orderParts))
	for i, part := range orderParts {
		splitedOrderParts[i] = strings.Split(part, ":")
		tmp, _ := strconv.ParseInt(splitedOrderParts[i][1], 0, 32)
		totalDishes += tmp
	}

	encoder := json.NewEncoder(conn)
	responseMap := make(map[string]string)
	responseMap["total_dishes"] = fmt.Sprintf("%d", totalDishes)
	err := encoder.Encode(responseMap)

	for _, part := range splitedOrderParts {
		var dishId int
		var dishName string
		var dishPrice float64
		dbConn.QueryRow(fmt.Sprintf("SELECT id, nam	, price FROM dishes where id=%s", part[0])).
			Scan(&dishId, &dishName, &dishPrice)

		count, _ := strconv.ParseInt(part[1], 0, 32)
		for j := int64(0); j < count; j++ {
			responseMap = make(map[string]string)
			responseMap["id"] = fmt.Sprintf("%d", dishId)
			responseMap["name"] = dishName
			responseMap["price"] = fmt.Sprintf("%f", dishPrice)
			err := encoder.Encode(responseMap)
			if err != nil {
				log.Fatal("Error on encode resoponse map: ", err)
			}
		}
	}

	responseMap = make(map[string]string)
	responseMap["price"] = fmt.Sprintf("%f", (1.0-discount)*price)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encode resoponse map: ", err)
	}
}
