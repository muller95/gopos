package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"strings"

	"strconv"

	"os"

	_ "github.com/go-sql-driver/mysql"
)

var mutex = &sync.Mutex{}

func handleOrderCreate(requestMap map[string]string, conn net.Conn) {
	var orderId int64
	var orderString string
	var orderParts []string

	mutex.Lock()

	rows, err := dbConn.Query(fmt.Sprintf("SELECT * FROM tables WHERE number=%s and "+
		"current_order=-1", requestMap["table_number"]))
	if err != nil {
		log.Fatal("Error on getting tables: ", err)
	}

	dbConn.QueryRow(fmt.Sprintf("SELECT current_order FROM tables WHERE number=%s",
		requestMap["table_number"])).Scan(&orderId)

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

	orderId = 0
	rows, err = dbConn.Query("SELECT id from orders ORDER BY id ASC")
	for ; rows.Next(); orderId++ {
		var currId int
		err := rows.Scan(&currId)
		if err != nil {
			log.Fatal("Error on handling sql response")
		}

		if orderId != int64(currId) {
			break
		}
	}

	tx, err := dbConn.Begin()
	_, err = tx.Exec(fmt.Sprintf("INSERT INTO orders VALUES(%d, '%s', %s, 0.0);", orderId,
		requestMap["order_string"], requestMap["price"]))
	if err != nil {
		log.Fatal("Error on inserting order: ", err)
	}

	_, err = tx.Exec(fmt.Sprintf("UPDATE tables SET current_order=%d WHERE number=%s", orderId,
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

	file, err := os.Create(fmt.Sprintf("%s%s_%d.html", goposOrderPath,
		requestMap["table_number"], orderId))
	_, err = file.Write([]byte("<html>"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte(fmt.Sprintf("<head><meta charset=\"utf-8\"></meta><style> th, td {"+
		"border: 1px solid black;} table { border: 1px solid black; } "+
		" @font-face { font-family: \"GoposFont\"; src: url(%s); -fs-pdf-font-embed: embed; "+
		"-fs-pdf-font-encoding: Identity-H; } * { font-family: GoposFont; }"+
		"@page { width: %s; } </style></head>",
		goposPrintserviceFont, goposOrderWidth)))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte("<body style=\"margin-top:1cm\">"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte(fmt.Sprintf("<H3>Новый заказ: %d. Столик: %s</H3>", orderId,
		requestMap["table_number"])[:]))

	_, err = file.Write([]byte("<table style=\"border-style:solid;width:100%\">"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte("<tr><th>Блюдо</th><th>Цена за единицу</th><th>Количество</th>" +
		"<th>Итого</th></tr>"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	orderString = requestMap["order_string"]
	orderParts = strings.Split(orderString, " ")
	for _, part := range orderParts {
		var dishId int
		var dishName string
		var dishPrice float64

		orderPartInfo := strings.Split(part, ":")

		dbConn.QueryRow(fmt.Sprintf("SELECT id, name, price FROM dishes where id=%s",
			orderPartInfo[0])).Scan(&dishId, &dishName, &dishPrice)

		count, _ := strconv.ParseInt(orderPartInfo[1], 0, 32)
		_, err = file.Write([]byte(fmt.Sprintf("<tr><td>%s</td><td>%.2f</td><td>%d</td>"+
			"<td>%.2f</td></tr>", dishName, dishPrice, count, dishPrice*float64(count))))
		if err != nil {
			log.Fatal("Error on writing to file: ", err)
		}
	}

	_, err = file.Write([]byte("</table>"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte("</body>"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte("</html>"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	file.Sync()
	file.Close()

	responseMap["result"] = fmt.Sprintf("OK")
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encode resoponse map: ", err)
	}

	mutex.Unlock()
}

func handleOrderGet(requestMap map[string]string, conn net.Conn) {
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
		dbConn.QueryRow(fmt.Sprintf("SELECT id, name, price FROM dishes where id=%s", part[0])).
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
	responseMap["price"] = fmt.Sprintf("%f", price)
	responseMap["discount"] = fmt.Sprintf("%f", discount)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encode resoponse map: ", err)
	}
}

func handleAddDiscount(requestMap map[string]string, conn net.Conn) {
	var orderId int
	var discount float64

	dbConn.QueryRow(fmt.Sprintf("SELECT current_order FROM tables WHERE number=%s",
		requestMap["table_number"])).Scan(&orderId)
	rows, err := dbConn.Query(fmt.Sprintf("SELECT discount FROM cards where number='%s'",
		requestMap["card_number"]))
	if err != nil {
		log.Fatal("Error on getting card numbers")
	}

	responseMap := make(map[string]string)

	if !rows.Next() {
		responseMap["result"] = "ERR"
		responseMap["error"] = "Нет карты с таким номером"
		encoder := json.NewEncoder(conn)
		err := encoder.Encode(responseMap)
		if err != nil {
			log.Fatal("Error on encoding response json: ", err)
		}

		return
	}

	rows.Scan(&discount)
	_, err = dbConn.Exec(fmt.Sprintf("UPDATE orders SET discount=%f WHERE id=%d",
		discount, orderId))
	if err != nil {
		log.Fatal("Error on setting discount")
	}

	responseMap["result"] = "OK"
	responseMap["discount"] = fmt.Sprintf("%f", discount)
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encoding response json: ", err)
	}

}

func handleDeleteDiscount(requestMap map[string]string, conn net.Conn) {
	var orderId int

	dbConn.QueryRow(fmt.Sprintf("SELECT current_order FROM tables WHERE number=%s",
		requestMap["table_number"])).Scan(&orderId)

	_, err := dbConn.Exec(fmt.Sprintf("UPDATE orders SET discount=0.0 WHERE id=%d", orderId))
	if err != nil {
		log.Fatal("Error on deleting discount")
	}

	responseMap := make(map[string]string)
	responseMap["result"] = "OK"
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encoding response json: ", err)
	}
}

func handleOrderClose(requestMap map[string]string, conn net.Conn) {
	var orderId int64
	var orderString string
	var price, discount float64
	var orderParts []string

	dbConn.QueryRow(fmt.Sprintf("SELECT current_order FROM tables WHERE number=%s",
		requestMap["table_number"])).Scan(&orderId)
	dbConn.QueryRow(fmt.Sprintf("SELECT order_string, price, discount FROM orders "+
		"WHERE id=%d", orderId)).Scan(&orderString, &price, &discount)

	file, err := os.Create(fmt.Sprintf("%s%s_%d.html", goposCheckPath,
		requestMap["table_number"], orderId))
	_, err = file.Write([]byte("<html>"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte(fmt.Sprintf("<head><meta charset=\"utf-8\"></meta><style> th, td {"+
		"border: 1px solid black;} table { border: 1px solid black; } "+
		" @font-face { font-family: \"GoposFont\"; src: url(%s); -fs-pdf-font-embed: embed; "+
		"-fs-pdf-font-encoding: Identity-H; } * { font-family: GoposFont; }"+
		"@page { width: %s; } </style></head>",
		goposPrintserviceFont, goposCheckWidth)))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte("<body style=\"margin-top:1cm\">"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte(fmt.Sprintf("<H3>Заказ: %d. Столик: %s</H3>", orderId,
		requestMap["table_number"])[:]))

	_, err = file.Write([]byte("<table style=\"border-style:solid;width:100%\">"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte("<tr><th>Блюдо</th><th>Цена за единицу</th><th>Количество</th>" +
		"<th>Итого</th></tr>"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	orderParts = strings.Split(orderString, " ")
	for _, part := range orderParts {
		var dishId int
		var dishName string
		var dishPrice float64

		orderPartInfo := strings.Split(part, ":")

		dbConn.QueryRow(fmt.Sprintf("SELECT id, name, price FROM dishes where id=%s",
			orderPartInfo[0])).Scan(&dishId, &dishName, &dishPrice)

		count, _ := strconv.ParseInt(orderPartInfo[1], 0, 32)
		_, err = file.Write([]byte(fmt.Sprintf("<tr><td>%s</td><td>%.2f</td><td>%d</td>"+
			"<td>%.2f</td></tr>", dishName, dishPrice, count, dishPrice*float64(count))))
		if err != nil {
			log.Fatal("Error on writing to file: ", err)
		}
	}

	_, err = file.Write([]byte("</table>"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte(fmt.Sprintf("<b>Итого:</b> %.2f <br/> <b>Скидка:</b> %.2f <br/>"+
		"<b>Итого со скидкой:</b> %.2f", price, discount, (1.0-discount)*price)))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte("</body>"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte("</html>"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	file.Sync()
	file.Close()

	_, err = dbConn.Exec(fmt.Sprintf("UPDATE tables SET current_order=-1 WHERE number=%s",
		requestMap["table_number"]))
	if err != nil {
		log.Fatal("Error on setting discount")
	}

	responseMap := make(map[string]string)
	responseMap["result"] = "OK"
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encoding response json: ", err)
	}
}

func handleOrderUpdate(requestMap map[string]string, conn net.Conn) {
	var orderId int
	var prevOrderString, newOrderString string

	mutex.Lock()

	dbConn.QueryRow(fmt.Sprintf("SELECT current_order FROM tables WHERE number=%s",
		requestMap["table_number"])).Scan(&orderId)
	dbConn.QueryRow(fmt.Sprintf("SELECT order_string FROM orders "+
		"WHERE id=%d", orderId)).Scan(&prevOrderString)
	newOrderString = requestMap["order_string"]

	file, err := os.Create(fmt.Sprintf("%s%s_%d.html", goposOrderPath,
		requestMap["table_number"], orderId))
	_, err = file.Write([]byte("<html>"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte(fmt.Sprintf("<head><meta charset=\"utf-8\"></meta><style> th, td {"+
		"border: 1px solid black;} table { border: 1px solid black; } "+
		" @font-face { font-family: \"GoposFont\"; src: url(%s); -fs-pdf-font-embed: embed; "+
		"-fs-pdf-font-encoding: Identity-H; } * { font-family: GoposFont; }"+
		"@page { width: %s; } </style></head>",
		goposPrintserviceFont, goposOrderWidth)))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte("<body style=\"margin-top:1cm\">"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte(fmt.Sprintf("<H3>Обновление заказа: %d. Столик: %s</H3>", orderId,
		requestMap["table_number"])[:]))

	_, err = file.Write([]byte("<table style=\"border-style:solid;width:100%\">"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte("<tr><th>Блюдо</th><th>Цена за единицу</th><th>Количество</th>" +
		"<th>Итого</th></tr>"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	prevOrderParts := strings.Split(prevOrderString, " ")
	newOrderParts := strings.Split(newOrderString, " ")
	for _, newPart := range newOrderParts {
		var dishId int
		var dishName string
		var dishPrice float64

		newOrderPartInfo := strings.Split(newPart, ":")

		dbConn.QueryRow(fmt.Sprintf("SELECT id, name, price FROM dishes where id=%s",
			newOrderPartInfo[0])).Scan(&dishId, &dishName, &dishPrice)
		newCount, _ := strconv.ParseInt(newOrderPartInfo[1], 0, 32)
		for _, prevPart := range prevOrderParts {
			prevOrderPartInfo := strings.Split(prevPart, ":")
			if prevOrderPartInfo[0] == newOrderPartInfo[0] {
				prevCount, _ := strconv.ParseInt(prevOrderPartInfo[1], 0, 32)
				newCount -= prevCount
			}
		}

		if newCount > 0 {
			_, err = file.Write([]byte(fmt.Sprintf("<tr><td>%s</td><td>%.2f</td><td>%d</td>"+
				"<td>%.2f</td></tr>", dishName, dishPrice, newCount, dishPrice*float64(newCount))))
			if err != nil {
				log.Fatal("Error on writing to file: ", err)
			}
		}
	}

	_, err = file.Write([]byte("</table>"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte("</body>"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = file.Write([]byte("</html>"[:]))
	if err != nil {
		log.Fatal("Error on writing to file: ", err)
	}

	_, err = dbConn.Exec(fmt.Sprintf("UPDATE orders SET order_string='%s', price=%s WHERE id=%d",
		requestMap["order_string"], requestMap["price"], orderId))

	responseMap := make(map[string]string)
	responseMap["result"] = fmt.Sprintf("OK")
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(responseMap)
	if err != nil {
		log.Fatal("Error on encode resoponse map: ", err)
	}

	mutex.Unlock()
}
