package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	db := database{data: map[string]dollars{"shoes": 50, "socks": 5}}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)

	//Additional Handlers
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/create", db.create)
	mux.HandleFunc("/read", db.read)
	mux.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database struct {
	data map[string]dollars
	mu   sync.RWMutex
}

// Prind out the entire database
func (db *database) list(w http.ResponseWriter, req *http.Request) {

	//Lock Before Reading
	db.mu.RLock()
	for item, price := range db.data {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
	db.mu.RUnlock()
}

// Retrieve the price of an item in the database
func (db *database) price(w http.ResponseWriter, req *http.Request) {
	/*
	 * Retrieve the item name from the url
	 * Make sure the item is in the database
	 * Print out the item's price
	 */
	item := req.URL.Query().Get("item")

	// Lock for reading
	db.mu.RLock()
	price, ok := db.data[item]
	db.mu.Unlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}

	fmt.Fprintf(w, "%s\n", price)

}

// Update an item's price in the database
func (db *database) update(w http.ResponseWriter, req *http.Request) {

	/*
	 * Retrieve the item name from the url
	 * Retrieve the new price from the url
	 * Make sure the item is in the database
	 * Make sure the new price is valid
	 * Update the price of the item in the database
	 */

	item := req.URL.Query().Get("item")
	updatedPrice := req.URL.Query().Get("price")

	price, err := strconv.ParseFloat(updatedPrice, 32)

	if err != nil {
		fmt.Println("Error encountered: ", err)
		fmt.Fprintf(w, "'%s' is not a valid price, please try again.\n", updatedPrice)
		return
	}
	//Lock for reading
	db.mu.RLock()
	_, ok := db.data[item]
	db.mu.RUnlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}

	//Lock for writing
	db.mu.Lock()
	db.data[item] = dollars(price)
	db.mu.Unlock()

	fmt.Fprintf(w, "Updated the price for %s to %s\n", item, db.data[item])
}

// Creates a new item and price in the database
func (db *database) create(w http.ResponseWriter, req *http.Request) {

	/*
	 * Retrieve the item name from the url
	 * Retrieve the price from the url
	 * Add the new item and price of the item to the database
	 */

	item := req.URL.Query().Get("item")
	PRICE := req.URL.Query().Get("price")

	price, err := strconv.ParseFloat(PRICE, 32)

	if err != nil {
		fmt.Println("Error encountered: ", err)
		fmt.Fprintf(w, "'%s' is not a valid price, please try again.\n", PRICE)
		return
	}
	// Lock for writing
	db.mu.Lock()
	db.data[item] = dollars(price)
	db.mu.Unlock()

	fmt.Fprintf(w, "New item added!\nItem Name: %s \nPrice: %s\n", item, db.data[item])
}

// Reads the whole database
func (db *database) read(w http.ResponseWriter, req *http.Request) {
	db.list(w, req)
}

// Deletes and item from the database
func (db *database) delete(w http.ResponseWriter, req *http.Request) {

	/*
	 * Retrieve the item name from the url
	 * Make sure the item is actually in the database
	 * Delete the item from the database
	 */

	item := req.URL.Query().Get("item")
	db.mu.RLock()
	_, ok := db.data[item]
	db.mu.RUnlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}

	// Lock for writing
	db.mu.Lock()
	delete(db.data, item)
	db.mu.Unlock()

	fmt.Fprintf(w, "Item %s Removed!\n", item)
}
