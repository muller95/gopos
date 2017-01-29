//main package for gopos database client
//this file contains common functions
package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func createColumn(title string, id int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer: ", err)
	}

	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		log.Fatal("Unable to create cell column")
	}

	return column
}
