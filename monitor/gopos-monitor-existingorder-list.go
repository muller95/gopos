package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"

	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var existingOrderTreeView *gtk.TreeView
var existingOrderListStore *gtk.ListStore
var existingOrderPriceLabel *gtk.Label
var existingOrderListWindow *gtk.Window

var existingOrderPrice float64

const (
	COLUMN_EXISTING_ORDER_LIST_DISH_ID = iota
	COLUMN_EXISTING_ORDER_LIST_DISH_NAME
	COLUMN_EXISTING_ORDER_LIST_DISH_PRICE
)

func createExistingOrderTreeView() {
	var err error
	existingOrderTreeView, err = gtk.TreeViewNew()

	if err != nil {
		log.Fatal("Unable to create tables tree view: ", err)
	}

	existingOrderTreeView.AppendColumn(createColumn("ID", COLUMN_EXISTING_ORDER_LIST_DISH_ID))
	existingOrderTreeView.AppendColumn(createColumn("Название блюда", COLUMN_EXISTING_ORDER_LIST_DISH_NAME))
	existingOrderTreeView.AppendColumn(createColumn("Цена", COLUMN_EXISTING_ORDER_LIST_DISH_PRICE))

	existingOrderListStore, err = gtk.ListStoreNew(glib.TYPE_INT, glib.TYPE_STRING, glib.TYPE_DOUBLE)
	if err != nil {
		log.Fatal("Unable to create tables list store: ", err)
	}

	existingOrderTreeView.SetModel(existingOrderListStore)
}

func existingOrderAddRow(id int, name string, price float64) {
	iter := existingOrderListStore.Append()
	fmt.Println("here")
	err := existingOrderListStore.Set(iter, []int{COLUMN_EXISTING_ORDER_LIST_DISH_ID, COLUMN_EXISTING_ORDER_LIST_DISH_NAME,
		COLUMN_EXISTING_ORDER_LIST_DISH_PRICE}, []interface{}{id, name, price})

	if err != nil {
		log.Fatal("Unable to add tables row: ", err)
	}
}

func dishExistingOrderDeleteSelectedButtonClicked(btn *gtk.Button, passwordEntry *gtk.Entry) {
	password, err := passwordEntry.GetText()
	if err != nil {
		log.Fatal("Error on getting password")
	}

	if password != goposServerPassword {
		messageDialog := gtk.MessageDialogNew(mainWindow, gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING,
			gtk.BUTTONS_OK, "Неверный пароль.")
		messageDialog.Run()
		messageDialog.Destroy()
		return
	}

	selection, err := existingOrderTreeView.GetSelection()
	if err != nil {
		log.Fatal("Error on getting existing order selection")
	}

	rows := selection.GetSelectedRows(existingOrderListStore)
	if rows == nil {
		return
	}
	path := rows.Data().(*gtk.TreePath)
	iter, err := existingOrderListStore.GetIter(path)
	if err != nil {
		log.Fatal("Error on getting iter: ", err)
	}
	value, err := existingOrderListStore.GetValue(iter, 2)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	existingOrderPrice -= value.GetDouble()
	existingOrderPriceLabel.SetText(fmt.Sprintf("Цена: %.2f", existingOrderPrice))
	existingOrderListStore.Remove(iter)
}

func getOrder() {
	var responseMap map[string]string

	existingOrderListStore.Clear()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
		goposServerPort))

	if err != nil {
		log.Fatal("Unable to connect to server")
	}

	requestMap := make(map[string]string)
	requestMap["group"] = "ORDER"
	requestMap["action"] = "GET"
	requestMap["password"] = goposServerPassword
	requestMap["table_number"] = fmt.Sprintf("%d", orderedTableNumber)
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(requestMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", requestMap)
	}

	decoder := json.NewDecoder(conn)
	responseMap = make(map[string]string)
	err = decoder.Decode(&responseMap)
	if err != nil {
		log.Fatal("Error on decoding response: ", err)
	}
	totalDishes, _ := strconv.Atoi(responseMap["total_dishes"])
	for i := 0; i < totalDishes; i++ {
		responseMap = make(map[string]string)
		err = decoder.Decode(&responseMap)
		fmt.Println(responseMap)
		if err != nil {
			log.Fatal("Error on decoding response: ", err)
		}

		id, _ := strconv.Atoi(responseMap["id"])
		price, _ := strconv.ParseFloat(responseMap["price"], 64)
		existingOrderAddRow(id, responseMap["name"], price)
	}

	responseMap = make(map[string]string)
	err = decoder.Decode(&responseMap)
	if err != nil {
		log.Fatal("Error on decoding response: ", err)
	}
	existingOrderPriceLabel.SetText("Цена: " + responseMap["price"])
}

func existingOrderCloseButtonClicked(btn *gtk.Button, passwordEntry *gtk.Entry) {
	password, err := passwordEntry.GetText()
	if err != nil {
		log.Fatal("Error on getting password")
	}

	if password != goposServerPassword {
		messageDialog := gtk.MessageDialogNew(mainWindow, gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING,
			gtk.BUTTONS_OK, "Неверный пароль.")
		messageDialog.Run()
		messageDialog.Destroy()
		return
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
		goposServerPort))

	if err != nil {
		log.Fatal("Unable to connect to server")
	}

	requestMap := make(map[string]string)
	requestMap["group"] = "ORDER"
	requestMap["action"] = "CLOSE"
	requestMap["password"] = goposServerPassword
	requestMap["table_number"] = fmt.Sprintf("%d", orderedTableNumber)
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(requestMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", requestMap)
	}

	decoder := json.NewDecoder(conn)
	responseMap := make(map[string]string)
	err = decoder.Decode(&responseMap)
	if err != nil {
		log.Fatal("Error on decoding response: ", err)
	}

	if responseMap["result"] == "OK" {
		existingOrderWindow.Destroy()
		existingOrderListWindow.Destroy()
	}
}

func discountAddButtonClicked(btn *gtk.Button, cardNumberEntry *gtk.Entry) {
	cardNumber, err := cardNumberEntry.GetText()
	if err != nil {
		log.Fatal("Unable to get card number entry text: ", err)
	}
	cardNumber = strings.Trim(cardNumber, " ")

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
		goposServerPort))

	if err != nil {
		log.Fatal("Unable to connect to server")
	}

	requestMap := make(map[string]string)
	requestMap["group"] = "ORDER"
	requestMap["action"] = "ADD DISCOUNT"
	requestMap["password"] = goposServerPassword
	requestMap["card_number"] = cardNumber
	requestMap["table_namber"] = fmt.Sprintf("%d", orderedTableNumber)
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(requestMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", requestMap)
	}

	decoder := json.NewDecoder(conn)
	responseMap := make(map[string]string)
	err = decoder.Decode(&responseMap)
	if err != nil {
		log.Fatal("Error on decoding response: ", err)
	}

	if responseMap["result"] != "OK" {
		messageDialog := gtk.MessageDialogNew(mainWindow,
			gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK,
			responseMap["error"])
		messageDialog.Run()
		messageDialog.Destroy()
		conn.Close()
		return
	} else {
		getOrder()
	}
}

func discountDeleteButtonClicked() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
		goposServerPort))

	if err != nil {
		log.Fatal("Unable to connect to server")
	}

	requestMap := make(map[string]string)
	requestMap["group"] = "ORDER"
	requestMap["action"] = "DELETE DISCOUNT"
	requestMap["password"] = goposServerPassword
	requestMap["table_namber"] = fmt.Sprintf("%d", orderedTableNumber)
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(requestMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", requestMap)
	}

	decoder := json.NewDecoder(conn)
	responseMap := make(map[string]string)
	err = decoder.Decode(&responseMap)
	if err != nil {
		log.Fatal("Error on decoding response: ", err)
	}

	if responseMap["result"] != "OK" {
		messageDialog := gtk.MessageDialogNew(mainWindow,
			gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK,
			responseMap["error"])
		messageDialog.Run()
		messageDialog.Destroy()
		conn.Close()
		return
	} else {
		getOrder()
	}
}

func existingOrderUpdateButtonClicked() {
	existingOrderMap := make(map[int]int)

	iter, isNotEmpty := existingOrderListStore.GetIterFirst()
	if !isNotEmpty {
		return
	}

	for {
		value, err := existingOrderListStore.GetValue(iter, COLUMN_NEW_ORDER_LIST_DISH_ID)
		if err != nil {
			log.Fatal("Error on getting value: ", err)
		}
		dishId := value.GetInt()
		existingOrderMap[dishId]++
		if !existingOrderListStore.IterNext(iter) {
			break
		}
	}

	orderString := ""
	for id, count := range existingOrderMap {
		orderString += fmt.Sprintf(" %d:%d", id, count)
	}
	orderString = strings.Trim(orderString, " ")
	fmt.Println(existingOrderMap)
	fmt.Println(orderString)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
		goposServerPort))

	if err != nil {
		log.Fatal("Unable to connect to server")
	}

	requestMap := make(map[string]string)
	requestMap["group"] = "ORDER"
	requestMap["action"] = "UPDATE"
	requestMap["password"] = goposServerPassword
	requestMap["order_string"] = orderString
	requestMap["table_number"] = fmt.Sprintf("%d", orderedTableNumber)
	requestMap["price"] = fmt.Sprintf("%f", existingOrderPrice)
	encoder := json.NewEncoder(conn)
	err = encoder.Encode(requestMap)
	if err != nil {
		log.Fatal("Error on encode request map: ", requestMap)
	}

	decoder := json.NewDecoder(conn)
	responseMap := make(map[string]string)
	err = decoder.Decode(&responseMap)
	if err != nil {
		log.Fatal("Error on decoding response: ", err)
	}

	if responseMap["result"] != "OK" {
		messageDialog := gtk.MessageDialogNew(mainWindow,
			gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK,
			responseMap["error"])
		messageDialog.Run()
		messageDialog.Destroy()
		conn.Close()
		return
	} else {
		existingOrderWindow.Destroy()
		existingOrderListWindow.Destroy()
	}
}

func existingOrderListCreateWindow() *gtk.Window {
	var err error
	existingOrderListWindow, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	existingOrderListWindow.SetTitle("gopos-monitor-existingorder-lsit")
	existingOrderListWindow.SetDefaultSize(800, 600)

	existingOrderPriceLabel, err = gtk.LabelNew("")
	if err != nil {
		log.Fatal("Error on creating price label")
	}

	existingOrderVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create existing order vertical box: ", err)
	}

	createExistingOrderTreeView()
	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating workers scrolled window")
	}
	scrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scrolledWindow.Add(existingOrderTreeView)

	passwordFormHbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		log.Fatal("Unable to create password form horizontal box: ", err)
	}

	passwordLabel, err := gtk.LabelNew("Пароль:")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	passwordEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry: ", err)
	}
	passwordEntry.SetVisibility(false)

	cardFormHbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		log.Fatal("Unable to create cards form horizontal box: ", err)
	}

	passwordFormHbox.PackStart(passwordLabel, false, false, 3)
	passwordFormHbox.PackStart(passwordEntry, true, true, 3)

	cardNumberLabel, err := gtk.LabelNew("Номер карты:")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	cardNumberEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry: ", err)
	}

	discountAddButton, err := gtk.ButtonNewWithLabel("Добавить скидку")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
	discountAddButton.Connect("clicked", discountAddButtonClicked, cardNumberEntry)

	discountDeleteButton, err := gtk.ButtonNewWithLabel("Удалить скидку")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
	discountDeleteButton.Connect("clicked", discountDeleteButtonClicked)

	cardFormHbox.PackStart(cardNumberLabel, false, false, 3)
	cardFormHbox.PackStart(cardNumberEntry, true, true, 3)
	cardFormHbox.PackStart(discountAddButton, true, true, 3)
	cardFormHbox.PackStart(discountDeleteButton, true, true, 3)

	dishDeleteButton, err := gtk.ButtonNewWithLabel("Удалить из заказа")
	if err != nil {
		log.Fatal("Unable to create delete button: ", err)
	}
	dishDeleteButton.Connect("clicked", dishExistingOrderDeleteSelectedButtonClicked,
		passwordEntry)

	existingOrderUpdate, err := gtk.ButtonNewWithLabel("Обновить заказ")
	if err != nil {
		log.Fatal("Unable to create update button: ", err)
	}
	existingOrderUpdate.Connect("clicked", existingOrderUpdateButtonClicked)

	existingOrderClose, err := gtk.ButtonNewWithLabel("Закрыть заказ")
	if err != nil {
		log.Fatal("Unable to create close button: ", err)
	}
	existingOrderClose.Connect("clicked", existingOrderCloseButtonClicked, passwordEntry)

	existingOrderVbox.PackStart(scrolledWindow, true, true, 3)
	existingOrderVbox.PackStart(existingOrderPriceLabel, false, false, 3)
	existingOrderVbox.PackStart(passwordFormHbox, false, true, 3)
	existingOrderVbox.PackStart(cardFormHbox, false, true, 3)
	existingOrderVbox.PackStart(dishDeleteButton, false, true, 3)
	existingOrderVbox.PackStart(existingOrderUpdate, false, true, 3)
	existingOrderVbox.PackStart(existingOrderClose, false, true, 3)

	existingOrderListWindow.Add(existingOrderVbox)

	return existingOrderListWindow
}
