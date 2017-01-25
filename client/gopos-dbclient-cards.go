//main package for gopos database client
//this file contains functions for cards tree view
package main

import (
//	"encoding/json"
//	"fmt"
//	"io"
	"log"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
//	"net"
	"strconv"
	"strings"
)

type CardInfo struct{
	CardNumberEntry *gtk.Entry
	CardHolderNameEntry *gtk.Entry
	CardDiscountEntry *gtk.Entry
}

var cardsTreeView *gtk.TreeView
var cardsListStore *gtk.ListStore

const (
	COLUMN_CARDS_NUMBER = iota
	COLUMN_CARDS_HOLDER_NAME
	COLUMN_CARDS_DISCOUNT
)

func createCardsTreeView() {
	var err error
	cardsTreeView, err = gtk.TreeViewNew()
	if err != nil {
		log.Fatal("Unable to create tree view: ", err)
	}

	cardsTreeView.AppendColumn(createColumn("Номер карты", COLUMN_CARDS_NUMBER))
	cardsTreeView.AppendColumn(createColumn("Имя валадельца", COLUMN_CARDS_HOLDER_NAME))
	cardsTreeView.AppendColumn(createColumn("Скидка", COLUMN_CARDS_DISCOUNT))

	cardsListStore, err = gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING,
		glib.TYPE_DOUBLE)
	if err != nil {
		log.Fatal("Unable to create list store: ", err)
	}

	cardsTreeView.SetModel(cardsListStore)
}

func cardAddRow(number string, holderName string, discount float64) {
	iter := cardsListStore.Append()

	err := cardsListStore.Set(iter, []int{ COLUMN_CARDS_NUMBER, COLUMN_CARDS_HOLDER_NAME,
		COLUMN_CARDS_DISCOUNT }, []interface{} { number, holderName, discount })

	if err != nil {
		log.Fatal("Unable to add row: ", err)
	}
}

func cardAddButtonClicked(btn *gtk.Button, cardInfo *CardInfo) {
	cardNumber, err := cardInfo.CardNumberEntry.GetText()
	if err != nil {
		log.Fatal("Unable to get card number entry text: ", err)
	}
	cardNumber = strings.Trim(cardNumber, " ")

	cardHolderName, err := cardInfo.CardHolderNameEntry.GetText()
	if err != nil {
		log.Fatal("Unable to get card holder name entry text: ", err)
	}
	cardHolderName = strings.Trim(cardHolderName, " ")

	cardDiscount, err := cardInfo.CardDiscountEntry.GetText()
	if err != nil {
		log.Fatal("Unable to get card discount entry text: ", err)
	}
	cardDiscount = strings.Trim(cardDiscount, " ")

	errMessage := ""
	if len(cardNumber) == 0 {
		errMessage += "Введите номер карты. "
	}

	if len(cardHolderName) == 0 {
		errMessage += "Введите имя владельца карты. "
	}

	discount := 0.0
	if len(cardDiscount) == 0 {
		errMessage += "Введите скидку карты."
	} else {
		discount, err = strconv.ParseFloat(cardDiscount, 64)
		if err != nil  || discount < 0.0 || discount >= 1.0{
			errMessage += "Скидка должна быть положительным вещественным числом, " +
				"непревосходящим единицы, например, 0.15."
		}
	}

	if errMessage != "" {
		messageDialog := gtk.MessageDialogNew(mainWindow, gtk.DIALOG_MODAL,
			gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, errMessage)
		messageDialog.Run()
		messageDialog.Destroy()
		return
	}

/*	path := rows.Data().(*gtk.TreePath)
	iter, err := categoriesListStore.GetIter(path)
	if err != nil {
		log.Fatal("Error on getting iter: ", err)
	}
	value, err := categoriesListStore.GetValue(iter, 0)
	if err != nil {
		log.Fatal("Error on getting value: ", err)
	}
	categoryId := value.GetInt()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", goposServerIp,
		goposServerPort))

	if err != nil {
		log.Fatal("Unable to connect to server")
	}

	requestMap := make(map[string]string)
	requestMap["group"] = "DISH"
	requestMap["action"] = "ADD"
	requestMap["password"] = goposServerPassword
	requestMap["name"] = dishName
	requestMap["price"] = fmt.Sprintf("%f", price)
	requestMap["category_id"] = fmt.Sprintf("%d", categoryId)
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

	if responseMap["result"] == "ERR" {
		messageDialog := gtk.MessageDialogNew(mainWindow,
			gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK,
			responseMap["error"])
			messageDialog.Run()
		messageDialog.Destroy()
		return
	}

	id, err := strconv.Atoi(responseMap["id"])
	if err != nil {
		log.Fatal("Error on converting id to int: ", err)
	}
	dishAddRow(id, dishName, price)
	conn.Close()*/
}

func cardsCreatePage() *gtk.Box {
	//creates cards tabpage
	cardsVbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		log.Fatal("Unable to create main vertical box: ", err)
	}

	createCardsTreeView()
	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatalf("Error on creating cards scrolled window")
	}
	scrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scrolledWindow.Add(cardsTreeView)

	cardsFormHbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	if err != nil {
		log.Fatal("Unable to create cards form horizontal box: ", err)
	}

	cardNumberLabel, err := gtk.LabelNew("Номер карты:")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	cardNumberEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry: ", err)
	}

	cardHolderNameLabel, err := gtk.LabelNew("Имя владельца карты:")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	cardHolderNameEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry: ", err)
	}

	cardDiscountLabel, err := gtk.LabelNew("Скидка:")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	cardDiscountEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry: ", err)
	}

	cardAddButton, err := gtk.ButtonNewWithLabel("Добавить")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
	cardInfo := new(CardInfo)
	cardInfo.CardNumberEntry = cardNumberEntry
	cardInfo.CardHolderNameEntry = cardHolderNameEntry
	cardInfo.CardDiscountEntry = cardDiscountEntry
	cardAddButton.Connect("clicked", cardAddButtonClicked, cardInfo)

	cardDeleteSelectedButton, err := gtk.ButtonNewWithLabel("Удалить карту")
	if err != nil {
		log.Fatal("Unable to create add button: ", err)
	}
//	cardDeleteSelectedButton.Connect("clicked", cardDeleteSelectedButtonClicked, nil)

	cardsFormHbox.PackStart(cardNumberLabel, false, false, 3)
	cardsFormHbox.PackStart(cardNumberEntry, true, true, 3)
	cardsFormHbox.PackStart(cardHolderNameLabel, false, false, 3)
	cardsFormHbox.PackStart(cardHolderNameEntry, true, true, 3)
	cardsFormHbox.PackStart(cardDiscountLabel, false, false, 3)
	cardsFormHbox.PackStart(cardDiscountEntry, true, true, 3)
	cardsFormHbox.PackStart(cardAddButton, true, true, 3)
	cardsFormHbox.PackStart(cardDeleteSelectedButton, true, true, 3)


	cardsVbox.PackStart(scrolledWindow, true, true, 3)
	cardsVbox.PackStart(cardsFormHbox, false, false, 3)

	return cardsVbox
}
